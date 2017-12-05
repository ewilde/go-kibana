package kibana

import (
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"testing"
)

const defaultElkVersion = "6.0.0"

var authForContainerVersion = map[string]AuthenticationHandler{
	"5.5.3": &BasicAuthenticationHandler{"elastic", "changeme"},
	"6.0.0": &NoAuthenticationHandler{},
}

func Test_NewClient(t *testing.T) {
	kibanaClient := defaultTestKibanaClient()

	assert.NotNil(t, kibanaClient)
	assert.Equal(t, os.Getenv(EnvKibanaUri), kibanaClient.Config.HostAddress)
}

func TestMain(m *testing.M) {
	testContext, err := startKibana(GetEnvVarOrDefault("ELK_VERSION", defaultElkVersion), defaultTestKibanaClient())
	if err != nil {
		log.Fatalf("Could start kibana: %v", err)
	}

	err = os.Setenv(EnvKibanaUri, testContext.KibanaUri)
	if err != nil {
		log.Fatalf("Could not set kibana uri env variable: %v", err)
	}

	err = os.Setenv(EnvKibanaIndexId, testContext.KibanaIndexId)
	if err != nil {
		log.Fatalf("Could not set kibana index id env variable: %v", err)
	}

	code := m.Run()

	stopKibana(testContext)

	os.Exit(code)

}

func defaultTestKibanaClient() *KibanaClient {
	kibanaClient := NewClient(NewDefaultConfig())
	kibanaClient.SetAuth(authForContainerVersion[kibanaClient.Config.KibanaVersion])
	return kibanaClient
}
