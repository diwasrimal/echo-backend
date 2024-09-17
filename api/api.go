package api

import (
	"net/http"

	"github.com/diwasrimal/echo-backend/types"
	"github.com/diwasrimal/echo-backend/utils"
)

type Response struct {
	Status  int
	Payload types.Json
}

var NotImplementedResp = Response{
	Status:  http.StatusNotImplemented,
	Payload: types.Json{"message": "Not implemented!"},
}

var InternalErrorResp = Response{
	Status:  http.StatusInternalServerError,
	Payload: types.Json{},
}

// Creates response struct with 400 status code and given message in payload
func BadRequestResp(message string) Response {
	return Response{
		Status:  http.StatusBadRequest,
		Payload: types.Json{"message": message},
	}
}

type apiFunc func(http.ResponseWriter, *http.Request) Response

// Turns our custom api function signatures to http.Handler
func MakeHandler(f apiFunc) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		resp := f(w, r)
		utils.SendJsonResp(w, resp.Status, resp.Payload)
	}
	return http.HandlerFunc(fn)
}
