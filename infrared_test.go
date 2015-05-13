package main

import (
	"./lib"
	"./lib/config"
	"bytes"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var serverConfig *config.Server = config.Server{":8080", {"localhost", "infrared_test"}}

func TestConfigShouldReturn400ForBadNodeType(t *testing.T) {
	config := strings.NewReader("{ \"test\": true }")
	req, _ := http.NewRequest("PUT", "/api/v1//config", config)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	infrared.Setup(serverConfig).ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code, "Expected a 400 status code")
}

func TestConfigShouldAllowCreation(t *testing.T) {
	config := strings.NewReader("{ \"test\": true }")
	req, _ := http.NewRequest("PUT", "/api/v1/test/config", config)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	infrared.Setup(serverConfig).ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code, "Expected a 200 status code")
}

func TestConfigShouldReturn404IfNoEntry(t *testing.T) {
	req, _ := http.NewRequest("GET", "/api/v1/notfound/config", nil)
	w := httptest.NewRecorder()
	infrared.Setup(&serverConfig).ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code, "Expected a 404 status code")
}

func TestConfigShouldReturn200ForValidEntry(t *testing.T) {
	req, _ := http.NewRequest("GET", "/api/v1/test/config", nil)
	w := httptest.NewRecorder()
	infrared.Setup(serverConfig).ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code, "Expected a 200 status code")
}

func TestConfigShouldReturnCorrectConfiguration(t *testing.T) {
	req, _ := http.NewRequest("GET", "/api/v1/test/config", nil)
	w := httptest.NewRecorder()
	infrared.Setup(serverConfig).ServeHTTP(w, req)

	bodyData := new(bytes.Buffer)
	bodyData.ReadFrom(w.Body)

	assert.Equal(t, "{\n  \"test\": true\n}", bodyData.String(), "Expected returned configuration to match set one.")
}
