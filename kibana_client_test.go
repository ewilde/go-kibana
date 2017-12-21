package kibana

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func Test_NewClient(t *testing.T) {
	kibanaClient := DefaultTestKibanaClient()

	assert.NotNil(t, kibanaClient)
	assert.Equal(t, os.Getenv(EnvKibanaUri), kibanaClient.Config.KibanaBaseUri)
}

func TestMain(m *testing.M) {
	client := DefaultTestKibanaClient()

	if client.Config.KibanaType == KibanaTypeVanilla {
		RunTestsWithContainers(m, client)
	} else {
		RunTestsWithoutContainers(m)
	}
}
