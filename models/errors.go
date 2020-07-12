package models

import (
	"github.com/donohutcheon/gowebserver/controllers/response/types"
	"net/http"

	e "github.com/donohutcheon/gowebserver/controllers/errors"
)

var (
	ErrLoginFailed = e.NewError("Invalid login credentials", []types.ErrorField{
		{Name: "email", Message: "Invalid login credentials"},
		{Name: "password", Message: "Invalid login credentials"},
	},
		http.StatusForbidden)

	ErrUserNotConfirmed = e.NewError("User has not confirmed their email address", []types.ErrorField{},
		http.StatusForbidden)

	ErrValidationEmail = e.NewError("Email address is required", []types.ErrorField{
		{Name: "email", Message: "A valid email address is required"},
	},
		http.StatusBadRequest)

	ErrValidationPassword = e.NewError("Password is required", []types.ErrorField{
		{Name: "password", Message: "Password must be longer than 6 characters"},
	}, http.StatusBadRequest)

	ErrUserDoesNotExist = e.NewError("User does not exist", nil, http.StatusForbidden)

	ErrEmailExists      = e.NewError("Email address already exists", []types.ErrorField{
		{Name: "email", Message: "Email address already exists"},
	}, http.StatusBadRequest)

	ErrValidationFailed = e.NewError("Invalid request, validation failed", nil, http.StatusBadRequest)

	ErrValidationName = e.NewError("Contact name is required", []types.ErrorField{
		{Name: "name", Message: "Contact name is required"},
	}, http.StatusBadRequest)

	ErrValidationPhone = e.NewError("Contact phone number is required", []types.ErrorField{
		{Name: "phone", Message: "Contact phone number is required"},
	}, http.StatusBadRequest)
)
