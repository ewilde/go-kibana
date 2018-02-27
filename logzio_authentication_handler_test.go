package kibana

import (
	"github.com/parnurzeal/gorequest"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func Test_LogzAuthentication_handler(t *testing.T) {
	if skip := testPreCheckForLogz(t); skip {
		t.Skip()
		return
	}

	handler := createLogzAuthenticationHandler()

	handler.Initialize(gorequest.New())

	assert.NotEmpty(t, handler.sessionToken, "Session token should not be empty")
}


func testPreCheckForLogz(t *testing.T) bool {
	config := NewDefaultConfig()
	if config.KibanaType == KibanaTypeVanilla {
		return true
	}

	if v := os.Getenv(EnvLogzClientId); v == "" {
		t.Fatalf("%s must be set for this test", EnvLogzClientId)
	}
	if v := os.Getenv(EnvKibanaUserName); v == "" {
		t.Fatalf("%s must be set for this test", EnvKibanaUserName)
	}
	if v := os.Getenv(EnvKibanaPassword); v == "" {
		t.Fatalf("%s must be set for this test", EnvKibanaPassword)
	}

	return false
}
