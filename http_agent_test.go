package kibana

import (
	"github.com/parnurzeal/gorequest"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func Test_LogzAuthentication_handler(t *testing.T) {
	testPreCheckForLogz(t)

	handler := createLogzAuthenticationHandler()

	handler.Initialize(gorequest.New())

	assert.NotEmpty(t, handler.sessionToken, "Session token should not be empty")
}

func createLogzAuthenticationHandler() *LogzAuthenticationHandler {
	return &LogzAuthenticationHandler{
		Auth0Uri: "https://logzio.auth0.com",
		LogzUri:  "https://app-eu.logz.io",
		ClientId: os.Getenv("LOGZ_CLIENT_ID"),
		UserName: os.Getenv("LOGZ_USERNAME"),
		Password: os.Getenv("LOGZ_PASSWORD"),
	}
}

func testPreCheckForLogz(t *testing.T) {
	if v := os.Getenv("LOGZ_CLIENT_ID"); v == "" {
		t.Fatal("LOGZ_CLIENT_ID must be set for this test")
	}
	if v := os.Getenv("LOGZ_USERNAME"); v == "" {
		t.Fatal("LOGZ_USERNAME must be set for this test")
	}
	if v := os.Getenv("LOGZ_PASSWORD"); v == "" {
		t.Fatal("LOGZ_PASSWORD must be set for this test")
	}
}
