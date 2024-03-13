package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type HttpUtils struct{}

func (e *HttpUtils) WriteJSONResponse(resp http.ResponseWriter, statusCode int, data interface{}) {
	resp.Header().Set("Content-Type", "application/json")
	resp.WriteHeader(statusCode)
	err := json.NewEncoder(resp).Encode(data)
	if err != nil {
		e.HandleError(resp, fmt.Errorf("failed to encode JSON response: %v", err), http.StatusInternalServerError)
	}
}

func (e *HttpUtils) WriteResponse(resp http.ResponseWriter, statusCode int, message string) {
	resp.WriteHeader(statusCode)
	_, err := resp.Write([]byte(message))
	if err != nil {
		e.HandleError(resp, fmt.Errorf("failed to write response: %v", err), http.StatusInternalServerError)
	}
}

func (e *HttpUtils) Decode(body io.Reader, v interface{}) (interface{}, error) {
	err := json.NewDecoder(body).Decode(v)
	return v, err
}

func (e *HttpUtils) GetIDFromRequest(req *http.Request) (uint64, error) {
	vars := mux.Vars(req)
	idStr := vars["id"]
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse ID: %v", err)
	}
	return id, nil
}

func (e *HttpUtils) HandleError(resp http.ResponseWriter, err error, statusCode int) {
	http.Error(resp, fmt.Sprintf("Error: %v", err), statusCode)
}
