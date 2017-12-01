package kibana

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_SearchCreate(t *testing.T) {
	client := NewClient(NewDefaultConfig())

	requestSearch, err := NewSearchSourceBuilder().
		WithIndexId(client.config.DefaultIndexId).
		WithFilter(&SearchFilter{
			Query: &SearchFilterQuery{
				Match: map[string]*SearchFilterQueryAttributes{
					"geo.src": {
						Query: "CN",
						Type:  "phrase",
					},
				},
			},
		}).
		Build()

	assert.Nil(t, err)

	request, err := NewRequestBuilder().
		WithTitle("Geography filter on china").
		WithDisplayColumns([]string{"_source"}).
		WithSortColumns([]string{"@timestamp"}, Descending).
		WithSearchSource(requestSearch).
		Build()

	assert.Nil(t, err)

	response, err := NewClient(NewDefaultConfig()).Search().Create(request)

	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, request.Attributes.Title, response.Attributes.Title)
	assert.Equal(t, request.Attributes.Columns, response.Attributes.Columns)
	assert.Equal(t, request.Attributes.Sort, response.Attributes.Sort)
	assert.NotEmpty(t, request.Attributes.KibanaSavedObjectMeta.SearchSourceJSON)

	responseSearch := &SearchSource{}
	json.Unmarshal([]byte(request.Attributes.KibanaSavedObjectMeta.SearchSourceJSON), responseSearch)
	assert.Equal(t, requestSearch.IndexId, responseSearch.IndexId)
	assert.Len(t, responseSearch.Filter, len(requestSearch.Filter))
	assert.Equal(t, requestSearch.Filter[0].Query.Match["geo.src"].Query, responseSearch.Filter[0].Query.Match["geo.src"].Query)
	assert.Equal(t, requestSearch.Filter[0].Query.Match["geo.src"].Type, responseSearch.Filter[0].Query.Match["geo.src"].Type)
}

func Test_SearchCreate_with_two_filters(t *testing.T) {
	client := NewClient(NewDefaultConfig())

	requestSearch, err := NewSearchSourceBuilder().
		WithIndexId(client.config.DefaultIndexId).
		WithFilter(&SearchFilter{
			Query: &SearchFilterQuery{
				Match: map[string]*SearchFilterQueryAttributes{
					"geo.src": {
						Query: "CN",
						Type:  "phrase",
					},
				},
			},
		}).
		WithFilter(&SearchFilter{
			Query: &SearchFilterQuery{
				Match: map[string]*SearchFilterQueryAttributes{
					"@tags": {
						Query: "error",
						Type:  "phrase",
					},
				},
			},
		}).
		Build()

	assert.Nil(t, err)

	request, err := NewRequestBuilder().
		WithTitle("Geography filter on china with errors").
		WithDisplayColumns([]string{"_source"}).
		WithSortColumns([]string{"@timestamp"}, Descending).
		WithSearchSource(requestSearch).
		Build()

	assert.Nil(t, err)

	response, err := NewClient(NewDefaultConfig()).Search().Create(request)

	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, request.Attributes.Title, response.Attributes.Title)
	assert.Equal(t, request.Attributes.Columns, response.Attributes.Columns)
	assert.Equal(t, request.Attributes.Sort, response.Attributes.Sort)
	assert.NotEmpty(t, request.Attributes.KibanaSavedObjectMeta.SearchSourceJSON)

	responseSearch := &SearchSource{}
	json.Unmarshal([]byte(request.Attributes.KibanaSavedObjectMeta.SearchSourceJSON), responseSearch)
	assert.Equal(t, requestSearch.IndexId, responseSearch.IndexId)
	assert.Len(t, responseSearch.Filter, len(requestSearch.Filter))
	assert.Equal(t, requestSearch.Filter[0].Query.Match["geo.src"].Query, responseSearch.Filter[0].Query.Match["geo.src"].Query)
	assert.Equal(t, requestSearch.Filter[0].Query.Match["geo.src"].Type, responseSearch.Filter[0].Query.Match["geo.src"].Type)
}
