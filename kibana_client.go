package kibana

import (
	"github.com/google/go-querystring/query"
	"net/url"
	"os"
	"reflect"
	"strings"
)

const EnvKibanaUri = "KIBANA_URI"
const EnvKibanaVersion = "ELK_VERSION"
const EnvKibanaIndexId = "KIBANA_INDEX_ID"
const DefaultKibanaUri = "http://localhost:5601"
const DefaultKibanaVersion = "6.0.0"

type Config struct {
	KibanaBaseUri  string
	DefaultIndexId string
	KibanaVersion  string
}

type KibanaClient struct {
	Config *Config
	client *HttpAgent
}

var indexClientFromVersion = map[string]func(kibanaClient *KibanaClient) IndexPatternClient{
	"6.0.0": func(kibanaClient *KibanaClient) IndexPatternClient {
		return &IndexPatternClient600{config: kibanaClient.Config, client: kibanaClient.client}
	},
	"5.5.3": func(kibanaClient *KibanaClient) IndexPatternClient {
		return &IndexPatternClient553{config: kibanaClient.Config, client: kibanaClient.client}
	},
}

var seachClientFromVersion = map[string]func(kibanaClient *KibanaClient) SearchClient{
	"6.0.0": func(kibanaClient *KibanaClient) SearchClient {
		return &SearchClient600{config: kibanaClient.Config, client: kibanaClient.client}
	},
	"5.5.3": func(kibanaClient *KibanaClient) SearchClient {
		return &SearchClient553{config: kibanaClient.Config, client: kibanaClient.client}
	},
}

func NewDefaultConfig() *Config {
	config := &Config{
		KibanaBaseUri: DefaultKibanaUri,
		KibanaVersion: DefaultKibanaVersion,
	}

	if os.Getenv(EnvKibanaUri) != "" {
		config.KibanaBaseUri = strings.TrimRight(os.Getenv(EnvKibanaUri), "/")
	}

	if os.Getenv(EnvKibanaVersion) != "" {
		config.KibanaVersion = os.Getenv(EnvKibanaVersion)
	}

	if os.Getenv(EnvKibanaIndexId) != "" {
		config.DefaultIndexId = os.Getenv(EnvKibanaIndexId)
	}

	return config
}

func NewClient(config *Config) *KibanaClient {
	return &KibanaClient{
		Config: config,
		client: NewHttpAgent(config, &NoAuthenticationHandler{}),
	}
}

func (kibanaClient *KibanaClient) SetAuth(handler AuthenticationHandler) *KibanaClient {
	kibanaClient.client.authHandler = handler
	return kibanaClient
}

func (kibanaClient *KibanaClient) Search() SearchClient {
	return seachClientFromVersion[kibanaClient.Config.KibanaVersion](kibanaClient)
}

func (kibanaClient *KibanaClient) IndexPattern() IndexPatternClient {
	return indexClientFromVersion[kibanaClient.Config.KibanaVersion](kibanaClient)
}

func (kibanaClient *KibanaClient) SavedObjects() *SavedObjectsClient {
	return &SavedObjectsClient{
		config: kibanaClient.Config,
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
