package controllers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/donohutcheon/gowebserver/router/auth"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/donohutcheon/gowebserver/models"
	"github.com/donohutcheon/gowebserver/state"
	"github.com/donohutcheon/gowebserver/state/facotory"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type AuthResponse struct {
	Message string             `json:"message"`
	Status  bool               `json:"status"`
	Token   auth.TokenResponse `json:"token"`
}

func TestAuthenticate(t *testing.T) {
	tests := []struct {
		name          string
		request       []byte
		expResp       AuthResponse
		expStatus     int
		expTokenValid bool
	}{
		{
			name:    "Success",
			request: []byte(`{"email": "subzero@dreamrealm.com", "password": "secret"}`),
			expResp: AuthResponse{
				Message: "Logged In",
				Status:  true,
			},
			expStatus:     http.StatusOK,
			expTokenValid: true,
		},
		{
			name:    "Non-existent User",
			request: []byte(`{"email": "skeletor@eternia.com", "password": "secret"}`),
			expResp: AuthResponse{
				Message: "Invalid login credentials",
				Status:  false,
			},
			expStatus: http.StatusForbidden,
		},
		{
			name:    "Wrong Password",
			request: []byte(`{"email": "subzero@dreamrealm.com", "password": "wrong"}`),
			expResp: AuthResponse{
				Message: "Invalid login credentials",
				Status:  false,
			},
			expStatus: http.StatusForbidden,
		},
		{
			name:    "Garbage Request",
			request: []byte(`garbage`),
			expResp: AuthResponse{
				Message: "Invalid request format",
				Status:  false,
			},
			expStatus: http.StatusBadRequest,
		},
	}

	callbacks := state.NewMockCallbacks(mailCallback)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			now := time.Now()
			state := facotory.NewForTesting(t, callbacks)
			ctx := state.Context

			req, err := http.NewRequestWithContext(ctx, http.MethodPost, state.URL+"/api/auth/login", nil)
			assert.NoError(t, err)

			req.Body = ioutil.NopCloser(bytes.NewReader(test.request))
			defer req.Body.Close()

			cl := new(http.Client)
			res, err := cl.Do(req)
			require.NoError(t, err)

			body, err := ioutil.ReadAll(res.Body)
			gotResp := new(AuthResponse)
			err = json.Unmarshal(body, gotResp)
			require.NoError(t, err)

			assert.Equal(t, test.expStatus, res.StatusCode)
			assert.Equal(t, test.expResp.Message, gotResp.Message)
			assert.Equal(t, test.expResp.Status, gotResp.Status)
			if test.expTokenValid {
				assert.NotEmpty(t, gotResp.Token.AccessToken)
				assert.NotEmpty(t, gotResp.Token.RefreshToken)
				assert.Less(t, now.Unix(), gotResp.Token.ExpiresIn)
			} else {
				assert.Empty(t, gotResp.Token)
			}
		})
	}
}

type AuthParameters struct {
	authRequest   models.User
	expHTTPStatus int
	expLoginResp  AuthResponse
}

type RefreshTokenParameters struct {
	request       auth.RefreshJWTReq
	expHTTPStatus int
	expResponse   AuthResponse
}

func login(t *testing.T, ctx context.Context, cl *http.Client, url string, params AuthParameters) *AuthResponse {
	testTime := time.Now()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url+"/api/auth/login", nil)
	assert.NoError(t, err)

	b, err := json.Marshal(params.authRequest)
	assert.NoError(t, err)

	req.Body = ioutil.NopCloser(bytes.NewReader(b))
	defer req.Body.Close()
	res, err := cl.Do(req)
	require.NoError(t, err)

	body, err := ioutil.ReadAll(res.Body)
	gotAuthResp := new(AuthResponse)
	err = json.Unmarshal(body, gotAuthResp)
	require.NoError(t, err)

	assert.Equal(t, params.expHTTPStatus, res.StatusCode)
	assert.Equal(t, params.expLoginResp.Message, gotAuthResp.Message)
	assert.Equal(t, params.expLoginResp.Status, gotAuthResp.Status)
	if params.expHTTPStatus == http.StatusOK {
		require.NotEmpty(t, gotAuthResp.Token.AccessToken)
		require.NotEmpty(t, gotAuthResp.Token.RefreshToken)
		require.Less(t, testTime.Unix(), gotAuthResp.Token.ExpiresIn)
	} else {
		assert.Empty(t, gotAuthResp.Token)
	}
	return gotAuthResp
}

func refreshToken(t *testing.T, ctx context.Context, cl *http.Client, url string, params RefreshTokenParameters) *AuthResponse {
	testTime := time.Now()
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url+"/api/auth/refresh", nil)
	assert.NoError(t, err)

	b, err := json.Marshal(params.request)
	assert.NoError(t, err)

	req.Body = ioutil.NopCloser(bytes.NewReader(b))
	defer req.Body.Close()
	res, err := cl.Do(req)
	require.NoError(t, err)

	body, err := ioutil.ReadAll(res.Body)
	gotAuthResp := new(AuthResponse)
	err = json.Unmarshal(body, gotAuthResp)
	require.NoError(t, err)

	assert.Equal(t, params.expHTTPStatus, res.StatusCode)
	assert.Equal(t, params.expResponse.Message, gotAuthResp.Message)
	assert.Equal(t, params.expResponse.Status, gotAuthResp.Status)
	if params.expHTTPStatus == http.StatusOK {
		require.NotEmpty(t, gotAuthResp.Token.AccessToken)
		require.NotEmpty(t, gotAuthResp.Token.RefreshToken)
		require.Less(t, testTime.Unix(), gotAuthResp.Token.ExpiresIn)
	} else {
		assert.Empty(t, gotAuthResp.Token)
	}
	return gotAuthResp
}

func TestRefreshToken(t *testing.T) {
	tests := []struct {
		name               string
		authParams         AuthParameters
		refreshTokenParams RefreshTokenParameters
	}{
		{
			name: "Golden",
			authParams: AuthParameters{
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
			refreshTokenParams: RefreshTokenParameters{
				expHTTPStatus: http.StatusOK,
				expResponse: AuthResponse{
					Message: "Tokens refreshed",
					Status:  true,
				},
			},
		},
	}

	ctx := context.Background()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cl := new(http.Client)

			callbacks := state.NewMockCallbacks(mailCallback)
			state := facotory.NewForTesting(t, callbacks)

			gotAuthResp := login(t, ctx, cl, state.URL, test.authParams)
			test.refreshTokenParams.request.RefreshToken = gotAuthResp.Token.RefreshToken
			test.refreshTokenParams.request.GrantType = "refresh_token"
			refreshToken(t, ctx, cl, state.URL, test.refreshTokenParams)
		})
	}
}
