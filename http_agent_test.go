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
		ClientId: os.Getenv(EnvLogzClientId),
		UserName: os.Getenv(EnvKibanaUserName),
		Password: os.Getenv(EnvKibanaPassword),
	}
}

func testPreCheckForLogz(t *testing.T) {
	if v := os.Getenv("EnvLogzClientId"); v == "" {
		t.Fatal("EnvLogzClientId must be set for this test")
	}
	if v := os.Getenv("KIBANA_USERNAME"); v == "" {
		t.Fatal("KIBANA_USERNAME must be set for this test")
	}
	if v := os.Getenv("KIBANA_PASSWORD"); v == "" {
		t.Fatal("KIBANA_PASSWORD must be set for this test")
	}
}
