package kibana

import (
	"encoding/json"
	"fmt"
	"errors"
)

const savedObjectsPath = "/api/saved_objects/"

type SavedObjectRequest struct {
	Type    string   `json:"type" url:"type"`
	Fields  []string `json:"fields" url:"fields"`
	PerPage int      `json:"per_page" url:"per_page"`
}

type SavedObjectRequestBuilder struct {
	objectType string
	fields     []string
	perPage    int
}

type SavedObjectsClient struct {
	config *Config
	client *HttpAgent
}

type SavedObjectResponse struct {
	Page         int            `json:"page"`
	PerPage      int            `json:"per_page"`
	Total        int            `json:"total"`
	SavedObjects []*SavedObject `json:"saved_objects"`
}

type SavedObject struct {
	Id         string                 `json:"id"`
	Type       string                 `json:"type"`
	Version    int                    `json:"version"`
	Attributes map[string]interface{} `json:"attributes"`
}

func (api *SavedObjectsClient) GetByType(request *SavedObjectRequest) (*SavedObjectResponse, error) {
	address, err := addQueryString(api.config.KibanaBaseUri+savedObjectsPath, request)

	if err != nil {
		return nil, fmt.Errorf("could not build query string for get saved objects by type, error: %v", err)
	}

	apiResponse, body, errs := api.client.Get(address).End()
	if errs != nil {
		return nil, fmt.Errorf("could not get saved objects, error: %v", errs)
	}

	if apiResponse.StatusCode >= 300 {
		return nil, errors.New(fmt.Sprintf("Status: %d, %s", apiResponse.StatusCode, body))
	}


	response := &SavedObjectResponse{}
	err = json.Unmarshal([]byte(body), response)
	if err != nil {
		return nil, fmt.Errorf("could not parse saved objects response, error: %v, response body: %s", err, apiResponse)
	}

	return response, nil
}

func NewSavedObjectRequestBuilder() *SavedObjectRequestBuilder {
	return &SavedObjectRequestBuilder{perPage: 20}
}

func (builder *SavedObjectRequestBuilder) WithType(objectType string) *SavedObjectRequestBuilder {
	builder.objectType = objectType
	return builder
}

func (builder *SavedObjectRequestBuilder) WithFields(fields []string) *SavedObjectRequestBuilder {
	builder.fields = fields
	return builder
}

func (builder *SavedObjectRequestBuilder) WithPerPage(perPage int) *SavedObjectRequestBuilder {
	builder.perPage = perPage
	return builder
}

func (builder *SavedObjectRequestBuilder) Build() *SavedObjectRequest {
	return &SavedObjectRequest{
		Fields:  builder.fields,
		Type:    builder.objectType,
		PerPage: builder.perPage,
	}
}
