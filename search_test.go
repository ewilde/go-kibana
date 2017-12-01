package kibana

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_SearchCreate(t *testing.T) {
	client := NewClient(NewDefaultConfig())

	requestSearch, err := NewSearchSourceBuilder().
		WithIndexId(client.Config.DefaultIndexId).
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

	response, err := client.Search().Create(request)

	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, request.Attributes.Title, response.Attributes.Title)
	assert.Equal(t, request.Attributes.Columns, response.Attributes.Columns)
	assert.Equal(t, request.Attributes.Sort, response.Attributes.Sort)
	assert.NotEmpty(t, request.Attributes.KibanaSavedObjectMeta.SearchSourceJSON)

	responseSearch := &SearchSource{}
	json.Unmarshal([]byte(response.Attributes.KibanaSavedObjectMeta.SearchSourceJSON), responseSearch)
	assert.Equal(t, requestSearch.IndexId, responseSearch.IndexId)
	assert.Len(t, responseSearch.Filter, len(requestSearch.Filter))
	assert.Equal(t, requestSearch.Filter[0].Query.Match["geo.src"].Query, responseSearch.Filter[0].Query.Match["geo.src"].Query)
	assert.Equal(t, requestSearch.Filter[0].Query.Match["geo.src"].Type, responseSearch.Filter[0].Query.Match["geo.src"].Type)
}

func Test_SearchCreate_with_two_filters(t *testing.T) {
	client := NewClient(NewDefaultConfig())

	requestSearch, err := NewSearchSourceBuilder().
		WithIndexId(client.Config.DefaultIndexId).
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

	response, err := client.Search().Create(request)

	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, request.Attributes.Title, response.Attributes.Title)
	assert.Equal(t, request.Attributes.Columns, response.Attributes.Columns)
	assert.Equal(t, request.Attributes.Sort, response.Attributes.Sort)
	assert.NotEmpty(t, request.Attributes.KibanaSavedObjectMeta.SearchSourceJSON)

	responseSearch := &SearchSource{}
	json.Unmarshal([]byte(response.Attributes.KibanaSavedObjectMeta.SearchSourceJSON), responseSearch)
	assert.Equal(t, requestSearch.IndexId, responseSearch.IndexId)
	assert.Len(t, responseSearch.Filter, len(requestSearch.Filter))
	assert.Equal(t, requestSearch.Filter[0].Query.Match["geo.src"].Query, responseSearch.Filter[0].Query.Match["geo.src"].Query)
	assert.Equal(t, requestSearch.Filter[0].Query.Match["geo.src"].Type, responseSearch.Filter[0].Query.Match["geo.src"].Type)
	assert.Equal(t, requestSearch.Filter[1].Query.Match["@tags"].Query, responseSearch.Filter[1].Query.Match["@tags"].Query)
	assert.Equal(t, requestSearch.Filter[1].Query.Match["@tags"].Type, responseSearch.Filter[1].Query.Match["@tags"].Type)
}

func Test_SearchRead(t *testing.T) {
	client := NewClient(NewDefaultConfig())

	requestSearch, err := NewSearchSourceBuilder().
		WithIndexId(client.Config.DefaultIndexId).
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

	createdSearch, err := client.Search().Create(request)
	assert.Nil(t, err, "Error creating search")

	readSearch, err := client.Search().GetById(createdSearch.Id)
	assert.Nil(t, err, "Error getting search by id")
	assert.NotNil(t, readSearch, "Search retrieved from get by id was null.")

	assert.Equal(t, request.Attributes.Title, createdSearch.Attributes.Title)
	assert.Equal(t, request.Attributes.Columns, createdSearch.Attributes.Columns)
	assert.Equal(t, request.Attributes.Sort, createdSearch.Attributes.Sort)
	assert.NotEmpty(t, request.Attributes.KibanaSavedObjectMeta.SearchSourceJSON)

	responseSearch := &SearchSource{}
	json.Unmarshal([]byte(createdSearch.Attributes.KibanaSavedObjectMeta.SearchSourceJSON), responseSearch)
	assert.Equal(t, requestSearch.IndexId, responseSearch.IndexId)
	assert.Len(t, responseSearch.Filter, len(requestSearch.Filter))
	assert.Equal(t, requestSearch.Filter[0].Query.Match["geo.src"].Query, responseSearch.Filter[0].Query.Match["geo.src"].Query)
	assert.Equal(t, requestSearch.Filter[0].Query.Match["geo.src"].Type, responseSearch.Filter[0].Query.Match["geo.src"].Type)
}
