package kibana

import (
	"github.com/google/go-querystring/query"
	"github.com/hashicorp/go-hclog"
	"github.com/parnurzeal/gorequest"
	"net/url"
	"os"
	"reflect"
	"strings"
)

const EnvKibanaUri = "KIBANA_URI"
const EnvKibanaIndexId = "KIBANA_INDEX_ID"

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

func (kibanaClient *KibanaClient) SavedObjects() *SavedObjectsClient {
	return &SavedObjectsClient{
		config: kibanaClient.config,
		client: kibanaClient.client,
	}
}

func addQueryString(currentUrl string, filter interface{}) (string, error) {
	v := reflect.ValueOf(filter)
	if v.Kind() == reflect.Ptr && v.IsNil() {
		return currentUrl, nil
	}

	url, err := url.Parse(currentUrl)
	if err != nil {
		return currentUrl, err
	}

	queryStringValues, err := query.Values(filter)
	if err != nil {
		return currentUrl, err
	}

	url.RawQuery = queryStringValues.Encode()
	return url.String(), nil
}
