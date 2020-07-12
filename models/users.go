package models

import (
	"database/sql"
	"fmt"
	e "github.com/donohutcheon/gowebserver/controllers/errors"
	"github.com/donohutcheon/gowebserver/controllers/response/types"
	"github.com/donohutcheon/gowebserver/datalayer"
	"github.com/donohutcheon/gowebserver/router/auth"
	"github.com/donohutcheon/gowebserver/state"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strings"
)

type Settings struct {
	ID int `json:"id"`
	ThemeName string `json:"themeName"`
}

type User struct {
	datalayer.Model
	serverState  *state.ServerState
	Email        string    `json:"email"`
	FirstName    string    `json:"firstName"`
	Surname      string    `json:"surname"`
	Age          int       `json:"age"`
	Address      string    `json:"address"`
	Roles        []string  `json:"roles"`
	Settings     Settings  `json:"settings"`
	Password     string    `json:"password,omitempty"`
	/*AccessToken  string    `json:"accessToken,omitempty" sql:"-"`
	RefreshToken string    `json:"refreshToken,omitempty" sql:"-"`
	LoggedOutAt  time.Time `json:"loggedOutAt,omitempty"`*/
}

func NewUser(state *state.ServerState) *User {
	user := new(User)
	user.serverState = state
	return user
}

func (u *User) convert(user datalayer.User) {
	u.ID = user.ID
	u.CreatedAt = user.CreatedAt
	u.UpdatedAt = user.UpdatedAt
	u.DeletedAt = user.DeletedAt
	if user.Email.Valid {
		u.Email = user.Email.String
	}
	if user.Password.Valid {
		u.Password = user.Password.String
	}
	u.Roles = []string{"ADMIN","USER"}
	u.Settings.ID = 0
	u.Settings.ThemeName = "default"
}

//Validate incoming user details...
func (u *User) validate() error {
	if !strings.Contains(u.Email, "@") {
		return ErrValidationEmail
	}

	if len(u.Password) < 6 {
		return ErrValidationPassword
	}

	//Email must be unique
	//check for errors and duplicate emails
	dl := u.serverState.DataLayer
	_, err := dl.GetUserByEmail(u.Email)
	if err != sql.ErrNoRows {
		return ErrEmailExists
	}

	return nil
}

func (u *User) Create() (*User, error) {
	logger := u.serverState.Logger
	err := u.validate()
	if err != nil {
		return nil, err
	}

	// TODO: Add some sort of salt
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	u.Password = string(hashedPassword)

	dl := u.serverState.DataLayer
	// TODO: Include roles
	id, err :=  dl.CreateUser(u.Email, u.Password)
	if err != nil {
		logger.Fatal(err) // TODO: remove
		return nil, err
	}

	dbUser, err := dl.GetUserByID(id)
	if err != nil && err != datalayer.ErrNoData {
		return nil, err
	}

	// Send confirmation Email
	u.serverState.Channels.ConfirmUsers <- *dbUser

	user := new(User)
	user.convert(*dbUser)

	user.Password = "" //delete password

	return user, nil
}

func (u *User) Login(email, password string) (*auth.TokenResponse, error) {
	dataLayer := u.serverState.DataLayer
	dbUser, err := dataLayer.GetUserByEmail(email)
	if err == sql.ErrNoRows {
		return nil, ErrLoginFailed
	} else if err != nil {
		return nil, err
	}

	if datalayer.UserState(dbUser.State.String) != datalayer.UserStateConfirmed {
		return nil, ErrUserNotConfirmed
	}

	u.convert(*dbUser)

	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword { //Password does not match!
		return nil, ErrLoginFailed
	}
	// Worked! Logged In
	u.Password = ""

	// Create JWT token
	tokenResp, err := auth.CreateToken(u.ID)
	if err != nil {
		return nil, e.Wrap("token creation failed", http.StatusInternalServerError, err)
	}

	return tokenResp, nil
}

func (u *User) GetUser(id int64) (error) {
	dl := u.serverState.DataLayer
	dbUser, err := dl.GetUserByID(id)
	if err == datalayer.ErrNoData {
		return ErrUserDoesNotExist
	} else if err != nil {
		return e.Wrap(fmt.Sprintf("Failed to query user [%d] from database", id), http.StatusInternalServerError, err)
	}

	u.convert(*dbUser)

	return nil
}

func (u *User) ConfirmUser(nonce string) error {
	logger := u.serverState.Logger
	dl := u.serverState.DataLayer

	signUp, err := dl.LookupSignUpConfirmation(nonce)
	if err != nil {
		return e.NewError("User confirmation not found", []types.ErrorField{
			{Name: "nonce", Message: "Nonce not found"},
		}, http.StatusNotFound)
	}
	logger.Printf("Received nonce confirmation for user %d %s", signUp.UserID, nonce)

	err = dl.SetUserStateByID(signUp.UserID, datalayer.UserStateConfirmed)
	if err != nil {
		return e.Wrap(fmt.Sprintf("Failed to confirm user [%d]", signUp.UserID), http.StatusInternalServerError, err)
	}

	return nil
}

func (u *User) GetAPIToken() (*auth.APITokenResponse, error) {
	tokenResp, err := auth.CreateAPIToken(u.ID)
	if err != nil {
		return nil, e.Wrap("api token creation failed", http.StatusInternalServerError, err)
	}

	return tokenResp, nil
}