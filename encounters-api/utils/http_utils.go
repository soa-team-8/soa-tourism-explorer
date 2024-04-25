package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

type HttpUtils struct{}

func (e *HttpUtils) Decode(body io.Reader, v interface{}) (interface{}, error) {
	err := json.NewDecoder(body).Decode(v)
	return v, err
}

func (e *HttpUtils) GetUInt64FromRequest(req *http.Request, paramName string) (uint64, error) {
	vars := mux.Vars(req)
	paramValueStr := vars[paramName]
	return strconv.ParseUint(paramValueStr, 10, 64)
}

func (e *HttpUtils) GetIntFromRequest(req *http.Request, paramName string) (int, error) {
	vars := mux.Vars(req)
	paramValueStr := vars[paramName]
	return strconv.Atoi(paramValueStr)
}

func (e *HttpUtils) GetDoubleFromForm(req *http.Request, paramName string) (float64, error) {
	paramValueStr := req.FormValue(paramName)
	return strconv.ParseFloat(paramValueStr, 64)
}

func (e *HttpUtils) GetIntFromForm(req *http.Request, paramName string) (int, error) {
	paramValueStr := req.FormValue(paramName)
	return strconv.Atoi(paramValueStr)
}

func (e *HttpUtils) GetObjectFromForm(req *http.Request, paramName string, obj interface{}) error {
	jsonData := req.FormValue(paramName)

	if err := json.NewDecoder(strings.NewReader(jsonData)).Decode(obj); err != nil {
		return err
	}

	return nil
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

func (e *HttpUtils) GetFilesFromForm(req *http.Request, paramName string) ([]*multipart.FileHeader, error) {
	files, ok := req.MultipartForm.File[paramName]
	if !ok {
		return nil, fmt.Errorf("no files found for parameter %s", paramName)
	}

	return files, nil
}

func (e *HttpUtils) HandleError(resp http.ResponseWriter, err error, statusCode int) {
	http.Error(resp, fmt.Sprintf("Error: %v", err), statusCode)
}

func (e *HttpUtils) WriteJSONResponse(resp http.ResponseWriter, statusCode int, data interface{}) {
	resp.Header().Set("Content-Type", "application/json")
	resp.WriteHeader(statusCode)
	if err := json.NewEncoder(resp).Encode(data); err != nil {
		log.Printf("Error encoding JSON response: %v", err)
		http.Error(resp, "Error encoding JSON response", http.StatusInternalServerError)
	}
}

func (e *HttpUtils) WriteResponse(resp http.ResponseWriter, statusCode int, message string) {
	resp.WriteHeader(statusCode)
	_, err := resp.Write([]byte(message))
	if err != nil {
		log.Printf("Error writing response: %v", err)
	}
}

func (e *HttpUtils) GetDoubleFromQuery(req *http.Request, paramName string) (float64, error) {
	paramValueStr := req.URL.Query().Get(paramName)
	return strconv.ParseFloat(paramValueStr, 64)
}

func (e *HttpUtils) GetUint64SliceFromQuery(req *http.Request, paramName string) ([]uint64, error) {
	paramValues := req.URL.Query()[paramName]
	if len(paramValues) == 0 {
		return nil, fmt.Errorf("parameter %s not found", paramName)
	}

	var uint64Slice []uint64
	for _, paramValue := range paramValues {
		values := strings.Split(paramValue, ",")
		for _, val := range values {
			uintVal, err := strconv.ParseUint(val, 10, 64)
			if err != nil {
				return nil, fmt.Errorf("error parsing parameter %s: %v", paramName, err)
			}
			uint64Slice = append(uint64Slice, uintVal)
		}
	}

	return uint64Slice, nil
}
