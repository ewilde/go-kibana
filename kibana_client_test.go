package kibana

import (
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"testing"
)

const defaultElkVersion = "6.0.0"

var authForContainerVersion = map[string]map[KibanaType]AuthenticationHandler{
	"5.5.3": {
		KibanaTypeVanilla: &BasicAuthenticationHandler{"elastic", "changeme"},
		KibanaTypeLogzio:  createLogzAuthenticationHandler(),
	},
	"6.0.0": {KibanaTypeVanilla: &NoAuthenticationHandler{}},
}

func Test_NewClient(t *testing.T) {
	kibanaClient := defaultTestKibanaClient()

	assert.NotNil(t, kibanaClient)
	assert.Equal(t, os.Getenv(EnvKibanaUri), kibanaClient.Config.KibanaBaseUri)
}

func TestMain(m *testing.M) {
	client := defaultTestKibanaClient()

	if client.Config.KibanaType == KibanaTypeVanilla {
		runTestsWithContainers(m, client)
	} else {
		runTestsWithoutContainers(m)
	}
}
func runTestsWithoutContainers(m *testing.M) {
	code := m.Run()
	os.Exit(code)
}

func runTestsWithContainers(m *testing.M, client *KibanaClient) {
	testContext, err := startKibana(GetEnvVarOrDefault("ELK_VERSION", defaultElkVersion), client)
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

	if client.Config.KibanaType == KibanaTypeVanilla {
		stopKibana(testContext)
	}

	os.Exit(code)
}

func defaultTestKibanaClient() *KibanaClient {
	kibanaClient := NewClient(NewDefaultConfig())
	kibanaClient.SetAuth(authForContainerVersion[kibanaClient.Config.KibanaVersion][kibanaClient.Config.KibanaType])
	return kibanaClient
}
