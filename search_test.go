package kibana

import (
	"encoding/json"
	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_SearchCreate(t *testing.T) {
	client := DefaultTestKibanaClient()

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

	searchApi := client.Search()
	response, err := searchApi.Create(request)
	defer searchApi.Delete(response.Id)

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
	client := DefaultTestKibanaClient()

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

	searchApi := client.Search()
	response, err := searchApi.Create(request)
	defer searchApi.Delete(response.Id)

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
	client := DefaultTestKibanaClient()

	request, requestSearch, err := createSearchRequest(client, t)

	assert.Nil(t, err)

	searchClient := client.Search()
	createdSearch, err := searchClient.Create(request)
	defer searchClient.Delete(createdSearch.Id)
	assert.Nil(t, err, "Error creating search")

	readSearch, err := searchClient.GetById(createdSearch.Id)

	assert.Nil(t, err, "Error getting search by id")
	assert.NotNil(t, readSearch, "Search retrieved from get by id was null.")

	assert.Equal(t, request.Attributes.Title, readSearch.Attributes.Title)
	assert.Equal(t, request.Attributes.Columns, readSearch.Attributes.Columns)
	assert.Equal(t, request.Attributes.Sort, readSearch.Attributes.Sort)
	assert.NotEmpty(t, readSearch.Attributes.KibanaSavedObjectMeta.SearchSourceJSON)

	responseSearch := &SearchSource{}
	json.Unmarshal([]byte(readSearch.Attributes.KibanaSavedObjectMeta.SearchSourceJSON), responseSearch)
	assert.Equal(t, requestSearch.IndexId, responseSearch.IndexId)
	assert.Len(t, responseSearch.Filter, len(requestSearch.Filter))
	assert.Equal(t, requestSearch.Filter[0].Query.Match["geo.src"].Query, responseSearch.Filter[0].Query.Match["geo.src"].Query)
	assert.Equal(t, requestSearch.Filter[0].Query.Match["geo.src"].Type, responseSearch.Filter[0].Query.Match["geo.src"].Type)
}

func Test_SearchRead_Unknown_Search_Returns_404(t *testing.T) {
	client := DefaultTestKibanaClient()

	searchClient := client.Search()
	_, err := searchClient.GetById(uuid.NewV4().String())

	assert.NotNil(t, err, "Expected to get a 404 error")
	httpErr, ok := err.(*HttpError)
	if !ok {
		t.Error("Expected http error")
	}

	assert.Equal(t, 404, httpErr.Code)
}

func Test_Update(t *testing.T) {
	client := DefaultTestKibanaClient()

	request, _, err := createSearchRequest(client, t)
	assert.Nil(t, err)
	searchClient := client.Search()
	search, err := searchClient.Create(request)
	assert.Nil(t, err)
	defer func() {
		err = searchClient.Delete(search.Id)
		assert.Nil(t, err, "Delete returned error:%+v", err)
	}()

	search.Attributes.Title = "China updated"
	search, err = searchClient.Update(search.Id, &UpdateSearchRequest{Attributes: search.Attributes})
	assert.Nil(t, err)
	assert.Equal(t, "China updated", search.Attributes.Title)
}

func Test_Delete(t *testing.T) {
	client := DefaultTestKibanaClient()

	request, _, err := createSearchRequest(client, t)
	assert.Nil(t, err)
	response, err := client.Search().Create(request)
	assert.Nil(t, err)

	err = client.Search().Delete(response.Id)
	assert.Nil(t, err, "Delete returned error:%+v", err)

	response, err = client.Search().GetById(response.Id)
	assert.Nil(t, response, "Response should be nil after being deleted")
}

func createSearchRequest(client *KibanaClient, t *testing.T) (*CreateSearchRequest, *SearchSource, error) {
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

	return request, requestSearch, err
}
