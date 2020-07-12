package controllers_test

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
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

type CreateCardTransactionControllerResponse struct {
	Message         string                 `json:"message"`
	Status          bool                   `json:"status"`
	CardTransaction models.CardTransaction `json:"cardTransaction"`
}

type GetCardTransactionControllerResponse struct {
	Message          string                   `json:"message"`
	Status           bool                     `json:"status"`
	CardTransactions []models.CardTransaction `json:"cardTransactions"`
}

type CreateCardTransactionParameters struct {
	skip          bool
	request       models.CardTransaction
	expResponse   CreateCardTransactionControllerResponse
	expHTTPStatus int
}

type GetCardTransactionParameters struct {
	skip          bool
	expResponse   GetCardTransactionControllerResponse
	expHTTPStatus int
}

func TestCardTransactions(t *testing.T) {
	testTime := time.Now()

	tests := []struct {
		name                        string
		authParameters              AuthParameters
		createCardTransactionParams CreateCardTransactionParameters
		getCardTransactionParams    GetCardTransactionParameters
	}{
		{
			name: "Golden",
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
			createCardTransactionParams: CreateCardTransactionParameters{
				request: models.CardTransaction{
					Model:    datalayer.Model{},
					DateTime: time.Date(2020, 04, 25, 19, 46, 23, 33, time.UTC),
					Amount: models.CurrencyValue{
						Value: 400,
						Scale: 2,
					},
					CurrencyCode:         "ZAR",
					Reference:            "simulation",
					MerchantName:         "Dwelms en Dinges",
					MerchantCity:         "Hillbrow",
					MerchantCountryCode:  "ZA",
					MerchantCountryName:  "South Africa",
					MerchantCategoryCode: "contraband",
					MerchantCategoryName: "Contraband",
					UserID:               1,
				},
				expResponse: CreateCardTransactionControllerResponse{
					Message: "success",
					Status:  true,
					CardTransaction: models.CardTransaction{
						Model: datalayer.Model{
							ID: 0,
							CreatedAt: datalayer.JsonNullTime{
								NullTime: sql.NullTime{
									Time:  testTime,
									Valid: true,
								},
							},
						},
						DateTime: time.Date(2020, 04, 25, 19, 46, 23, 33, time.UTC),
						Amount: models.CurrencyValue{
							Value: 400,
							Scale: 2,
						},
						CurrencyCode:         "ZAR",
						Reference:            "simulation",
						MerchantName:         "Dwelms en Dinges",
						MerchantCity:         "Hillbrow",
						MerchantCountryCode:  "ZA",
						MerchantCountryName:  "South Africa",
						MerchantCategoryCode: "contraband",
						MerchantCategoryName: "Contraband",
					},
				},
				expHTTPStatus: http.StatusOK,
			},
			getCardTransactionParams: GetCardTransactionParameters{
				expResponse: GetCardTransactionControllerResponse{
					Message: "success",
					Status:  true,
					CardTransactions: []models.CardTransaction{
						{
							Model: datalayer.Model{
								ID: 0,
								CreatedAt: datalayer.JsonNullTime{
									NullTime: sql.NullTime{
										Time:  testTime,
										Valid: true,
									},
								},
							},
							DateTime: time.Date(2020, 04, 25, 19, 46, 23, 33, time.UTC),
							Amount: models.CurrencyValue{
								Value: 400,
								Scale: 2,
							},
							CurrencyCode:         "ZAR",
							Reference:            "simulation",
							MerchantName:         "Dwelms en Dinges",
							MerchantCity:         "Hillbrow",
							MerchantCountryCode:  "ZA",
							MerchantCountryName:  "South Africa",
							MerchantCategoryCode: "contraband",
							MerchantCategoryName: "Contraband",
						},
					},
				},
				expHTTPStatus: http.StatusOK,
			},
		},
		{
			name: "No data",
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
			createCardTransactionParams: CreateCardTransactionParameters{
				skip: true,
			},
			getCardTransactionParams: GetCardTransactionParameters{
				expResponse: GetCardTransactionControllerResponse{
					Message:          "success",
					Status:           true,
					CardTransactions: []models.CardTransaction{},
				},
				expHTTPStatus: http.StatusOK,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cl := new(http.Client)

			callbacks := state.NewMockCallbacks(mailCallback)

			state := facotory.NewForTesting(t, callbacks)
			ctx := state.Context

			gotAuthResp := login(t, ctx, cl, state.URL, test.authParameters)
			createCardTransaction(t, ctx, cl, state.URL, gotAuthResp, &test.createCardTransactionParams)
			getCardTransactions(t, ctx, cl, state.URL, gotAuthResp, &test.getCardTransactionParams)
		})
	}
}

func createCardTransaction(t *testing.T, ctx context.Context, cl *http.Client,
	url string, auth *AuthResponse, params *CreateCardTransactionParameters) {
	if params.skip {
		return
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url+"/api/card-transactions/new", nil)
	assert.NoError(t, err)
	req.Header.Add("Authorization", "Bearer "+auth.Token.AccessToken)

	b, err := json.Marshal(params.request)
	require.NoError(t, err)

	req.Body = ioutil.NopCloser(bytes.NewReader(b))
	defer req.Body.Close()

	res, err := cl.Do(req)
	require.NoError(t, err)

	body, err := ioutil.ReadAll(res.Body)
	gotResp := new(CreateCardTransactionControllerResponse)
	err = json.Unmarshal(body, gotResp)
	require.NoError(t, err)

	assert.Equal(t, params.expResponse.Status, gotResp.Status)
	assert.Equal(t, params.expResponse.Message, gotResp.Message)
	assert.Equal(t, params.expResponse.CardTransaction.Amount.Value, gotResp.CardTransaction.Amount.Value)
}

func getCardTransactions(t *testing.T, ctx context.Context, cl *http.Client,
	url string, auth *AuthResponse, params *GetCardTransactionParameters) {
	if params.skip {
		return
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url+"/api/me/card-transactions", nil)
	assert.NoError(t, err)
	req.Header.Add("Authorization", "Bearer "+auth.Token.AccessToken)

	res, err := cl.Do(req)
	require.NoError(t, err)

	body, err := ioutil.ReadAll(res.Body)
	gotResp := new(GetCardTransactionControllerResponse)
	err = json.Unmarshal(body, gotResp)
	require.NoError(t, err)

	assert.Equal(t, params.expResponse.Status, gotResp.Status)
	assert.Equal(t, params.expResponse.Message, gotResp.Message)
	require.Equal(t, len(params.expResponse.CardTransactions), len(gotResp.CardTransactions))
	for i, x := range params.expResponse.CardTransactions {
		// Negate datetime fields we can't control.
		x.CreatedAt = datalayer.JsonNullTime{}
		gotResp.CardTransactions[i].CreatedAt = datalayer.JsonNullTime{}

		assert.Equal(t, x, gotResp.CardTransactions[i])
	}
}