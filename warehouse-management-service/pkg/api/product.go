package api

import wms "warehouse-management-service"

type ProductResponse struct {
	Message  string      `json:"message,omitempty"`
	Error    string      `json:"error,omitempty"`
	Response wms.Product `json:"response,omitempty"`
}
