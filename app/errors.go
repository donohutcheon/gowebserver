package app

import (
	"net/http"

	"github.com/donohutcheon/gowebserver/controllers/response"
)

var NotFoundHandler = func(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		resp := response.New(false, "The resource was not found on our server")
		resp.Respond(w)
		next.ServeHTTP(w, r)
	})
}
