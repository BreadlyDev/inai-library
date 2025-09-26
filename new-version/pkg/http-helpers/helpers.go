package httphelpers

import (
	"net/http"
	"strconv"
)

type Response struct {
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
	Status  int    `json:"status"`
}

func NewResponse(message string, data any, status int) Response {
	return Response{
		Message: message,
		Data:    data,
		Status:  status,
	}
}

func NewErrResponse(message string, status int) Response {
	return Response{
		Message: message,
		Data:    nil,
		Status:  status,
	}
}

func ParseIdFromPath(r *http.Request) (int, error) {
	id, err := strconv.Atoi(r.PathValue("id"))

	if err != nil {
		return -1, err
	}

	return id, nil
}
