package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMainHandlerSuccess(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)
	req := httptest.NewRequest("GET", "/cafe?count=2&city=moscow", nil)
	responseRecorder := httptest.NewRecorder()
	handler.ServeHTTP(responseRecorder, req)

	require.Equal(t, http.StatusOK, responseRecorder.Code, "Expected status code 200")
	assert.NotEmpty(t, responseRecorder.Body.String(), "Response body should not be empty")
}

func TestMainHandlerUnsupportedCity(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)
	req := httptest.NewRequest("GET", "/cafe?count=2&city=unknowncity", nil)
	responseRecorder := httptest.NewRecorder()
	handler.ServeHTTP(responseRecorder, req)

	require.Equal(t, http.StatusBadRequest, responseRecorder.Code, "Expected status code 400")
	assert.Equal(t, "wrong city value", responseRecorder.Body.String(), "Expected error message 'wrong city value'")
}

func TestMainHandlerWhenCountMoreThanTotal(t *testing.T) {
	totalCount := 4
	handler := http.HandlerFunc(mainHandle)
	req := httptest.NewRequest("GET", "/cafe?count=10&city=moscow", nil)
	responseRecorder := httptest.NewRecorder()
	handler.ServeHTTP(responseRecorder, req)

	require.Equal(t, http.StatusOK, responseRecorder.Code, "Expected status code 200")

	body := responseRecorder.Body.String()
	list := strings.Split(body, ",")

	assert.Len(t, list, totalCount, "Expected cafe count %d, got %d", totalCount, len(list))
}
