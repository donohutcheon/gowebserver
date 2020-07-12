package models

import (
	e "github.com/donohutcheon/gowebserver/controllers/errors"
	"github.com/donohutcheon/gowebserver/controllers/response/types"
	"github.com/donohutcheon/gowebserver/datalayer"
	"github.com/donohutcheon/gowebserver/state"
	"net/http"
)

type SignUpConfirmation struct {
	datalayer.Model
	serverState  *state.ServerState
	Nonce	string
	UserID  int64
}

func NewSignUpConfirmation(state *state.ServerState) *SignUpConfirmation {
	signUp := new(SignUpConfirmation)
	signUp.serverState = state
	return signUp
}

func (s *SignUpConfirmation) convert(signUp datalayer.SignUpConfirmation) {
	s.ID = signUp.ID
	s.CreatedAt = signUp.CreatedAt
	s.UpdatedAt = signUp.UpdatedAt
	s.DeletedAt = signUp.DeletedAt
	s.UserID = signUp.UserID
	s.Nonce = signUp.Nonce
}

func (s *SignUpConfirmation) LookupUsingNonce(nonce string) error {
	dl := s.serverState.DataLayer

	dbSignUp, err := dl.LookupSignUpConfirmation(nonce)
	if err != nil {
		return e.NewError("User confirmation not found", []types.ErrorField{
			{Name: "nonce", Message: "Nonce not found"},
		}, http.StatusNotFound)
	}

	s.convert(*dbSignUp)
	return nil
}