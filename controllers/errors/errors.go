package errors

import (
	"fmt"
	"github.com/donohutcheon/gowebserver/controllers/response/types"
	"net/http"

	"github.com/donohutcheon/gowebserver/controllers/response"
)

type ControllerError struct {
	ErrorMessage string
	Fields []types.ErrorField
	StatusCode int
	Err error
}

func NewError(errorMessage string, fields []types.ErrorField, statusCode int) *ControllerError {
	m := new(ControllerError)
	m.ErrorMessage = errorMessage
	m.Fields = fields
	m.StatusCode = statusCode
	return m
}

func Wrap(errorMessage string, statusCode int, err error) *ControllerError {
	m := new(ControllerError)
	m.ErrorMessage = errorMessage
	m.StatusCode = statusCode
	m.Err = err
	return m
}

func (e *ControllerError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%d %s %s", e.StatusCode, e.ErrorMessage, e.Err)
	} else {
		return fmt.Sprintf("%d %s", e.StatusCode, e.ErrorMessage)
	}
}

func WriteError(w http.ResponseWriter, err error, defaultStatusCode ...int) {
	if err, ok := err.(*ControllerError); ok {
		resp := response.NewWithFieldsList(false, err.ErrorMessage, err.Fields)
		w.WriteHeader(err.StatusCode)
		resp.Respond(w)
		return
	}

	var statusCode = http.StatusInternalServerError
	if len(defaultStatusCode) > 0 {
		statusCode = defaultStatusCode[0]
	}

	resp := response.New(false, err.Error())
	w.WriteHeader(statusCode)
	resp.Respond(w)
}
