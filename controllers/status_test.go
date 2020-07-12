package controllers_test

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/donohutcheon/gowebserver/controllers/response"
	"github.com/donohutcheon/gowebserver/state"
	"github.com/donohutcheon/gowebserver/state/facotory"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetStatus(t *testing.T) {
	state := facotory.NewForTesting(t, state.NewMockCallbacks(mailCallback))
	ctx := state.Context

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, state.URL + "/api/status", nil)
	assert.NoError(t, err)

	cl := new(http.Client)
	res, err := cl.Do(req)
	require.NoError(t, err)

	body, err := ioutil.ReadAll(res.Body)

	expResp := response.New(true, "Service is up")
	expected, err := json.Marshal(expResp)
	require.NoError(t, err)

	assert.Equal(t, string(expected), string(body))
}
