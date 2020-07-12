package controllers

import (
	"encoding/json"
	"github.com/donohutcheon/gowebserver/controllers/response/types"
	"github.com/gorilla/mux"
	"net/http"

	e "github.com/donohutcheon/gowebserver/controllers/errors"
	"github.com/donohutcheon/gowebserver/controllers/response"
	"github.com/donohutcheon/gowebserver/models"
	"github.com/donohutcheon/gowebserver/state"
)

func CreateUser(w http.ResponseWriter, r *http.Request, state *state.ServerState) error {
	if r.Method == http.MethodOptions {
		return nil
	}

	user := models.NewUser(state)
	err := json.NewDecoder(r.Body).Decode(user) //decode the request body into struct and failed if any error occur
	if err != nil {
		err = e.Wrap("Invalid request", http.StatusBadRequest, err)
		e.WriteError(w, err)
		return err
	}

	data, err := user.Create()
	if err != nil {
		e.WriteError(w, err)
		return err
	}

	resp := response.New(true, "User has been created")
	resp.Set("user", data)
	resp.Respond(w)

	return nil
}

// TODO: Move into usersController
func GetCurrentUser(w http.ResponseWriter, r *http.Request, state *state.ServerState) error {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("X-FRAME-OPTIONS", "SAMEORIGIN")
	w.Header().Set("X-XSS-Protection", "1; mode=block")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set( "Access-Control-Allow-Methods", "OPTIONS,GET,POST,PUT,DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "authorization")

	if r.Method == http.MethodOptions {
		return nil
	}
	id := r.Context().Value("userID").(int64)

	user := models.NewUser(state)
	err := user.GetUser(id)
	if err != nil {
		e.WriteError(w, err)
		return err
	}

	user.Password = ""

	resp := response.New(true, "success")
	resp.Set("user", user)

	err = resp.Respond(w)
	if err != nil {
		return err
	}

	return nil
}

func ConfirmUserSignUp(w http.ResponseWriter, r *http.Request, state *state.ServerState) error {
	if r.Method == http.MethodOptions {
		return nil
	}
	vars := mux.Vars(r)
	nonce, ok := vars["nonce"]
	if !ok {
		err := e.NewError("Path variable 'nonce' is required", []types.ErrorField{
			{Name: "nonce", Message: "Path variable 'nonce' is required"},
		}, http.StatusBadRequest)
		e.WriteError(w, err)
		return err
	}

	signUp := models.NewSignUpConfirmation(state)
	err := signUp.LookupUsingNonce(nonce)
	if err != nil {
		err := e.NewError("No match found for nonce", []types.ErrorField{
			{Name: "nonce", Message: "No match found for nonce"},
		}, http.StatusBadRequest)
		e.WriteError(w, err)
		return err
	}

	user := models.NewUser(state)
	err = user.GetUser(signUp.UserID)
	if err != nil {
		err := e.NewError("User not found", []types.ErrorField{},
			http.StatusInternalServerError)
		e.WriteError(w, err)
		return err
	}

	err = user.ConfirmUser(nonce)
	if err != nil {
		err := e.NewError("Failed to confirm user", []types.ErrorField{},
			http.StatusInternalServerError)
		e.WriteError(w, err)
		return err
	}

	resp := response.New(true, "User's email has been confirmed")
	resp.Respond(w)

	return nil
}