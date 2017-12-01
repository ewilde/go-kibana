package kibana

import (
	"github.com/ewilde/go-kibana/containers"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"testing"
)

func Test_Newclient(t *testing.T) {
	result := NewClient(NewDefaultConfig())

	assert.NotNil(t, result)
	assert.Equal(t, os.Getenv(EnvKibanaUri), result.Config.HostAddress)
}

func TestMain(m *testing.M) {

	testContext, err := containers.StartKibana()
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

	containers.StopKibana(testContext)

	os.Exit(code)

}
