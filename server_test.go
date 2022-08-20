package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetupRoutes(t *testing.T) {
	SetupRoutes()

	w := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/getUserData", nil)

	getUserData(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "{\"Email\":\"toto@gmail.com\",\"Name\":\"Frank\"}", w.Body.String())
}
