package controllers_test

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/donohutcheon/gowebserver/datalayer"
	"github.com/donohutcheon/gowebserver/models"
	"github.com/donohutcheon/gowebserver/state"
	"github.com/donohutcheon/gowebserver/state/facotory"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type UserControllerResponse struct {
	Message string      `json:"message"`
	Status  bool        `json:"status"`
	User    models.User `json:"user"`
}

type GetCurrentUserParameters struct {
	skip          bool
	expResponse   UserControllerResponse
	expHTTPStatus int
}

type CreateUserParameters struct {
	skip          bool
	createUserReq models.User
	expResponse   UserControllerResponse
	expHTTPStatus int
}

func getCurrentUser(t *testing.T, ctx context.Context, cl *http.Client,
	url string, auth *AuthResponse, params *GetCurrentUserParameters) *UserControllerResponse {
	if params.skip {
		return nil
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url + "/api/users/current", nil)
	assert.NoError(t, err)
	req.Header.Add("Authorization", "Bearer " + auth.Token.AccessToken)
	res, err := cl.Do(req)
	require.NoError(t, err)

	body, err := ioutil.ReadAll(res.Body)
	gotResp := new(UserControllerResponse)
	err = json.Unmarshal(body, gotResp)
	require.NoError(t, err)

	assert.Equal(t, params.expHTTPStatus, res.StatusCode)
	assert.Equal(t, params.expResponse.Status, gotResp.Status)
	assert.Equal(t, params.expResponse.Message, gotResp.Message)
	assert.Equal(t, params.expResponse.User.Email, gotResp.User.Email)

	return gotResp
}

func createUser(t *testing.T, ctx context.Context, cl *http.Client,
	url string, params *CreateUserParameters) *UserControllerResponse {
	if params.skip {
		return nil
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url + "/api/auth/sign-up", nil)

	assert.NoError(t, err)

	b, err := json.Marshal(params.createUserReq)
	require.NoError(t, err)

	req.Body = ioutil.NopCloser(bytes.NewReader(b))
	defer req.Body.Close()

	res, err := cl.Do(req)
	require.NoError(t, err)

	body, err := ioutil.ReadAll(res.Body)
	gotCreateUserResp := new(UserControllerResponse)
	err = json.Unmarshal(body, gotCreateUserResp)
	require.NoError(t, err)
	assert.Equal(t, params.expHTTPStatus, res.StatusCode)
	assert.Equal(t, params.expResponse.Message, gotCreateUserResp.Message)
	assert.Equal(t, params.expResponse.Status, gotCreateUserResp.Status)
	assert.Equal(t, params.expResponse.User.Email, gotCreateUserResp.User.Email)
	return gotCreateUserResp
}

func TestGetCurrentUser(t *testing.T) {
	tests := []struct {
		name                 string
		authParameters       AuthParameters
		getCurrentUserParams GetCurrentUserParameters
	}{
		{
			name: "Success",
			authParameters: AuthParameters{
				authRequest: models.User{
					Email:    "subzero@dreamrealm.com",
					Password: "secret",
				},
				expHTTPStatus: http.StatusOK,
				expLoginResp: AuthResponse{
					Message: "Logged In",
					Status:  true,
				},
			},
			getCurrentUserParams: GetCurrentUserParameters{
				expResponse: UserControllerResponse{
					Message: "success",
					Status:  true,
					User: models.User{
						Model: datalayer.Model{
							ID: 0,
							CreatedAt: datalayer.JsonNullTime{
								NullTime: sql.NullTime{
									Time:  time.Now(),
									Valid: true,
								},
							},
						},
						Email:    "subzero@dreamrealm.com",
						Password: "$2a$10$NkTUeL6hkTRZ7M13tKYLqOmg7pAQaGPdpch9b5UoTSoO77MHjbPjm",
						Roles:    []string{"ADMIN", "USER"},
						Settings: models.Settings{
							ID:        0,
							ThemeName: "default",
						},
					},
				},
				expHTTPStatus: http.StatusOK,
			},
		},
	}

	callbacks := state.NewMockCallbacks(mailCallback)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cl := new(http.Client)
			//now := time.Now()
			state := facotory.NewForTesting(t, callbacks)
			ctx := state.Context
			gotAuthResp := login(t, ctx, cl, state.URL, test.authParameters)
			getCurrentUser(t, ctx, cl, state.URL, gotAuthResp, &test.getCurrentUserParams)
		})
	}
}

func TestCreateUser(t *testing.T) {
	tests := []struct {
		name              string
		authParameters    AuthParameters
		createUserParams  CreateUserParameters
	}{
		{
			name: "Golden",
			authParameters: AuthParameters{
				authRequest: models.User{
					Email:    "jade@edenia.com",
					Password: "secret",
				},
				expHTTPStatus: http.StatusOK,
				expLoginResp: AuthResponse{
					Message: "Logged In",
					Status:  true,
				},
			},
			createUserParams: CreateUserParameters{
				createUserReq: models.User{
					Email:        "jade@edenia.com",
					Password:     "secret",
				},
				expResponse: UserControllerResponse{
					Message: "User has been created",
					Status:  true,
					User: models.User{
						Email: "jade@edenia.com",
					},
				},
				expHTTPStatus: http.StatusOK,
			},
		},
		{
			name: "Incomplete Email",
			createUserParams: CreateUserParameters{
				createUserReq: models.User{
					Email:        "",
					Password:     "secret",
				},
				expResponse: UserControllerResponse{
					Message: "Email address is required",
					Status:  false,
				},
				expHTTPStatus: http.StatusBadRequest,
			},
		},
		{
			name: "Incomplete Password",
			createUserParams: CreateUserParameters{
				createUserReq: models.User{
					Email:        "sindel@dreamrealm.com",
				},
				expResponse: UserControllerResponse{
					Message: "Password is required",
					Status:  false,
				},
				expHTTPStatus: http.StatusBadRequest,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			callbacks := state.NewMockCallbacks(mailCallback)
			state := facotory.NewForTesting(t, callbacks)
			ctx := state.Context
			cl := new(http.Client)

			createUser(t, ctx, cl, state.URL, &test.createUserParams)

			// End the test if createUser is supposed to fail.
			if !test.createUserParams.expResponse.Status {
				return
			}

			callbacks.MockMailWG.Wait()

			fmt.Println("Got mail, lets login...")

			login(t, ctx, cl, state.URL, test.authParameters)
		})
	}
}