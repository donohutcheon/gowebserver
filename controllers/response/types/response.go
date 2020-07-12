package types

type Response map[string]interface{}

type ErrorField struct {
	Name    string `json:"name"`
	Message string `json:"message"`
	Direct  bool   `json:"direct"`
}
