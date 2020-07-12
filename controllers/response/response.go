package response

import (
	"encoding/json"
	"github.com/donohutcheon/gowebserver/controllers/response/types"
	"net/http"
)

type Response map[string]interface{}

func New(status bool, message string) Response {
	m := make(Response)
	m["status"] = status
	m["message"] = message
	return m
}

func NewWithFieldsList(status bool, message string, fields []types.ErrorField) Response {
	m := make(Response)
	m["status"] = status
	m["message"] = message
	if len(fields) > 0 {
		m["fields"] = fields
	}
	return m
}

func (m Response) SetResponse(status bool, message string) {
	m["status"] = status
	m["message"] = message
}

func (m Response) Set(key string, value interface{}) {
	m[key] = value
}

func (m Response) SetString(key string, value string) {
	m[key] = value
}

func (m Response) Respond(w http.ResponseWriter) error {
	w.Header().Add("Content-Type", "application/json")

	bytes, err := json.Marshal(m)
	if err != nil {
		return err
	}

	_, err = w.Write(bytes)
	if err != nil {
		return err
	}

	return nil
}

