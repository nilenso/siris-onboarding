package handler

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/golang/mock/gomock"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	wms "warehouse-management-service"
	mock "warehouse-management-service/internal/handler/mock"
	"warehouse-management-service/pkg/api"
	"warehouse-management-service/pkg/log"
)

var h *handler

var warehouse = wms.Warehouse{
	Id:        "85bd3b85-ad4d-4224-b589-fb2a80a6ce45",
	Name:      "test_get_by_id",
	Latitude:  12.9716,
	Longitude: 77.5946,
}

func TestMain(m *testing.M) {
	logger := log.New()
	logger.SetLevel("debug")
	h = &handler{
		logger: logger,
	}

	os.Exit(m.Run())
}

func TestPing(t *testing.T) {
	request, err := http.NewRequest(
		"GET",
		"/ping",
		nil,
	)
	if err != nil {
		t.Error(err)
	}
	response := executeRequest(request)
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		t.Error(err)
	}
	if string(responseBody) != "pong" {
		t.Errorf("want: %v, got: %v", "pong", string(responseBody))
	}
}

func TestGetWarehouseById(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockObj := mock.NewMockWarehouseService(mockCtrl)

	tests := []struct {
		getWarehouseByIdRequest string
		warehouseByIdResponse   *wms.Warehouse
		warehouseByIdErr        error
		wantStatusCode          int
		wantResponse            interface{}
	}{
		{
			getWarehouseByIdRequest: warehouse.Id,
			warehouseByIdResponse:   &warehouse,
			warehouseByIdErr:        nil,
			wantStatusCode:          http.StatusOK,
			wantResponse:            api.GetWarehouseResponse{Response: warehouse},
		},
		{
			getWarehouseByIdRequest: warehouse.Id,
			warehouseByIdResponse:   nil,
			warehouseByIdErr:        sql.ErrConnDone,
			wantStatusCode:          http.StatusInternalServerError,
			wantResponse:            api.GetWarehouseResponse{Error: "Failed to get warehouse"},
		},
		{
			getWarehouseByIdRequest: warehouse.Id,
			warehouseByIdResponse:   nil,
			warehouseByIdErr:        wms.WarehouseDoesNotExist,
			wantStatusCode:          http.StatusNotFound,
			wantResponse: api.GetWarehouseResponse{Error: fmt.Sprintf(
				"failed to get, warehouse: %s does not exist",
				warehouse.Id,
			)},
		},
	}
	for _, test := range tests {

		mockObj.EXPECT().GetWarehouseById(gomock.Any(), test.getWarehouseByIdRequest).Return(test.warehouseByIdResponse, test.warehouseByIdErr)
		h.warehouseService = mockObj

		requestURL := fmt.Sprintf("/warehouse/%s", test.getWarehouseByIdRequest)
		request, err := http.NewRequest(
			"GET",
			requestURL,
			nil,
		)
		if err != nil {
			t.Error(err)
		}

		response := executeRequest(request)
		responseBody, err := io.ReadAll(response.Body)
		if err != nil {
			t.Error(err)
		}
		var got api.GetWarehouseResponse
		err = json.Unmarshal(responseBody, &got)
		if err != nil {
			t.Error(err)
		}

		if response.StatusCode != test.wantStatusCode {
			t.Errorf("want: %v, got: %v", test.wantStatusCode, response.StatusCode)
		}

		if got != test.wantResponse {
			t.Errorf("want: %v, got: %v", test.wantResponse, got)
		}
	}
}

func TestCreateWarehouse(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockObj := mock.NewMockWarehouseService(mockCtrl)

	tests := []struct {
		createWarehouseRequest api.CreateWarehouseRequest
		warehouseByIdErr       error
		wantStatusCode         int
		wantResponse           api.WarehouseResponse
	}{
		{
			createWarehouseRequest: api.CreateWarehouseRequest{
				Name:      "test_create",
				Latitude:  12.989127,
				Longitude: 77.597088,
			},
			warehouseByIdErr: nil,
			wantStatusCode:   http.StatusOK,
			wantResponse: api.WarehouseResponse{
				Response: "Successfully created warehouse: 955fd0a8-1f2e-437f-b988-4b9d3d2acf81",
			},
		},
		{
			createWarehouseRequest: api.CreateWarehouseRequest{
				Name:      "test_create",
				Latitude:  12.989127,
				Longitude: 77.597088,
			},
			warehouseByIdErr: sql.ErrConnDone,
			wantStatusCode:   http.StatusInternalServerError,
			wantResponse:     api.WarehouseResponse{Error: "Failed to create warehouse"},
		},
	}
	for _, test := range tests {

		mockObj.EXPECT().CreateWarehouse(
			gomock.Any(),
			gomock.Any(),
		).Return(test.warehouseByIdErr)
		h.warehouseService = mockObj

		marshalledRequest, err := json.Marshal(test.createWarehouseRequest)
		if err != nil {
			t.Error(err)
		}
		requestBody := bytes.NewBuffer(marshalledRequest)
		request, err := http.NewRequest(
			"POST",
			"/warehouse",
			requestBody,
		)
		if err != nil {
			t.Error(err)
		}

		response := executeRequest(request)
		responseBody, err := io.ReadAll(response.Body)
		if err != nil {
			t.Error(err)
		}
		var got api.WarehouseResponse
		err = json.Unmarshal(responseBody, &got)
		if err != nil {
			t.Error(err)
		}

		if response.StatusCode != test.wantStatusCode {
			t.Errorf("want: %v, got: %v", test.wantStatusCode, response.StatusCode)
		}

		if got.Error != test.wantResponse.Error {
			t.Errorf("want: %v, got: %v", test.wantResponse, got)
		}
	}
}

func TestCreateWarehouseRequestError(t *testing.T) {
	tests := []struct {
		createWarehouseRequest interface{}
		want                   []string
		wantStatusCode         int
	}{
		{
			createWarehouseRequest: nil,
			want:                   []string{"request body cannot be empty"},
			wantStatusCode:         http.StatusBadRequest,
		},
		{
			createWarehouseRequest: map[string]string{"foo": "bar"},
			want:                   []string{"Failed to parse request"},
			wantStatusCode:         http.StatusBadRequest,
		},
		{
			createWarehouseRequest: api.CreateWarehouseRequest{
				Name:      "",
				Latitude:  -90.21,
				Longitude: -180.78,
			},
			want: []string{
				"Name: zero value",
				"Latitude: less than min",
				"Longitude: less than min",
			},
			wantStatusCode: http.StatusBadRequest,
		},
		{
			createWarehouseRequest: api.CreateWarehouseRequest{
				Name:      "test_update",
				Latitude:  100,
				Longitude: 190.89,
			},
			want: []string{
				"Latitude: greater than max",
				"Longitude: greater than max",
			},
			wantStatusCode: http.StatusBadRequest,
		},
	}
	for _, test := range tests {
		var requestBody io.Reader

		if test.createWarehouseRequest != nil {
			marshalledRequest, err := json.Marshal(test.createWarehouseRequest)
			if err != nil {
				t.Error(err)
			}
			requestBody = bytes.NewBuffer(marshalledRequest)
		}

		request, err := http.NewRequest(
			"POST",
			"/warehouse",
			requestBody,
		)
		if err != nil {
			t.Error(err)
		}

		response := executeRequest(request)
		responseBody, err := io.ReadAll(response.Body)
		if err != nil {
			t.Error(err)
		}
		var got api.WarehouseResponse
		err = json.Unmarshal(responseBody, &got)
		if err != nil {
			t.Error(err)
		}

		if response.StatusCode != test.wantStatusCode {
			t.Errorf("want: %v, got: %v", test.wantStatusCode, response.StatusCode)
		}

		for _, want := range test.want {
			if !strings.Contains(got.Error, want) {
				t.Errorf("want: %v, got: %v", test.want, got)
			}
		}
	}
}

func TestUpdateWarehouseRequestError(t *testing.T) {
	tests := []struct {
		updateWarehouseRequest interface{}
		want                   []string
		wantStatusCode         int
	}{
		{
			updateWarehouseRequest: nil,
			want:                   []string{"request body cannot be empty"},
			wantStatusCode:         http.StatusBadRequest,
		},
		{
			updateWarehouseRequest: map[string]string{"foo": "bar"},
			want:                   []string{"Failed to parse request"},
			wantStatusCode:         http.StatusBadRequest,
		},
		{
			updateWarehouseRequest: api.UpdateWarehouseRequest{
				Id:        "",
				Name:      "",
				Latitude:  -90.21,
				Longitude: -180.78,
			},
			want: []string{
				"Latitude: less than min",
				"Longitude: less than min",
				"Id: zero value",
				"Name: zero value",
			},
			wantStatusCode: http.StatusBadRequest,
		},
		{
			updateWarehouseRequest: api.UpdateWarehouseRequest{
				Id:        "85bd3b85-ad4d-4224-b589-fb2a80a6ce45",
				Name:      "test_update",
				Latitude:  100,
				Longitude: 190.89,
			},
			want: []string{
				"Latitude: greater than max",
				"Longitude: greater than max",
			},
			wantStatusCode: http.StatusBadRequest,
		},
	}
	for _, test := range tests {
		var requestBody io.Reader

		if test.updateWarehouseRequest != nil {
			marshalledRequest, err := json.Marshal(test.updateWarehouseRequest)
			if err != nil {
				t.Error(err)
			}
			requestBody = bytes.NewBuffer(marshalledRequest)
		}

		request, err := http.NewRequest(
			"PUT",
			"/warehouse",
			requestBody,
		)
		if err != nil {
			t.Error(err)
		}

		response := executeRequest(request)
		responseBody, err := io.ReadAll(response.Body)
		if err != nil {
			t.Error(err)
		}
		var got api.WarehouseResponse
		err = json.Unmarshal(responseBody, &got)
		if err != nil {
			t.Error(err)
		}

		if response.StatusCode != test.wantStatusCode {
			t.Errorf("want: %v, got: %v", test.wantStatusCode, response.StatusCode)
		}

		for _, want := range test.want {
			if !strings.Contains(got.Error, want) {
				t.Errorf("want: %v, got: %v", want, got)
			}
		}
	}
}

func TestUpdateWarehouse(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockObj := mock.NewMockWarehouseService(mockCtrl)

	tests := []struct {
		updateWarehouseRequest  api.UpdateWarehouseRequest
		updateWarehouseErr      error
		updateWarehouseResponse api.WarehouseResponse
		wantStatusCode          int
	}{
		{
			updateWarehouseRequest: api.UpdateWarehouseRequest{
				Id:        "85bd3b85-ad4d-4224-b589-fb2a80a6ce45",
				Name:      "test_update",
				Latitude:  12.5678,
				Longitude: 77.8901,
			},
			updateWarehouseErr: nil,
			updateWarehouseResponse: api.WarehouseResponse{Response: fmt.Sprintf(
				"Successfully updated warehouse: %s",
				"85bd3b85-ad4d-4224-b589-fb2a80a6ce45",
			)},
			wantStatusCode: http.StatusOK,
		},
		{
			updateWarehouseRequest: api.UpdateWarehouseRequest{
				Id:        "85bd3b85-ad4d-4224-b589-fb2a80a6ce45",
				Name:      "test_update",
				Latitude:  12.5678,
				Longitude: 77.8901,
			},
			updateWarehouseErr:      sql.ErrConnDone,
			updateWarehouseResponse: api.WarehouseResponse{Error: "Failed to update warehouse"},
			wantStatusCode:          http.StatusInternalServerError,
		},
		{
			updateWarehouseRequest: api.UpdateWarehouseRequest{
				Id:        "85bd3b85-ad4d-4224-b589-fb2a80a6ce45",
				Name:      "test_update",
				Latitude:  12.5678,
				Longitude: 77.8901,
			},
			updateWarehouseErr: wms.WarehouseDoesNotExist,
			updateWarehouseResponse: api.WarehouseResponse{Error: fmt.Sprintf(
				"failed to update, warehouse: %s does not exist",
				"85bd3b85-ad4d-4224-b589-fb2a80a6ce45",
			)},
			wantStatusCode: http.StatusNotFound,
		},
	}
	for _, test := range tests {
		warehouse := &wms.Warehouse{
			Id:        test.updateWarehouseRequest.Id,
			Name:      test.updateWarehouseRequest.Name,
			Latitude:  test.updateWarehouseRequest.Latitude,
			Longitude: test.updateWarehouseRequest.Longitude,
		}

		mockObj.EXPECT().UpdateWarehouse(
			gomock.Any(),
			warehouse,
		).Return(test.updateWarehouseErr)
		h.warehouseService = mockObj

		marshalledRequest, err := json.Marshal(test.updateWarehouseRequest)
		if err != nil {
			t.Error(err)
		}
		requestBody := bytes.NewBuffer(marshalledRequest)
		request, err := http.NewRequest(
			"PUT",
			"/warehouse",
			requestBody,
		)
		if err != nil {
			t.Error(err)
		}

		response := executeRequest(request)
		responseBody, err := io.ReadAll(response.Body)
		if err != nil {
			t.Error(err)
		}
		var got api.WarehouseResponse
		err = json.Unmarshal(responseBody, &got)
		if err != nil {
			t.Error(err)
		}

		if response.StatusCode != test.wantStatusCode {
			t.Errorf("want: %v, got: %v", test.wantStatusCode, response.StatusCode)
		}

		if got != test.updateWarehouseResponse {
			t.Errorf("want: %v, got: %v", test.updateWarehouseResponse, got)
		}
	}
}

func TestDeleteWarehouse(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockObj := mock.NewMockWarehouseService(mockCtrl)

	tests := []struct {
		deleteWarehouseRequest string
		deleteWarehouseErr     error
		wantStatusCode         int
		wantResponse           api.WarehouseResponse
	}{
		{
			deleteWarehouseRequest: warehouse.Id,
			deleteWarehouseErr:     nil,
			wantStatusCode:         http.StatusOK,
			wantResponse: api.WarehouseResponse{Response: fmt.Sprintf(
				"Successfully deleted warehouse: %s",
				warehouse.Id,
			)},
		},
		{
			deleteWarehouseRequest: warehouse.Id,
			deleteWarehouseErr:     sql.ErrConnDone,
			wantStatusCode:         http.StatusInternalServerError,
			wantResponse:           api.WarehouseResponse{Error: "Failed to delete warehouse"},
		},
		{
			deleteWarehouseRequest: warehouse.Id,
			deleteWarehouseErr:     wms.WarehouseDoesNotExist,
			wantStatusCode:         http.StatusNotFound,
			wantResponse: api.WarehouseResponse{Error: fmt.Sprintf(
				"failed to delete, warehouse: %s does not exist",
				warehouse.Id,
			)},
		},
	}
	for _, test := range tests {

		mockObj.EXPECT().DeleteWarehouse(gomock.Any(), test.deleteWarehouseRequest).Return(test.deleteWarehouseErr)
		h.warehouseService = mockObj

		requestURL := fmt.Sprintf("/warehouse/%s", test.deleteWarehouseRequest)
		request, err := http.NewRequest(
			"DELETE",
			requestURL,
			nil,
		)
		if err != nil {
			t.Error(err)
		}

		response := executeRequest(request)
		responseBody, err := io.ReadAll(response.Body)
		if err != nil {
			t.Error(err)
		}
		var got api.WarehouseResponse
		err = json.Unmarshal(responseBody, &got)
		if err != nil {
			t.Error(err)
		}

		if response.StatusCode != test.wantStatusCode {
			t.Errorf("want: %v, got: %v", test.wantStatusCode, response.StatusCode)
		}

		if got != test.wantResponse {
			t.Errorf("want: %v, got: %v", test.wantResponse, got)
		}
	}
}

func TestGetShelfBlockById(t *testing.T) {
	shelfBlock := wms.ShelfBlock{
		Id:          "update_test",
		Aisle:       "1",
		Rack:        "1",
		StorageType: "regular",
		WarehouseId: "85bd3b85-ad4d-4224-b589-fb2a80a6ce45",
	}
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockObj := mock.NewMockShelfBlockService(mockCtrl)

	tests := []struct {
		getShelfBlockByIdRequest string
		shelfBlockByIdResponse   wms.ShelfBlock
		shelfBlockByIdErr        error
		wantStatusCode           int
		wantResponse             interface{}
	}{
		{
			getShelfBlockByIdRequest: shelfBlock.Id,
			shelfBlockByIdResponse:   shelfBlock,
			shelfBlockByIdErr:        nil,
			wantStatusCode:           http.StatusOK,
			wantResponse:             api.GetShelfBlockResponse{Response: shelfBlock},
		},
		{
			getShelfBlockByIdRequest: shelfBlock.Id,
			shelfBlockByIdResponse:   wms.ShelfBlock{},
			shelfBlockByIdErr:        sql.ErrConnDone,
			wantStatusCode:           http.StatusInternalServerError,
			wantResponse:             api.GetShelfBlockResponse{Error: "Failed to get shelfBlock"},
		},
		{
			getShelfBlockByIdRequest: shelfBlock.Id,
			shelfBlockByIdResponse:   wms.ShelfBlock{},
			shelfBlockByIdErr:        wms.ShelfBlockDoesNotExist,
			wantStatusCode:           http.StatusNotFound,
			wantResponse: api.GetShelfBlockResponse{Error: fmt.Sprintf(
				"failed to get, shelfBlock: %s does not exist",
				shelfBlock.Id,
			)},
		},
	}
	for _, test := range tests {

		mockObj.EXPECT().GetShelfBlockById(gomock.Any(), test.getShelfBlockByIdRequest).Return(test.shelfBlockByIdResponse, test.shelfBlockByIdErr)
		h.shelfBlockService = mockObj

		requestURL := fmt.Sprintf("/shelf_block/%s", test.getShelfBlockByIdRequest)
		request, err := http.NewRequest(
			"GET",
			requestURL,
			nil,
		)
		if err != nil {
			t.Error(err)
		}

		response := executeRequest(request)
		responseBody, err := io.ReadAll(response.Body)
		if err != nil {
			t.Error(err)
		}
		var got api.GetShelfBlockResponse
		err = json.Unmarshal(responseBody, &got)
		if err != nil {
			t.Error(err)
		}

		if response.StatusCode != test.wantStatusCode {
			t.Errorf("want: %v, got: %v", test.wantStatusCode, response.StatusCode)
		}

		if got != test.wantResponse {
			t.Errorf("want: %v, got: %v", test.wantResponse, got)
		}
	}
}

func TestCreateShelfBlock(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockObj := mock.NewMockShelfBlockService(mockCtrl)

	tests := []struct {
		createShelfBlockRequest api.CreateShelfBlockRequest
		createShelfBlockErr     error
		wantStatusCode          int
		wantResponse            api.ShelfBlockResponse
	}{
		{
			createShelfBlockRequest: api.CreateShelfBlockRequest{
				Aisle:       "2",
				Rack:        "3",
				StorageType: "refrigerated",
				WarehouseId: "xx",
			},
			createShelfBlockErr: nil,
			wantStatusCode:      http.StatusOK,
			wantResponse: api.ShelfBlockResponse{
				Response: "Successfully created shelf_block: ",
			},
		},
		{
			createShelfBlockRequest: api.CreateShelfBlockRequest{
				Aisle:       "2",
				Rack:        "3",
				StorageType: "refrigerated",
				WarehouseId: "xx",
			},
			createShelfBlockErr: sql.ErrConnDone,
			wantStatusCode:      http.StatusInternalServerError,
			wantResponse:        api.ShelfBlockResponse{Error: "Failed to create shelf block"},
		},
		{
			createShelfBlockRequest: api.CreateShelfBlockRequest{
				Aisle:       "2",
				Rack:        "3",
				StorageType: "refrigerated",
				WarehouseId: "xx",
			},
			createShelfBlockErr: wms.InvalidWarehouse,
			wantStatusCode:      http.StatusBadRequest,
			wantResponse: api.ShelfBlockResponse{Error: fmt.Sprintf("%s: %s",
				wms.InvalidWarehouse.Error(),
				"xx",
			)},
		},
	}
	for _, test := range tests {

		mockObj.EXPECT().CreateShelfBlock(
			gomock.Any(),
			gomock.Any(),
		).Return(test.createShelfBlockErr)
		h.shelfBlockService = mockObj

		marshalledRequest, err := json.Marshal(test.createShelfBlockRequest)
		if err != nil {
			t.Error(err)
		}
		requestBody := bytes.NewBuffer(marshalledRequest)
		request, err := http.NewRequest(
			"POST",
			"/shelf_block",
			requestBody,
		)
		if err != nil {
			t.Error(err)
		}

		response := executeRequest(request)
		responseBody, err := io.ReadAll(response.Body)
		if err != nil {
			t.Error(err)
		}
		var got api.ShelfBlockResponse
		err = json.Unmarshal(responseBody, &got)
		if err != nil {
			t.Error(err)
		}

		if response.StatusCode != test.wantStatusCode {
			t.Errorf("want: %v, got: %v", test.wantStatusCode, response.StatusCode)
		}

		if got.Error != test.wantResponse.Error {
			t.Errorf("want: %v, got: %v", test.wantResponse, got)
		}
		if !strings.HasPrefix(got.Response, test.wantResponse.Response) {
			t.Errorf("want: %v, got: %v", test.wantResponse.Response, got.Response)
		}
	}
}

func TestCreateShelfBlockRequestError(t *testing.T) {
	tests := []struct {
		createShelfBlockRequest interface{}
		want                    []string
		wantStatusCode          int
	}{
		{
			createShelfBlockRequest: nil,
			want:                    []string{"request body cannot be empty"},
			wantStatusCode:          http.StatusBadRequest,
		},
		{
			createShelfBlockRequest: map[string]string{"foo": "bar"},
			want:                    []string{"Failed to parse request"},
			wantStatusCode:          http.StatusBadRequest,
		},
		{
			createShelfBlockRequest: api.CreateShelfBlockRequest{
				Aisle:       "",
				Rack:        "",
				StorageType: "",
				WarehouseId: "",
			},
			want: []string{
				"StorageType: zero value",
				"WarehouseId: zero value",
				"Aisle: zero value",
				"Rack: zero value",
			},
			wantStatusCode: http.StatusBadRequest,
		},
	}
	for _, test := range tests {
		var requestBody io.Reader

		if test.createShelfBlockRequest != nil {
			marshalledRequest, err := json.Marshal(test.createShelfBlockRequest)
			if err != nil {
				t.Error(err)
			}
			requestBody = bytes.NewBuffer(marshalledRequest)
		}

		request, err := http.NewRequest(
			"POST",
			"/shelf_block",
			requestBody,
		)
		if err != nil {
			t.Error(err)
		}

		response := executeRequest(request)
		responseBody, err := io.ReadAll(response.Body)
		if err != nil {
			t.Error(err)
		}
		var got api.ShelfBlockResponse
		err = json.Unmarshal(responseBody, &got)
		if err != nil {
			t.Error(err)
		}

		if response.StatusCode != test.wantStatusCode {
			t.Errorf("want: %v, got: %v", test.wantStatusCode, response.StatusCode)
		}

		for _, want := range test.want {
			if !strings.Contains(got.Error, want) {
				t.Errorf("want: %v, got: %v", test.want, got)
			}
		}
	}
}

func TestUpdateShelfBlock(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockObj := mock.NewMockShelfBlockService(mockCtrl)

	request := api.UpdateShelfBlockRequest{
		Id:          "8387eec6-040a-4eb8-b1b5-9277b2d1a72c",
		Aisle:       "2",
		Rack:        "3",
		StorageType: "refrigerated",
		WarehouseId: "xx",
	}

	tests := []struct {
		updateShelfBlockErr      error
		updateShelfBlockResponse api.ShelfBlockResponse
		wantStatusCode           int
	}{
		{
			updateShelfBlockErr: nil,
			updateShelfBlockResponse: api.ShelfBlockResponse{Response: fmt.Sprintf(
				"Successfully updated shelf_block: %s",
				"8387eec6-040a-4eb8-b1b5-9277b2d1a72c",
			)},
			wantStatusCode: http.StatusOK,
		},
		{
			updateShelfBlockErr:      sql.ErrConnDone,
			updateShelfBlockResponse: api.ShelfBlockResponse{Error: "Failed to update shelf_block"},
			wantStatusCode:           http.StatusInternalServerError,
		},
		{
			updateShelfBlockErr: wms.ShelfBlockDoesNotExist,
			updateShelfBlockResponse: api.ShelfBlockResponse{Error: fmt.Sprintf(
				"failed to update, shelf_block: %s does not exist",
				"8387eec6-040a-4eb8-b1b5-9277b2d1a72c",
			)},
			wantStatusCode: http.StatusNotFound,
		},
		{
			updateShelfBlockErr: wms.InvalidWarehouse,
			updateShelfBlockResponse: api.ShelfBlockResponse{Error: fmt.Sprintf("%s: %s",
				wms.InvalidWarehouse.Error(),
				"xx",
			)},
			wantStatusCode: http.StatusBadRequest,
		},
	}

	for _, test := range tests {
		mockObj.EXPECT().UpdateShelfBlock(
			gomock.Any(),
			wms.ShelfBlock{
				Id:          request.Id,
				Aisle:       request.Aisle,
				Rack:        request.Rack,
				StorageType: request.StorageType,
				WarehouseId: request.WarehouseId,
			},
		).Return(test.updateShelfBlockErr)

		h.shelfBlockService = mockObj

		marshalledRequest, err := json.Marshal(request)
		if err != nil {
			t.Error(err)
			return
		}
		requestBody := bytes.NewBuffer(marshalledRequest)
		request, err := http.NewRequest(
			"PUT",
			"/shelf_block",
			requestBody,
		)
		if err != nil {
			t.Error(err)
			return
		}

		response := executeRequest(request)
		responseBody, err := io.ReadAll(response.Body)
		if err != nil {
			t.Error(err)
			return
		}
		var got api.ShelfBlockResponse
		err = json.Unmarshal(responseBody, &got)
		if err != nil {
			t.Error(err)
		}

		if response.StatusCode != test.wantStatusCode {
			t.Errorf("want: %v, got: %v", test.wantStatusCode, response.StatusCode)
		}

		if got.Error != test.updateShelfBlockResponse.Error {
			t.Errorf("want: %v, got: %v", test.updateShelfBlockResponse, got)
		}
	}
}

func TestUpdateShelfBlockRequestError(t *testing.T) {
	tests := []struct {
		updateWarehouseRequest interface{}
		want                   []string
		wantStatusCode         int
	}{
		{
			updateWarehouseRequest: nil,
			want:                   []string{"request body cannot be empty"},
			wantStatusCode:         http.StatusBadRequest,
		},
		{
			updateWarehouseRequest: map[string]string{"foo": "bar"},
			want:                   []string{"Failed to parse request"},
			wantStatusCode:         http.StatusBadRequest,
		},
		{
			updateWarehouseRequest: api.UpdateShelfBlockRequest{},
			want: []string{
				"Id: zero value",
				"Aisle: zero value",
				"Rack: zero value",
				"StorageType: zero value",
				"WarehouseId: zero value",
			},
			wantStatusCode: http.StatusBadRequest,
		},
		{
			updateWarehouseRequest: api.UpdateShelfBlockRequest{
				Id:          "2db3d156-c2c3-4cd7-9af1-cf20047083e8",
				Aisle:       "2",
				Rack:        "3",
				StorageType: "",
				WarehouseId: "",
			},
			want: []string{
				"StorageType: zero value",
				"WarehouseId: zero value",
			},
			wantStatusCode: http.StatusBadRequest,
		},
	}

	for _, test := range tests {
		var requestBody io.Reader

		if test.updateWarehouseRequest != nil {
			marshalledRequest, err := json.Marshal(test.updateWarehouseRequest)
			if err != nil {
				t.Error(err)
			}
			requestBody = bytes.NewBuffer(marshalledRequest)
		}

		request, err := http.NewRequest(
			"PUT",
			"/shelf_block",
			requestBody,
		)
		if err != nil {
			t.Error(err)
		}

		response := executeRequest(request)
		responseBody, err := io.ReadAll(response.Body)
		if err != nil {
			t.Error(err)
		}
		var got api.ShelfBlockResponse
		err = json.Unmarshal(responseBody, &got)
		if err != nil {
			t.Error(err)
		}

		if response.StatusCode != test.wantStatusCode {
			t.Errorf("want: %v, got: %v", test.wantStatusCode, response.StatusCode)
		}

		for _, want := range test.want {
			if !strings.Contains(got.Error, want) {
				t.Errorf("want: %v, got: %v", test.want, got)
			}
		}
	}
}

func TestDeleteShelfBlock(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockObj := mock.NewMockShelfBlockService(mockCtrl)

	tests := []struct {
		deleteShelfBlockRequest string
		deleteShelfBlockErr     error
		wantStatusCode          int
		wantResponse            api.ShelfBlockResponse
	}{
		{
			deleteShelfBlockRequest: "8387eec6-040a-4eb8-b1b5-9277b2d1a72c",
			deleteShelfBlockErr:     nil,
			wantStatusCode:          http.StatusOK,
			wantResponse: api.ShelfBlockResponse{Response: fmt.Sprintf(
				"Successfully deleted shelf_block: %s",
				"8387eec6-040a-4eb8-b1b5-9277b2d1a72c",
			)},
		},
		{
			deleteShelfBlockRequest: "8387eec6-040a-4eb8-b1b5-9277b2d1a72c",
			deleteShelfBlockErr:     sql.ErrConnDone,
			wantStatusCode:          http.StatusInternalServerError,
			wantResponse:            api.ShelfBlockResponse{Error: "Failed to delete shelf_block"},
		},
		{
			deleteShelfBlockRequest: "8387eec6-040a-4eb8-b1b5-9277b2d1a72c",
			deleteShelfBlockErr:     wms.ShelfBlockDoesNotExist,
			wantStatusCode:          http.StatusNotFound,
			wantResponse: api.ShelfBlockResponse{Error: fmt.Sprintf(
				"failed to delete, shelf_block: %s does not exist",
				"8387eec6-040a-4eb8-b1b5-9277b2d1a72c",
			)},
		},
	}
	for _, test := range tests {

		mockObj.EXPECT().DeleteShelfBlockById(gomock.Any(), test.deleteShelfBlockRequest).Return(test.deleteShelfBlockErr)
		h.shelfBlockService = mockObj

		requestURL := fmt.Sprintf("/shelf_block/%s", test.deleteShelfBlockRequest)
		request, err := http.NewRequest(
			"DELETE",
			requestURL,
			nil,
		)
		if err != nil {
			t.Error(err)
		}

		response := executeRequest(request)
		responseBody, err := io.ReadAll(response.Body)
		if err != nil {
			t.Error(err)
		}
		var got api.ShelfBlockResponse
		err = json.Unmarshal(responseBody, &got)
		if err != nil {
			t.Error(err)
		}

		if response.StatusCode != test.wantStatusCode {
			t.Errorf("want: %v, got: %v", test.wantStatusCode, response.StatusCode)
		}

		if got != test.wantResponse {
			t.Errorf("want: %v, got: %v", test.wantResponse, got)
		}
	}
}

func TestGetShelfById(t *testing.T) {
	shelf := wms.Shelf{
		Id:           "get_handler_test",
		Label:        "12A",
		Section:      "A",
		Level:        "12",
		ShelfBlockId: "863e835b-a05b-4554-b0af-a45389ebbb78",
	}
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockObj := mock.NewMockShelfService(mockCtrl)

	tests := []struct {
		getShelfByIdRequest string
		shelfByIdResponse   wms.Shelf
		shelfByIdErr        error
		wantStatusCode      int
		wantResponse        interface{}
	}{
		{
			getShelfByIdRequest: shelf.Id,
			shelfByIdResponse:   shelf,
			shelfByIdErr:        nil,
			wantStatusCode:      http.StatusOK,
			wantResponse:        api.ShelfResponse{Response: shelf},
		},
		{
			getShelfByIdRequest: shelf.Id,
			shelfByIdResponse:   wms.Shelf{},
			shelfByIdErr:        sql.ErrConnDone,
			wantStatusCode:      http.StatusInternalServerError,
			wantResponse:        api.ShelfResponse{Error: "Failed to get shelf"},
		},
		{
			getShelfByIdRequest: shelf.Id,
			shelfByIdResponse:   wms.Shelf{},
			shelfByIdErr:        wms.ShelfDoesNotExist,
			wantStatusCode:      http.StatusNotFound,
			wantResponse: api.ShelfResponse{Error: fmt.Sprintf(
				"failed to get, shelf: %s does not exist",
				shelf.Id,
			)},
		},
	}

	for _, test := range tests {

		mockObj.EXPECT().GetShelfById(
			gomock.Any(),
			test.getShelfByIdRequest,
		).Return(test.shelfByIdResponse, test.shelfByIdErr)
		h.shelfService = mockObj

		requestURL := fmt.Sprintf("/shelf/%s", test.getShelfByIdRequest)
		request, err := http.NewRequest(
			"GET",
			requestURL,
			nil,
		)
		if err != nil {
			t.Error(err)
		}

		response := executeRequest(request)
		responseBody, err := io.ReadAll(response.Body)
		if err != nil {
			t.Error(err)
		}
		var got api.ShelfResponse
		err = json.Unmarshal(responseBody, &got)
		if err != nil {
			t.Error(err)
		}

		if response.StatusCode != test.wantStatusCode {
			t.Errorf("want: %v, got: %v", test.wantStatusCode, response.StatusCode)
		}

		if got != test.wantResponse {
			t.Errorf("want: %v, got: %v", test.wantResponse, got)
		}
	}
}

func TestCreateShelf(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockObj := mock.NewMockShelfService(mockCtrl)

	tests := []struct {
		createShelfRequest wms.Shelf
		createShelfErr     error
		wantStatusCode     int
		wantResponse       api.ShelfResponse
	}{
		{
			createShelfRequest: wms.Shelf{
				Label:        "12A",
				Section:      "A",
				Level:        "12",
				ShelfBlockId: "863e835b-a05b-4554-b0af-a45389ebbb78",
			},
			createShelfErr: nil,
			wantStatusCode: http.StatusOK,
			wantResponse: api.ShelfResponse{
				Message: "Successfully created shelf: ",
			},
		},
		{
			createShelfRequest: wms.Shelf{
				Label:        "12A",
				Section:      "A",
				Level:        "12",
				ShelfBlockId: "863e835b-a05b-4554-b0af-a45389ebbb78",
			},
			createShelfErr: sql.ErrConnDone,
			wantStatusCode: http.StatusInternalServerError,
			wantResponse:   api.ShelfResponse{Error: "Failed to create shelf"},
		},
		{
			createShelfRequest: wms.Shelf{
				Label:        "12A",
				Section:      "A",
				Level:        "12",
				ShelfBlockId: "xx",
			},
			createShelfErr: wms.InvalidShelfBlock,
			wantStatusCode: http.StatusBadRequest,
			wantResponse: api.ShelfResponse{Error: fmt.Sprintf("%s: %s",
				wms.InvalidShelfBlock.Error(),
				"xx",
			)},
		},
	}

	for _, test := range tests {

		mockObj.EXPECT().CreateShelf(
			gomock.Any(),
			gomock.Any(),
		).Return(test.createShelfErr)
		h.shelfService = mockObj

		marshalledRequest, err := json.Marshal(test.createShelfRequest)
		if err != nil {
			t.Error(err)
		}
		requestBody := bytes.NewBuffer(marshalledRequest)
		request, err := http.NewRequest(
			"POST",
			"/shelf",
			requestBody,
		)
		if err != nil {
			t.Error(err)
		}

		response := executeRequest(request)
		responseBody, err := io.ReadAll(response.Body)
		if err != nil {
			t.Error(err)
		}
		var got api.ShelfResponse
		err = json.Unmarshal(responseBody, &got)
		if err != nil {
			t.Error(err)
		}

		if response.StatusCode != test.wantStatusCode {
			t.Errorf("want: %v, got: %v", test.wantStatusCode, response.StatusCode)
		}

		if got.Error != test.wantResponse.Error {
			t.Errorf("want: %v, got: %v", test.wantResponse, got)
		}

		if !strings.HasPrefix(got.Message, test.wantResponse.Message) {
			t.Errorf("want: %v, got: %v", test.wantResponse.Message, got.Message)
		}
	}
}

func TestUpdateShelf(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockObj := mock.NewMockShelfService(mockCtrl)

	request := wms.Shelf{
		Id:           "00c4e998-41df-40d1-ba1e-0bdf870bcc5c",
		Label:        "12A",
		Section:      "A",
		Level:        "12",
		ShelfBlockId: "863e835b-a05b-4554-b0af-a45389ebbb78",
	}

	tests := []struct {
		updateShelfErr      error
		updateShelfResponse api.ShelfResponse
		wantStatusCode      int
	}{
		{
			updateShelfErr: nil,
			updateShelfResponse: api.ShelfResponse{Message: fmt.Sprintf(
				"Successfully updated shelf: %s",
				"00c4e998-41df-40d1-ba1e-0bdf870bcc5c",
			)},
			wantStatusCode: http.StatusOK,
		},
		{
			updateShelfErr:      sql.ErrConnDone,
			updateShelfResponse: api.ShelfResponse{Error: "Failed to update shelf"},
			wantStatusCode:      http.StatusInternalServerError,
		},
		{
			updateShelfErr: wms.ShelfDoesNotExist,
			updateShelfResponse: api.ShelfResponse{Error: fmt.Sprintf(
				"failed to update, shelf: %s does not exist",
				"00c4e998-41df-40d1-ba1e-0bdf870bcc5c",
			)},
			wantStatusCode: http.StatusNotFound,
		},
		{
			updateShelfErr: wms.InvalidShelfBlock,
			updateShelfResponse: api.ShelfResponse{Error: fmt.Sprintf("%s: %s",
				wms.InvalidShelfBlock.Error(),
				"863e835b-a05b-4554-b0af-a45389ebbb78",
			)},
			wantStatusCode: http.StatusBadRequest,
		},
	}

	for _, test := range tests {
		mockObj.EXPECT().UpdateShelf(
			gomock.Any(),
			wms.Shelf{
				Id:           request.Id,
				Label:        request.Label,
				Section:      request.Section,
				Level:        request.Level,
				ShelfBlockId: request.ShelfBlockId,
			},
		).Return(test.updateShelfErr)

		h.shelfService = mockObj

		marshalledRequest, err := json.Marshal(request)
		if err != nil {
			t.Error(err)
			return
		}
		requestBody := bytes.NewBuffer(marshalledRequest)
		request, err := http.NewRequest(
			"PUT",
			"/shelf",
			requestBody,
		)
		if err != nil {
			t.Error(err)
			return
		}

		response := executeRequest(request)
		responseBody, err := io.ReadAll(response.Body)
		if err != nil {
			t.Error(err)
			return
		}
		var got api.ShelfResponse
		err = json.Unmarshal(responseBody, &got)
		if err != nil {
			t.Error(err)
		}

		if response.StatusCode != test.wantStatusCode {
			t.Errorf("want: %v, got: %v", test.wantStatusCode, response.StatusCode)
		}

		if got.Error != test.updateShelfResponse.Error {
			t.Errorf("want: %v, got: %v", test.updateShelfResponse, got)
		}
	}
}

func TestDeleteShelf(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockObj := mock.NewMockShelfService(mockCtrl)
	tests := []struct {
		deleteShelfRequest string
		deleteShelfErr     error
		wantStatusCode     int
		wantResponse       api.ShelfResponse
	}{
		{
			deleteShelfRequest: "8387eec6-040a-4eb8-b1b5-9277b2d1a72c",
			deleteShelfErr:     nil,
			wantStatusCode:     http.StatusOK,
			wantResponse: api.ShelfResponse{Message: fmt.Sprintf(
				"Successfully deleted shelf: %s",
				"8387eec6-040a-4eb8-b1b5-9277b2d1a72c",
			)},
		},
		{
			deleteShelfRequest: "8387eec6-040a-4eb8-b1b5-9277b2d1a72c",
			deleteShelfErr:     sql.ErrConnDone,
			wantStatusCode:     http.StatusInternalServerError,
			wantResponse:       api.ShelfResponse{Error: "Failed to delete shelf"},
		},
		{
			deleteShelfRequest: "8387eec6-040a-4eb8-b1b5-9277b2d1a72c",
			deleteShelfErr:     wms.ShelfDoesNotExist,
			wantStatusCode:     http.StatusNotFound,
			wantResponse: api.ShelfResponse{Error: fmt.Sprintf(
				"failed to delete, shelf: %s does not exist",
				"8387eec6-040a-4eb8-b1b5-9277b2d1a72c",
			)},
		},
	}

	for _, test := range tests {

		mockObj.EXPECT().DeleteShelfById(gomock.Any(), test.deleteShelfRequest).Return(test.deleteShelfErr)
		h.shelfService = mockObj

		requestURL := fmt.Sprintf("/shelf/%s", test.deleteShelfRequest)
		request, err := http.NewRequest(
			"DELETE",
			requestURL,
			nil,
		)
		if err != nil {
			t.Error(err)
		}

		response := executeRequest(request)
		responseBody, err := io.ReadAll(response.Body)
		if err != nil {
			t.Error(err)
		}
		var got api.ShelfResponse
		err = json.Unmarshal(responseBody, &got)
		if err != nil {
			t.Error(err)
		}

		if response.StatusCode != test.wantStatusCode {
			t.Errorf("want: %v, got: %v", test.wantStatusCode, response.StatusCode)
		}

		if got != test.wantResponse {
			t.Errorf("want: %v, got: %v", test.wantResponse, got)
		}
	}
}

func executeRequest(request *http.Request) *http.Response {
	responseRecorder := httptest.NewRecorder()

	h.router().ServeHTTP(responseRecorder, request)
	return responseRecorder.Result()
}
