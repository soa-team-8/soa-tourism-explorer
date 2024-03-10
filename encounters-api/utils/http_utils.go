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

func (e *HttpUtils) Decode(body io.Reader, v interface{}) (interface{}, error) {
	err := json.NewDecoder(body).Decode(v)
	return v, err
}

func (e *HttpUtils) GetIDFromRequest(req *http.Request) (uint64, error) {
	vars := mux.Vars(req)
	idStr := vars["id"]
	return strconv.ParseUint(idStr, 10, 64)
}

func (e *HttpUtils) HandleError(resp http.ResponseWriter, err error, statusCode int) {
	http.Error(resp, fmt.Sprintf("Error: %v", err), statusCode)
}

func (e *HttpUtils) WriteJSONResponse(resp http.ResponseWriter, statusCode int, data interface{}) {
	resp.Header().Set("Content-Type", "application/json")
	resp.WriteHeader(statusCode)
	json.NewEncoder(resp).Encode(data)
}

func (e *HttpUtils) WriteResponse(resp http.ResponseWriter, statusCode int, message string) {
	resp.WriteHeader(statusCode)
	resp.Write([]byte(message))
}
