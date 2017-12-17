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
const EnvKibanaType = "KIBANA_TYPE"
const DefaultKibanaUri = "http://localhost:5601"
const DefaultKibanaVersion6 = "6.0.0"
const DefaultKibanaVersion553 = "5.5.3"
const DefaultKibanaVersion = DefaultKibanaVersion6
const DefaultKibanaIndexId = "logstash-*"
const DefaultKibanaIndexIdLogzio = "[logzioCustomerIndex]YYMMDD"

type KibanaType int

var kibanaTypeNames = map[string]KibanaType{
	KibanaTypeVanilla.String(): KibanaTypeVanilla,
	KibanaTypeLogzio.String():  KibanaTypeLogzio,
}

const (
	KibanaTypeUnknown KibanaType = iota
	KibanaTypeVanilla
	KibanaTypeLogzio
)

func parseKibanaType(value string) KibanaType {
	kibanaType, ok := kibanaTypeNames[value]

	if !ok {
		return KibanaTypeUnknown
	}

	return kibanaType
}

type Config struct {
	DefaultIndexId string
	KibanaBaseUri  string
	KibanaVersion  string
	KibanaType     KibanaType
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

var savedObjectsClientFromVersion = map[string]func(kibanaClient *KibanaClient) SavedObjectsClient{
	"6.0.0": func(kibanaClient *KibanaClient) SavedObjectsClient {
		return &savedObjectsClient600{config: kibanaClient.Config, client: kibanaClient.client}
	},
	"5.5.3": func(kibanaClient *KibanaClient) SavedObjectsClient {
		return &savedObjectsClient553{config: kibanaClient.Config, client: kibanaClient.client}
	},
}

func NewDefaultConfig() *Config {
	config := &Config{
		KibanaBaseUri: DefaultKibanaUri,
		KibanaVersion: DefaultKibanaVersion,
		KibanaType:    KibanaTypeVanilla,
	}

	if value := os.Getenv(EnvKibanaUri); value != "" {
		config.KibanaBaseUri = strings.TrimRight(value, "/")
	}

	if value := os.Getenv(EnvKibanaVersion); value != "" {
		config.KibanaVersion = value
	}

	if value := os.Getenv(EnvKibanaType); value != "" {
		config.KibanaType = parseKibanaType(value)
	}

	if value := os.Getenv(EnvKibanaIndexId); value != "" {
		config.DefaultIndexId = value
	} else {
		if config.KibanaType == KibanaTypeVanilla {
			config.DefaultIndexId = DefaultKibanaIndexId
		} else {
			config.DefaultIndexId = DefaultKibanaIndexIdLogzio
		}
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

func (kibanaClient *KibanaClient) SavedObjects() SavedObjectsClient {
	return savedObjectsClientFromVersion[kibanaClient.Config.KibanaVersion](kibanaClient)
}

func addQueryString(currentUrl string, filter interface{}) (string, error) {
	v := reflect.ValueOf(filter)
	if v.Kind() == reflect.Ptr && v.IsNil() {
		return currentUrl, nil
	}

	uri, err := url.Parse(currentUrl)
	if err != nil {
		return currentUrl, err
	}

	queryStringValues, err := query.Values(filter)
	if err != nil {
		return currentUrl, err
	}

	uri.RawQuery = queryStringValues.Encode()
	return uri.String(), nil
}
