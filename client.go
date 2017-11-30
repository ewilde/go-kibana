package kibana

import (
	"github.com/hashicorp/go-hclog"
	"github.com/parnurzeal/gorequest"
	"os"
	"strings"
)

const EnvKibanaUri = "KIBANA_URI"

type Config struct {
	HostAddress string
	Logger      hclog.Logger
}

type KibanaClient struct {
	config *Config
	client *gorequest.SuperAgent
}

func NewDefaultConfig() *Config {
	config := &Config{
		HostAddress: "http://localhost:5601",
	}

	if os.Getenv(EnvKibanaUri) != "" {
		config.HostAddress = strings.TrimRight(os.Getenv(EnvKibanaUri), "/")
	}

	return config
}

func NewClient(config *Config) *KibanaClient {
	return &KibanaClient{
		config: config,
		client: gorequest.New(),
	}
}

func (kibanaClient *KibanaClient) Discover() *DiscoverClient {
	return &DiscoverClient{
		config: kibanaClient.config,
		client: kibanaClient.client,
	}
}
