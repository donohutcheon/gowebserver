package routes

import (
	"net/http"

	"github.com/donohutcheon/gowebserver/controllers"
	"github.com/donohutcheon/gowebserver/state"
)

type MiddlewareFunc func(next http.Handler, state *state.ServerState, registry map[string]RouteEntry) http.Handler
type HandlerFunc func(w http.ResponseWriter, r *http.Request, handlerState *state.ServerState) error

type RouteEntry struct {
	Handler HandlerFunc
	Methods []string
	Public bool
}

func GetRouteRegistry() map[string]RouteEntry {
	return map[string]RouteEntry {
		"/" : {
			Public: true,
		},
		"/api/status" : {
			Handler: controllers.Status,
			Methods: []string{http.MethodGet},
			Public:  true,
		},
		"/api/auth/sign-up" : {
			Handler: controllers.CreateUser,
			Methods: []string{http.MethodPost, http.MethodOptions},
			Public:  true,
		},
		"/api/users/current" : {
			Handler: controllers.GetCurrentUser,
			Methods: []string{http.MethodGet, http.MethodOptions},
		},
		"/api/auth/login" : {
			Handler: controllers.Authenticate,
			Methods: []string{http.MethodPost, http.MethodOptions},
			Public:  true,
		},
		"/api/auth/api-token" : {
			Handler: controllers.GetAPIToken,
			Methods: []string{http.MethodGet, http.MethodOptions},
		},
		"/api/auth/refresh" : {
			Handler: controllers.RefreshToken,
			Methods: []string{http.MethodPost, http.MethodOptions},
			Public:  true,
		},
		"/api/card-transactions/new" : {
			Handler: controllers.CreateCardTransaction,
			Methods: []string{http.MethodPost, http.MethodOptions},
		},
		"/api/me/card-transactions" : {
			Handler: controllers.GetCardTransactions,
			Methods: []string{http.MethodGet, http.MethodOptions},
		},
		"/api/users/confirm/{nonce}" : {
			Handler: controllers.ConfirmUserSignUp,
			Methods: []string{http.MethodGet, http.MethodOptions},
			Public: true,
		},
	}
}