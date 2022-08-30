package controller

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

type response struct {
	Status string
	Data   string
}

func checkBody(t *testing.T, recorder *httptest.ResponseRecorder, e error) {
	var resp response
	json.Unmarshal(recorder.Body.Bytes(), &resp)
	require.Equal(t, e.Error(), resp.Data)
}
