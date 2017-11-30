package kibana

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_DiscoverCreate(t *testing.T) {

	apiRequest := &ApiRequest{
		Name:                   "test-" + uuid.NewV4().String(),
		Hosts:                  []string{"example.com"},
		Uris:                   []string{"/example"},
		Methods:                []string{"GET", "POST"},
		UpstreamUrl:            "http://localhost:4140/testservice",
		StripUri:               false,
		PreserveHost:           true,
		Retries:                3,
		UpstreamConnectTimeout: 1000,
		UpstreamSendTimeout:    2000,
		UpstreamReadTimeout:    3000,
		HttpsOnly:              true,
		HttpIfTerminated:       false,
	}

	result, err := NewClient(NewDefaultConfig()).Apis().Create(apiRequest)

	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, apiRequest.Name, result.Name)
	assert.Equal(t, apiRequest.Hosts, result.Hosts)
	assert.Equal(t, apiRequest.Uris, result.Uris)
	assert.Equal(t, apiRequest.Methods, result.Methods)
	assert.Equal(t, apiRequest.UpstreamUrl, result.UpstreamUrl)
	assert.Equal(t, apiRequest.StripUri, result.StripUri)
	assert.Equal(t, apiRequest.PreserveHost, result.PreserveHost)
	assert.Equal(t, apiRequest.UpstreamConnectTimeout, result.UpstreamConnectTimeout)
	assert.Equal(t, apiRequest.UpstreamSendTimeout, result.UpstreamSendTimeout)
	assert.Equal(t, apiRequest.UpstreamReadTimeout, result.UpstreamReadTimeout)
	assert.Equal(t, apiRequest.HttpsOnly, result.HttpsOnly)
	assert.Equal(t, apiRequest.HttpIfTerminated, result.HttpIfTerminated)
}
