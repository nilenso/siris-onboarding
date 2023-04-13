package api

import wms "warehouse-management-service"

type CreateShelfRequest struct {
	Label        string `json:"label"`
	Section      string `json:"section"`
	Level        string `json:"level"`
	ShelfBlockId string `json:"shelfBlockId"`
}

type UpdateShelfRequest struct {
	Id           string `json:"id"`
	Label        string `json:"label"`
	Section      string `json:"section"`
	Level        string `json:"level"`
	ShelfBlockId string `json:"shelfBlockId"`
}

type ShelfResponse struct {
	Message  string    `json:"message,omitempty"`
	Response wms.Shelf `json:"response,omitempty"`
	Error    string    `json:"error,omitempty"`
}
