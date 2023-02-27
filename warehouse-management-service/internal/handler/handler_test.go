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
	"testing"
	warehousemanagementservice "warehouse-management-service"
	mock "warehouse-management-service/internal/handler/mock"
	"warehouse-management-service/pkg/api"
	"warehouse-management-service/pkg/log"
)

var h *handler

var warehouse = warehousemanagementservice.Warehouse{
	Id:        "85bd3b85-ad4d-4224-b589-fb2a80a6ce45",
	Name:      "test_get_by_id",
	Latitude:  12.9716,
	Longitude: 77.5946,
}

var getWarehouseByIdTests = []struct {
	getWarehouseByIdRequest string
	warehouseByIdResponse   *warehousemanagementservice.Warehouse
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
		warehouseByIdErr:        warehousemanagementservice.WarehouseDoesNotExist,
		wantStatusCode:          http.StatusNotFound,
		wantResponse: api.GetWarehouseResponse{Error: fmt.Sprintf(
			"failed to get, warehouse: %s does not exist",
			warehouse.Id,
		)},
	},
}

var createWarehouseTests = []struct {
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

var createWarehouseBadRequestTests = []struct {
	createWarehouseRequest interface{}
	want                   api.WarehouseResponse
	wantStatusCode         int
}{
	{
		createWarehouseRequest: nil,
		want:                   api.WarehouseResponse{Error: "request body cannot be empty"},
		wantStatusCode:         http.StatusBadRequest,
	},
	{
		createWarehouseRequest: map[string]string{"foo": "bar"},
		want:                   api.WarehouseResponse{Error: "Failed to parse request"},
		wantStatusCode:         http.StatusBadRequest,
	},
	{
		createWarehouseRequest: api.CreateWarehouseRequest{
			Name:      "",
			Latitude:  -90.21,
			Longitude: -180.78,
		},
		want: api.WarehouseResponse{Error: "Invalid input: name cannot be empty, latitude has to be in the range [-90, 90], " +
			"longitude has to be in the range [-180, 180]"},
		wantStatusCode: http.StatusBadRequest,
	},
	{
		createWarehouseRequest: api.CreateWarehouseRequest{
			Name:      "test_update",
			Latitude:  100,
			Longitude: 190.89,
		},
		want: api.WarehouseResponse{Error: "Invalid input: latitude has to be in the range [-90, 90], " +
			"longitude has to be in the range [-180, 180]"},
		wantStatusCode: http.StatusBadRequest,
	},
}

var updateWarehouseTests = []struct {
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
		updateWarehouseErr: warehousemanagementservice.WarehouseDoesNotExist,
		updateWarehouseResponse: api.WarehouseResponse{Error: fmt.Sprintf(
			"failed to update, warehouse: %s does not exist",
			"85bd3b85-ad4d-4224-b589-fb2a80a6ce45",
		)},
		wantStatusCode: http.StatusNotFound,
	},
}

var updateWarehouseBadRequestTests = []struct {
	updateWarehouseRequest interface{}
	want                   api.WarehouseResponse
	wantStatusCode         int
}{
	{
		updateWarehouseRequest: nil,
		want:                   api.WarehouseResponse{Error: "request body cannot be empty"},
		wantStatusCode:         http.StatusBadRequest,
	},
	{
		updateWarehouseRequest: map[string]string{"foo": "bar"},
		want:                   api.WarehouseResponse{Error: "Failed to parse request"},
		wantStatusCode:         http.StatusBadRequest,
	},
	{
		updateWarehouseRequest: api.UpdateWarehouseRequest{
			Id:        "",
			Name:      "",
			Latitude:  -90.21,
			Longitude: -180.78,
		},
		want: api.WarehouseResponse{Error: "Invalid input: id cannot be empty, name cannot be empty, latitude has to be in the range [-90, 90], " +
			"longitude has to be in the range [-180, 180]"},
		wantStatusCode: http.StatusBadRequest,
	},
	{
		updateWarehouseRequest: api.UpdateWarehouseRequest{
			Id:        "85bd3b85-ad4d-4224-b589-fb2a80a6ce45",
			Name:      "test_update",
			Latitude:  100,
			Longitude: 190.89,
		},
		want: api.WarehouseResponse{Error: "Invalid input: latitude has to be in the range [-90, 90], " +
			"longitude has to be in the range [-180, 180]"},
		wantStatusCode: http.StatusBadRequest,
	},
}

var deleteWarehouseTests = []struct {
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
		deleteWarehouseErr:     warehousemanagementservice.WarehouseDoesNotExist,
		wantStatusCode:         http.StatusNotFound,
		wantResponse: api.WarehouseResponse{Error: fmt.Sprintf(
			"failed to delete, warehouse: %s does not exist",
			warehouse.Id,
		)},
	},
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

	for _, test := range getWarehouseByIdTests {

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

	for _, test := range createWarehouseTests {

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

		if response.StatusCode != test.wantStatusCode {
			t.Errorf("want: %v, got: %v", test.wantStatusCode, response.StatusCode)
		}

		if got.Error != test.wantResponse.Error {
			t.Errorf("want: %v, got: %v", test.wantResponse, got)
		}
	}
}

func TestCreateWarehouseRequestError(t *testing.T) {
	for _, test := range createWarehouseBadRequestTests {
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

		if response.StatusCode != test.wantStatusCode {
			t.Errorf("want: %v, got: %v", test.wantStatusCode, response.StatusCode)
		}

		if got != test.want {
			t.Errorf("want: %v, got: %v", test.want, got)
		}
	}
}

func TestUpdateWarehouseRequestError(t *testing.T) {
	for _, test := range updateWarehouseBadRequestTests {
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

		if response.StatusCode != test.wantStatusCode {
			t.Errorf("want: %v, got: %v", test.wantStatusCode, response.StatusCode)
		}

		if got != test.want {
			t.Errorf("want: %v, got: %v", test.want, got)
		}
	}
}

func TestUpdateWarehouse(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockObj := mock.NewMockWarehouseService(mockCtrl)

	for _, test := range updateWarehouseTests {
		warehouse := &warehousemanagementservice.Warehouse{
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

	for _, test := range deleteWarehouseTests {

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
