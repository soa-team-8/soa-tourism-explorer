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

func (e *HttpUtils) GetIDFromRequest(req *http.Request, paramName string) (uint64, error) {
	vars := mux.Vars(req)
	paramValueStr := vars[paramName]
	return strconv.ParseUint(paramValueStr, 10, 64)
}

func (e *HttpUtils) GetDoubleFromForm(req *http.Request, paramName string) (float64, error) {
	paramValueStr := req.FormValue(paramName)
	return strconv.ParseFloat(paramValueStr, 64)
}

func (e *HttpUtils) GetUint64SliceFromForm(req *http.Request, paramName string) ([]uint64, error) {
	paramValues := req.Form[paramName]
	if len(paramValues) == 0 {
		return nil, fmt.Errorf("parameter %s not found", paramName)
	}

	var uint64Slice []uint64
	for _, paramValue := range paramValues {
		val, err := strconv.ParseUint(paramValue, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("error parsing parameter %s: %v", paramName, err)
		}
		uint64Slice = append(uint64Slice, val)
	}

	return uint64Slice, nil
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
