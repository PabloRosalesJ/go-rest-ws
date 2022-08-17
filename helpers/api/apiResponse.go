package apiResponse

import (
	"encoding/json"
	"net/http"
)

type responseOK struct {
	Status string      `json:"status"`
	Data   interface{} `json:"data"`
}

type responseError struct {
	Status       string   `json:"status"`
	ErrorMessage string   `json:"error_mssg"`
	ErrorData    []string `json:"error_data"`
}

func ResponseOk(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(responseOK{
		Status: "success",
		Data:   data,
	})
}

func ResponseErr(w http.ResponseWriter, status int, errorMssg string, errorData []string) {
	w.Header().Set("ContentType", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(responseError{
		Status:       "error",
		ErrorMessage: errorMssg,
		ErrorData:    errorData,
	})
}
