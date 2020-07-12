package controllers_test

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func mailCallback(t *testing.T, ctx context.Context, to []string, from, subject, message string) {
	fmt.Println("mail callback: ", message)
	re := regexp.MustCompile("https?://.*/users/confirm/[a-z0-9]+")
	url := re.FindString(message)
	fmt.Println("confirmation url is: ", url)

	req, err := http.NewRequestWithContext( ctx, http.MethodGet, url, nil)
	assert.NoError(t, err)

	cl := new(http.Client)
	res, err := cl.Do(req)
	require.NoError(t, err)

	body, err := ioutil.ReadAll(res.Body)
	fmt.Println("Confirmation response body: ", string(body))
}