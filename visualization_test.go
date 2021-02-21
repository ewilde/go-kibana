package kibana

import (
	"testing"

	goversion "github.com/mcuadros/go-version"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func newTestVisualizationRequestBuilder() *VisualizationRequestBuilder {
	builder := NewVisualizationRequestBuilder().
		WithTitle("China errors").
		WithDescription("This visualization shows errors from china").
		WithVisualizationState("{\"title\":\"test kong vis\",\"type\":\"area\",\"params\":{\"grid\":{\"categoryLines\":false,\"style\":{\"color\":\"#eee\"}},\"categoryAxes\":[{\"id\":\"CategoryAxis-1\",\"type\":\"category\",\"position\":\"bottom\",\"show\":true,\"style\":{},\"scale\":{\"type\":\"linear\"},\"labels\":{\"show\":true,\"truncate\":100},\"title\":{\"text\":\"@timestamp date ranges\"}}],\"valueAxes\":[{\"id\":\"ValueAxis-1\",\"name\":\"LeftAxis-1\",\"type\":\"value\",\"position\":\"left\",\"show\":true,\"style\":{},\"scale\":{\"type\":\"linear\",\"mode\":\"normal\"},\"labels\":{\"show\":true,\"rotate\":0,\"filter\":false,\"truncate\":100},\"title\":{\"text\":\"Count\"}}],\"seriesParams\":[{\"show\":\"true\",\"type\":\"area\",\"mode\":\"stacked\",\"data\":{\"label\":\"Count\",\"id\":\"1\"},\"drawLinesBetweenPoints\":true,\"showCircles\":true,\"interpolate\":\"linear\",\"valueAxis\":\"ValueAxis-1\"}],\"addTooltip\":true,\"addLegend\":true,\"legendPosition\":\"right\",\"times\":[],\"addTimeMarker\":false},\"aggs\":[{\"id\":\"1\",\"enabled\":true,\"type\":\"count\",\"schema\":\"metric\",\"params\":{}},{\"id\":\"2\",\"enabled\":true,\"type\":\"date_range\",\"schema\":\"segment\",\"params\":{\"field\":\"@timestamp\",\"ranges\":[{\"from\":\"now-1h\",\"to\":\"now\"}]}}],\"listeners\":{}}").
		WithSavedSearchId("123")

	return builder
}

func Test_VisualizationCreateFromSavedSearch(t *testing.T) {
	client := DefaultTestKibanaClient()

	visualizationApi := client.Visualization()

	request, err := newTestVisualizationRequestBuilder().
		Build(client.Config.KibanaVersion)
	assert.Nil(t, err)

	response, err := visualizationApi.Create(request)
	defer visualizationApi.Delete(response.Id)

	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, request.Attributes.Title, response.Attributes.Title)
	assert.Equal(t, request.Attributes.VisualizationState, response.Attributes.VisualizationState)
	assert.Equal(t, request.Attributes.Version, response.Attributes.Version)
	assert.Equal(t, request.Attributes.SavedSearchId, response.Attributes.SavedSearchId)
}

func Test_VisualizationCreateWithReferences(t *testing.T) {
	client := DefaultTestKibanaClient()
	if goversion.Compare(client.Config.KibanaVersion, "7.0.0", "<") {
		t.SkipNow()
	}

	visualizationApi := client.Visualization()

	builder := newTestVisualizationRequestBuilder()
	builder.WithReferences([]*VisualizationReferences{
		{
			Id:   "456",
			Name: "TestIndex",
			Type: VisualizationReferencesTypeIndexPattern,
		},
		{
			Id:   "789",
			Name: "TestSearch",
			Type: VisualizationReferencesTypeSearch,
		},
	})
	request, err := builder.
		Build(client.Config.KibanaVersion)
	assert.Nil(t, err)

	response, err := visualizationApi.Create(request)
	defer visualizationApi.Delete(response.Id)

	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, request.Attributes.Title, response.Attributes.Title)
	assert.Equal(t, request.Attributes.VisualizationState, response.Attributes.VisualizationState)
	assert.Equal(t, request.Attributes.Version, response.Attributes.Version)
	assert.Equal(t, "", response.Attributes.SavedSearchId)
	assert.NotEmpty(t, response.References)
	assert.Equal(t, "456", response.References[0].Id)
	assert.Equal(t, "TestIndex", response.References[0].Name)
	assert.Equal(t, VisualizationReferencesTypeIndexPattern, response.References[0].Type)
	assert.Equal(t, "789", response.References[1].Id)
	assert.Equal(t, "TestSearch", response.References[1].Name)
	assert.Equal(t, VisualizationReferencesTypeSearch, response.References[1].Type)
}

func Test_VisualizationRead(t *testing.T) {
	client := DefaultTestKibanaClient()
	visualizationApi := client.Visualization()

	request, err := newTestVisualizationRequestBuilder().
		Build(client.Config.KibanaVersion)
	assert.Nil(t, err)

	createdVisualization, err := visualizationApi.Create(request)
	defer visualizationApi.Delete(createdVisualization.Id)
	assert.Nil(t, err, "Error creating visualization")

	readVisualization, err := visualizationApi.GetById(createdVisualization.Id)

	assert.Nil(t, err, "Error getting visualization by id")
	assert.NotNil(t, readVisualization, "Visualization retrieved from get by id was null.")

	assert.Equal(t, request.Attributes.Title, readVisualization.Attributes.Title)
	assert.Equal(t, request.Attributes.Description, readVisualization.Attributes.Description)
	assert.Equal(t, request.Attributes.VisualizationState, readVisualization.Attributes.VisualizationState)
	assert.Equal(t, request.Attributes.SavedSearchId, readVisualization.Attributes.SavedSearchId)
}

func Test_VisualizationRead_Unknown_Visualization_Returns_404(t *testing.T) {
	client := DefaultTestKibanaClient()

	visualizationClient := client.Visualization()
	_, err := visualizationClient.GetById(uuid.NewV4().String())

	assert.NotNil(t, err, "Expected to get a 404 error")
	httpErr, ok := err.(*HttpError)
	if !ok {
		t.Error("Expected http error")
	}

	assert.Equal(t, 404, httpErr.Code)
}

func Test_VisualizationList(t *testing.T) {
	client := DefaultTestKibanaClient()
	if goversion.Compare(client.Config.KibanaVersion, "6.3.0", "<") {
		t.SkipNow()
	}
	visualizationApi := client.Visualization()

	request, err := newTestVisualizationRequestBuilder().
		Build(client.Config.KibanaVersion)
	assert.Nil(t, err)

	createdVisualization, err := visualizationApi.Create(request)
	defer visualizationApi.Delete(createdVisualization.Id)
	assert.Nil(t, err, "Error creating visualization")

	listVisualization, err := visualizationApi.List()
	assert.Nil(t, err, "Error listing visualizations")
	assert.NotNil(t, listVisualization, "Response from list visualization is null")
	assert.NotEmpty(t, listVisualization, "Response from list visualization is empty")
}

func Test_VisualizationUpdate(t *testing.T) {
	client := DefaultTestKibanaClient()
	visualizationApi := client.Visualization()

	request, err := newTestVisualizationRequestBuilder().
		Build(client.Config.KibanaVersion)
	assert.Nil(t, err)

	createdVisualization, err := visualizationApi.Create(request)
	defer visualizationApi.Delete(createdVisualization.Id)
	assert.Nil(t, err, "Error creating visualization")

	createdVisualization.Attributes.Title = "China errors updated"
	updatedVisualization, err := visualizationApi.Update(createdVisualization.Id, &UpdateVisualizationRequest{Attributes: createdVisualization.Attributes})
	assert.Nil(t, err)
	assert.Equal(t, "China errors updated", updatedVisualization.Attributes.Title)
}
