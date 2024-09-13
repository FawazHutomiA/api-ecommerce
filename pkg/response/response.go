package response

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type RetError struct {
	Msg string `json:"msg"`
}

type BaseResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Errors  interface{} `json:"errors"`
	Data    interface{} `json:"data"`
}

type Response interface {
	JSON(w http.ResponseWriter) (err error)
}

func Success(status string, message string, data interface{}) (resp Response) {
	return &BaseResponse{
		Status:  status,
		Message: message,
		Errors:  nil,
		Data:    data,
	}
}

func Error(status string, message string, err interface{}) (resp Response) {
	switch v := err.(type) {
	case string:
		err = []RetError{{v}}
	case error:
		err = []RetError{{fmt.Sprintf("%v", err)}}
	case interface{}:
	default:
		err = fmt.Errorf("unknown error type")
	}

	return &BaseResponse{
		Status:  status,
		Message: message,
		Errors:  err,
		Data:    nil,
	}
}

func (r *BaseResponse) getStatusCode(status string) (statusCode int) {
	switch status {
	case StatusOK:
		return http.StatusOK
	case StatusCreated:
		return http.StatusCreated
	case StatusBadRequest:
		return http.StatusBadRequest
	case StatusUnauthorized:
		return http.StatusUnauthorized
	case StatusForbiddend:
		return http.StatusForbidden
	case StatusNotFound:
		return http.StatusNotFound
	case StatusConflicted:
		return http.StatusConflict
	case StatusUnprocessableEntity:
		return http.StatusUnprocessableEntity
	case StatusInternalServerError:
		return http.StatusInternalServerError
	case StatusMethodNotAllowed:
		return http.StatusMethodNotAllowed
	default:
		return http.StatusInternalServerError
	}
}

func (r *BaseResponse) JSON(w http.ResponseWriter) error {
	statusCode := r.getStatusCode(r.Status)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	return json.NewEncoder(w).Encode(r)
}
