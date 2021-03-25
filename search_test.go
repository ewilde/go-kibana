package kibana

import (
	"encoding/json"
	"testing"

	goversion "github.com/mcuadros/go-version"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_SortUnmarshalJSON(t *testing.T) {
	s := &Sort{}
	err := json.Unmarshal([]byte(`["@timestamp", "desc"]`), s)
	assert.NoError(t, err)
	assert.Equal(t, &Sort{"@timestamp", "desc"}, s)

	s = &Sort{}
	err = json.Unmarshal([]byte(`[["@name", "desc"]]`), s)
	assert.NoError(t, err)
	assert.Equal(t, &Sort{"@name", "desc"}, s)
}

func Test_SearchCreate(t *testing.T) {
	client := DefaultTestKibanaClient()
	searchApi := client.Search()

	requestSearch, err := searchApi.NewSearchSource().
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
			Meta: &SearchFilterMetaData{
				Index:    client.Config.DefaultIndexId,
				Negate:   false,
				Disabled: false,
				Alias:    "China",
				Type:     "phrase",
				Key:      "geo.src",
				Value:    "CN",
				Params: &SearchFilterQueryAttributes{
					Query: "CN",
					Type:  "phrase",
				},
			},
		}).
		Build()

	assert.Nil(t, err)

	request, err := NewSearchRequestBuilder().
		WithTitle("Geography filter on china").
		WithDisplayColumns([]string{"_source"}).
		WithSortColumns([]string{"@timestamp"}, Descending).
		WithSearchSource(requestSearch).
		Build()

	require.NoError(t, err)

	response, err := searchApi.Create(request)
	require.NoError(t, err)
	defer searchApi.Delete(response.Id)

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

	assert.Equal(t, requestSearch.Filter[0].Meta.Type, responseSearch.Filter[0].Meta.Type)
	assert.Equal(t, requestSearch.Filter[0].Meta.Key, responseSearch.Filter[0].Meta.Key)
	assert.Equal(t, requestSearch.Filter[0].Meta.Alias, responseSearch.Filter[0].Meta.Alias)
	assert.Equal(t, requestSearch.Filter[0].Meta.Disabled, responseSearch.Filter[0].Meta.Disabled)
	assert.Equal(t, requestSearch.Filter[0].Meta.Negate, responseSearch.Filter[0].Meta.Negate)
	assert.Equal(t, requestSearch.Filter[0].Meta.Index, responseSearch.Filter[0].Meta.Index)
	assert.Equal(t, requestSearch.Filter[0].Meta.Params.Query, responseSearch.Filter[0].Meta.Params.Query)
	assert.Equal(t, requestSearch.Filter[0].Meta.Params.Type, responseSearch.Filter[0].Meta.Params.Type)
}

func Test_SearchCreateWithReferences(t *testing.T) {
	client := DefaultTestKibanaClient()
	if goversion.Compare(client.Config.KibanaVersion, "7.0.0", "<") {
		t.SkipNow()
	}

	searchApi := client.Search()

	requestSearch, err := searchApi.NewSearchSource().
		WithIndexRefName("kibanaSavedObjectMeta.searchSourceJSON.index").
		WithFilter(&SearchFilter{
			Query: &SearchFilterQuery{
				Match: map[string]*SearchFilterQueryAttributes{
					"geo.src": {
						Query: "CN",
						Type:  "phrase",
					},
				},
			},
			Meta: &SearchFilterMetaData{
				IndexRefName: "kibanaSavedObjectMeta.searchSourceJSON.filter[0].meta.index",
				Negate:       false,
				Disabled:     false,
				Alias:        "China",
				Type:         "phrase",
				Key:          "geo.src",
				Value:        "CN",
				Params: &SearchFilterQueryAttributes{
					Query: "CN",
					Type:  "phrase",
				},
			},
		}).
		Build()

	require.NoError(t, err)

	request, err := NewSearchRequestBuilder().
		WithTitle("Geography filter on china").
		WithDisplayColumns([]string{"_source"}).
		WithSortColumns([]string{"@timestamp"}, Descending).
		WithSearchSource(requestSearch).
		WithReferences([]*SearchReferences{
			{
				Id:   "logzioCustomerIndex*",
				Name: "kibanaSavedObjectMeta.searchSourceJSON.index",
				Type: SearchReferencesTypeIndexPattern,
			},
			{
				Id:   "logzioCustomerIndex*",
				Name: "kibanaSavedObjectMeta.searchSourceJSON.filter[0].meta.index",
				Type: SearchReferencesTypeIndexPattern,
			},
		}).
		Build()

	require.NoError(t, err)

	response, err := searchApi.Create(request)
	require.NoError(t, err)
	defer searchApi.Delete(response.Id)

	assert.NotNil(t, response)
	assert.Equal(t, request.Attributes.Title, response.Attributes.Title)
	assert.Equal(t, request.Attributes.Columns, response.Attributes.Columns)
	assert.Equal(t, request.Attributes.Sort, response.Attributes.Sort)
	assert.NotEmpty(t, request.Attributes.KibanaSavedObjectMeta.SearchSourceJSON)

	responseSearch := &SearchSource{}
	json.Unmarshal([]byte(response.Attributes.KibanaSavedObjectMeta.SearchSourceJSON), responseSearch)
	assert.Equal(t, "kibanaSavedObjectMeta.searchSourceJSON.index", responseSearch.IndexRefName)

	assert.Len(t, responseSearch.Filter, len(requestSearch.Filter))
	assert.Equal(t, requestSearch.Filter[0].Query.Match["geo.src"].Query, responseSearch.Filter[0].Query.Match["geo.src"].Query)
	assert.Equal(t, requestSearch.Filter[0].Query.Match["geo.src"].Type, responseSearch.Filter[0].Query.Match["geo.src"].Type)

	assert.Equal(t, requestSearch.Filter[0].Meta.Type, responseSearch.Filter[0].Meta.Type)
	assert.Equal(t, requestSearch.Filter[0].Meta.Key, responseSearch.Filter[0].Meta.Key)
	assert.Equal(t, requestSearch.Filter[0].Meta.Alias, responseSearch.Filter[0].Meta.Alias)
	assert.Equal(t, requestSearch.Filter[0].Meta.Disabled, responseSearch.Filter[0].Meta.Disabled)
	assert.Equal(t, requestSearch.Filter[0].Meta.Negate, responseSearch.Filter[0].Meta.Negate)
	assert.Equal(t, "", responseSearch.Filter[0].Meta.Index)
	assert.Equal(t, "kibanaSavedObjectMeta.searchSourceJSON.filter[0].meta.index", responseSearch.Filter[0].Meta.IndexRefName)
	assert.Equal(t, requestSearch.Filter[0].Meta.Params.Query, responseSearch.Filter[0].Meta.Params.Query)
	assert.Equal(t, requestSearch.Filter[0].Meta.Params.Type, responseSearch.Filter[0].Meta.Params.Type)
	assert.Equal(t, "logzioCustomerIndex*", response.References[0].Id)
	assert.Equal(t, "kibanaSavedObjectMeta.searchSourceJSON.index", response.References[0].Name)
	assert.Equal(t, SearchReferencesTypeIndexPattern, response.References[0].Type)
	assert.Equal(t, "logzioCustomerIndex*", response.References[1].Id)
	assert.Equal(t, "kibanaSavedObjectMeta.searchSourceJSON.filter[0].meta.index", response.References[1].Name)
	assert.Equal(t, SearchReferencesTypeIndexPattern, response.References[1].Type)
}

func Test_SearchCreate_with_exists_field(t *testing.T) {
	client := DefaultTestKibanaClient()
	searchApi := client.Search()

	requestSearch, err := searchApi.NewSearchSource().
		WithIndexId(client.Config.DefaultIndexId).
		WithFilter(&SearchFilter{
			Exists: &SearchFilterExists{
				Field: "geo.ip",
			},
			Meta: &SearchFilterMetaData{
				Index:    client.Config.DefaultIndexId,
				Negate:   false,
				Disabled: false,
				Type:     "exists",
				Key:      "geo.ip",
				Value:    "exists",
			},
		}).
		Build()

	require.NoError(t, err)

	request, err := NewSearchRequestBuilder().
		WithTitle("Geography filter on china").
		WithDisplayColumns([]string{"_source"}).
		WithSortColumns([]string{"@timestamp"}, Descending).
		WithSearchSource(requestSearch).
		Build()

	require.NoError(t, err)

	response, err := searchApi.Create(request)
	require.NoError(t, err)
	defer searchApi.Delete(response.Id)

	assert.NotNil(t, response)
	assert.Equal(t, request.Attributes.Title, response.Attributes.Title)
	assert.Equal(t, request.Attributes.Columns, response.Attributes.Columns)
	assert.Equal(t, request.Attributes.Sort, response.Attributes.Sort)
	assert.NotEmpty(t, request.Attributes.KibanaSavedObjectMeta.SearchSourceJSON)

	responseSearch := &SearchSource{}
	json.Unmarshal([]byte(response.Attributes.KibanaSavedObjectMeta.SearchSourceJSON), responseSearch)
	assert.Equal(t, requestSearch.IndexId, responseSearch.IndexId)

	assert.Len(t, responseSearch.Filter, len(requestSearch.Filter))
	assert.Equal(t, requestSearch.Filter[0].Exists.Field, responseSearch.Filter[0].Exists.Field)

	assert.Equal(t, requestSearch.Filter[0].Meta.Type, responseSearch.Filter[0].Meta.Type)
	assert.Equal(t, requestSearch.Filter[0].Meta.Key, responseSearch.Filter[0].Meta.Key)
	assert.Equal(t, requestSearch.Filter[0].Meta.Alias, responseSearch.Filter[0].Meta.Alias)
	assert.Equal(t, requestSearch.Filter[0].Meta.Disabled, responseSearch.Filter[0].Meta.Disabled)
	assert.Equal(t, requestSearch.Filter[0].Meta.Negate, responseSearch.Filter[0].Meta.Negate)
	assert.Equal(t, requestSearch.Filter[0].Meta.Index, responseSearch.Filter[0].Meta.Index)
}

func Test_SearchCreate_with_two_filters(t *testing.T) {
	client := DefaultTestKibanaClient()

	searchApi := client.Search()
	requestSearch, err := searchApi.NewSearchSource().
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

	require.NoError(t, err)

	request, err := NewSearchRequestBuilder().
		WithTitle("Geography filter on china with errors").
		WithDisplayColumns([]string{"_source"}).
		WithSortColumns([]string{"@timestamp"}, Descending).
		WithSearchSource(requestSearch).
		Build()

	require.NoError(t, err)

	response, err := searchApi.Create(request)
	require.NoError(t, err)
	defer searchApi.Delete(response.Id)

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

func Test_SearchCreate_with_query(t *testing.T) {
	client := DefaultTestKibanaClient()
	searchApi := client.Search()

	requestSearch, err := searchApi.NewSearchSource().
		WithIndexId(client.Config.DefaultIndexId).
		WithQuery("geo.src:china").
		Build()

	require.NoError(t, err)

	request, err := NewSearchRequestBuilder().
		WithTitle("Geography search on china with errors").
		WithDisplayColumns([]string{"_source"}).
		WithSortColumns([]string{"@timestamp"}, Descending).
		WithSearchSource(requestSearch).
		Build()

	require.NoError(t, err)

	response, err := searchApi.Create(request)
	require.NoError(t, err)
	defer searchApi.Delete(response.Id)

	assert.NotNil(t, response)
	assert.Equal(t, request.Attributes.Title, response.Attributes.Title)
	assert.Equal(t, request.Attributes.Columns, response.Attributes.Columns)
	assert.Equal(t, request.Attributes.Sort, response.Attributes.Sort)
	assert.NotEmpty(t, request.Attributes.KibanaSavedObjectMeta.SearchSourceJSON)

	responseSearch := &SearchSource{}
	json.Unmarshal([]byte(response.Attributes.KibanaSavedObjectMeta.SearchSourceJSON), responseSearch)
	assert.Equal(t, requestSearch.IndexId, responseSearch.IndexId)
	assert.NotNil(t, responseSearch.Query)
}

func Test_SearchRead(t *testing.T) {
	client := DefaultTestKibanaClient()
	searchClient := client.Search()

	request, requestSearch, err := createSearchRequest(searchClient, client.Config.DefaultIndexId, t)

	require.NoError(t, err)

	createdSearch, err := searchClient.Create(request)
	require.NoError(t, err, "Error creating search")
	defer searchClient.Delete(createdSearch.Id)

	readSearch, err := searchClient.GetById(createdSearch.Id)

	assert.NoError(t, err, "Error getting search by id")
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

	require.Error(t, err, "Expected to get a 404 error")
	httpErr, ok := err.(*HttpError)
	require.True(t, ok, "Expected http error")
	assert.Equal(t, 404, httpErr.Code)
}

func Test_SearchList(t *testing.T) {
	client := DefaultTestKibanaClient()
	if goversion.Compare(client.Config.KibanaVersion, "6.3.0", "<") {
		t.SkipNow()
	}

	searchClient := client.Search()

	request, _, err := createSearchRequest(searchClient, client.Config.DefaultIndexId, t)
	require.NoError(t, err)

	createdSearch, err := searchClient.Create(request)
	require.NoError(t, err, "Error creating search")
	defer searchClient.Delete(createdSearch.Id)

	listSearch, err := searchClient.List()
	require.NoError(t, err, "Error listing searches")
	assert.NotNil(t, listSearch, "Response from list search is null")
	assert.NotEmpty(t, listSearch, "Response from list search is empty")
}

func Test_SearchUpdate(t *testing.T) {
	client := DefaultTestKibanaClient()
	searchClient := client.Search()

	request, _, err := createSearchRequest(searchClient, client.Config.DefaultIndexId, t)
	require.NoError(t, err)
	search, err := searchClient.Create(request)
	require.NoError(t, err)
	defer func() {
		err = searchClient.Delete(search.Id)
		assert.Nil(t, err, "Delete returned error:%+v", err)
	}()

	search.Attributes.Title = "China updated"
	search, err = searchClient.Update(search.Id, &UpdateSearchRequest{Attributes: search.Attributes})
	assert.Nil(t, err)
	assert.Equal(t, "China updated", search.Attributes.Title)
}

func Test_SearchDelete(t *testing.T) {
	client := DefaultTestKibanaClient()
	searchClient := client.Search()

	request, _, err := createSearchRequest(searchClient, client.Config.DefaultIndexId, t)
	require.NoError(t, err)
	response, err := searchClient.Create(request)
	require.NoError(t, err)

	err = searchClient.Delete(response.Id)
	require.NoError(t, err, "Delete returned error:%+v", err)

	response, err = searchClient.GetById(response.Id)
	require.Nil(t, response, "Response should be nil after being deleted")
}

func createSearchRequest(factory SearchSourceBuilderFactory, indexId string, t *testing.T) (*CreateSearchRequest, *SearchSource, error) {
	requestSearch, err := factory.NewSearchSource().
		WithIndexId(indexId).
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

	request, err := NewSearchRequestBuilder().
		WithTitle("Geography filter on china").
		WithDisplayColumns([]string{"_source"}).
		WithSortColumns([]string{"@timestamp"}, Descending).
		WithSearchSource(requestSearch).
		Build()

	return request, requestSearch, err
}
