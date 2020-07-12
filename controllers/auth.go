package controllers

import (
	"encoding/json"
	"github.com/donohutcheon/gowebserver/controllers/errors"
	"github.com/donohutcheon/gowebserver/controllers/response"
	"github.com/donohutcheon/gowebserver/models"
	"github.com/donohutcheon/gowebserver/router/auth"
	"github.com/donohutcheon/gowebserver/state"
	"net/http"
)

func Authenticate(w http.ResponseWriter, r *http.Request, state *state.ServerState) error {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("X-FRAME-OPTIONS", "SAMEORIGIN")
	w.Header().Set("X-XSS-Protection", "1; mode=block")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	//w.Header().Set("Sec-Fetch-Site", "same-site")

	if r.Method == http.MethodOptions {
		return nil
	}

	user := models.NewUser(state)
	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		err = errors.Wrap("Invalid request format", http.StatusBadRequest, err)
		errors.WriteError(w, err, http.StatusBadRequest)
		return err
	}

	data, err := user.Login(user.Email, user.Password)
	if err != nil {
		errors.WriteError(w, err)
		return err
	}

	resp := response.New(true, "Logged In")
	resp["token"] = data
	resp.Respond(w)

	return nil
}

func RefreshToken(w http.ResponseWriter, r *http.Request, state *state.ServerState) error {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("X-FRAME-OPTIONS", "SAMEORIGIN")
	w.Header().Set("X-XSS-Protection", "1; mode=block")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	//w.Header().Set("Sec-Fetch-Site", "same-site")

	if r.Method == http.MethodOptions {
		return nil
	}

	refreshTokenReq := new(auth.RefreshJWTReq)
	err := json.NewDecoder(r.Body).Decode(refreshTokenReq) //decode the request body into struct and failed if any error occur
	if err != nil {
		errors.WriteError(w, err, http.StatusBadRequest)
		return err
	}

	if refreshTokenReq.GrantType != "refresh_token" {
		errors.WriteError(w, errors.NewError("grant type not refresh_token", nil, http.StatusBadRequest))
		return err
	}

	data, err := auth.RefreshToken(refreshTokenReq.RefreshToken)
	if err != nil {
		errors.WriteError(w, err)
		return err
	}

	resp := response.New(true, "Tokens refreshed")
	resp.Set("token", data)
	resp.Respond(w)

	return nil
}

func GetAPIToken(w http.ResponseWriter, r *http.Request, state *state.ServerState) error {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("X-FRAME-OPTIONS", "SAMEORIGIN")
	w.Header().Set("X-XSS-Protection", "1; mode=block")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	//w.Header().Set("Sec-Fetch-Site", "same-site")

	if r.Method == http.MethodOptions {
		return nil
	}

	userID := r.Context().Value("userID").(int64)

	user := models.NewUser(state)
	user.ID = userID
	data, err := user.GetAPIToken()
	if err != nil {
		errors.WriteError(w, errors.NewError("could not generate API token", nil, http.StatusForbidden))
	}

	resp := response.New(true, "Logged In")
	resp["token"] = data
	resp.Respond(w)

	return nil
}