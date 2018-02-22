package kibana

import (
	"fmt"
	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_DashboardCreateFromSavedSearch(t *testing.T) {
	client := DefaultTestKibanaClient()

	searchRequest, _, err := createSearchRequest(client, t)
	assert.Nil(t, err)
	searchResponse, err := client.Search().Create(searchRequest)
	defer client.Search().Delete(searchResponse.Id)
	assert.Nil(t, err)

	visualizationApi := client.Visualization()

	visualizationRequest, err := NewVisualizationRequestBuilder().
		WithTitle("China errors").
		WithDescription("This visualization shows errors from china").
		WithVisualizationState("{\"title\":\"test kong vis\",\"type\":\"area\",\"params\":{\"grid\":{\"categoryLines\":false,\"style\":{\"color\":\"#eee\"}},\"categoryAxes\":[{\"id\":\"CategoryAxis-1\",\"type\":\"category\",\"position\":\"bottom\",\"show\":true,\"style\":{},\"scale\":{\"type\":\"linear\"},\"labels\":{\"show\":true,\"truncate\":100},\"title\":{\"text\":\"@timestamp date ranges\"}}],\"valueAxes\":[{\"id\":\"ValueAxis-1\",\"name\":\"LeftAxis-1\",\"type\":\"value\",\"position\":\"left\",\"show\":true,\"style\":{},\"scale\":{\"type\":\"linear\",\"mode\":\"normal\"},\"labels\":{\"show\":true,\"rotate\":0,\"filter\":false,\"truncate\":100},\"title\":{\"text\":\"Count\"}}],\"seriesParams\":[{\"show\":\"true\",\"type\":\"area\",\"mode\":\"stacked\",\"data\":{\"label\":\"Count\",\"id\":\"1\"},\"drawLinesBetweenPoints\":true,\"showCircles\":true,\"interpolate\":\"linear\",\"valueAxis\":\"ValueAxis-1\"}],\"addTooltip\":true,\"addLegend\":true,\"legendPosition\":\"right\",\"times\":[],\"addTimeMarker\":false},\"aggs\":[{\"id\":\"1\",\"enabled\":true,\"type\":\"count\",\"schema\":\"metric\",\"params\":{}},{\"id\":\"2\",\"enabled\":true,\"type\":\"date_range\",\"schema\":\"segment\",\"params\":{\"field\":\"@timestamp\",\"ranges\":[{\"from\":\"now-1h\",\"to\":\"now\"}]}}],\"listeners\":{}}").
		WithSavedSearchId(searchResponse.Id).
		Build()

	assert.Nil(t, err)

	visualizationResponse, err := visualizationApi.Create(visualizationRequest)
	defer visualizationApi.Delete(visualizationResponse.Id)

	dashboardApi := client.Dashboard()

	dashboardRequest, err := NewDashboardRequestBuilder().
		WithTitle("China errors").
		WithDescription("This dashboard shows errors from china").
		WithPanelsJson(fmt.Sprintf("[{\"size_x\":6,\"size_y\":3,\"panelIndex\":1,\"type\":\"visualization\",\"id\":\"%s\",\"col\":1,\"row\":1},{\"size_x\":6,\"size_y\":3,\"panelIndex\":2,\"type\":\"search\",\"id\":\"%s\",\"col\":7,\"row\":1,\"columns\":[\"_source\"],\"sort\":[\"@timestamp\",\"desc\"]}]", visualizationResponse.Id, searchResponse.Id)).
		WithOptionsJson("{\"darkTheme\":false}").
		WithUiStateJson("{\"P-1\":{\"vis\":{\"defaultColors\":{\"0 - 50\":\"rgb(0,104,55)\",\"50 - 75\":\"rgb(255,255,190)\",\"75 - 100\":\"rgb(165,0,38)\"}}}}").
		Build()

	assert.Nil(t, err)

	response, err := dashboardApi.Create(dashboardRequest)
	defer dashboardApi.Delete(response.Id)

	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, dashboardRequest.Attributes.Title, response.Attributes.Title)
	assert.Equal(t, dashboardRequest.Attributes.PanelsJson, response.Attributes.PanelsJson)
	assert.Equal(t, dashboardRequest.Attributes.OptionsJson, response.Attributes.OptionsJson)
	assert.Equal(t, dashboardRequest.Attributes.UiStateJSON, response.Attributes.UiStateJSON)
	assert.Equal(t, dashboardRequest.Attributes.Version, response.Attributes.Version)
}

func Test_DashboardRead(t *testing.T) {
	client := DefaultTestKibanaClient()
	dashboardApi := client.Dashboard()

	request, err := NewDashboardRequestBuilder().
		WithTitle("China errors").
		WithDescription("This dashboard shows errors from china").
		WithPanelsJson("[{\"size_x\":6,\"size_y\":3,\"panelIndex\":1,\"type\":\"visualization\",\"id\":\"bc8a1970-175b-11e8-accb-65182aaf9591\",\"col\":1,\"row\":1},{\"size_x\":6,\"size_y\":3,\"panelIndex\":2,\"type\":\"search\",\"id\":\"aca8b340-175b-11e8-accb-65182aaf9591\",\"col\":7,\"row\":1,\"columns\":[\"_source\"],\"sort\":[\"@timestamp\",\"desc\"]}]").
		WithOptionsJson("{\"darkTheme\":false}").
		WithUiStateJson("{\"P-1\":{\"vis\":{\"defaultColors\":{\"0 - 50\":\"rgb(0,104,55)\",\"50 - 75\":\"rgb(255,255,190)\",\"75 - 100\":\"rgb(165,0,38)\"}}}}").
		Build()

	assert.Nil(t, err)

	createdDashboard, err := dashboardApi.Create(request)
	defer dashboardApi.Delete(createdDashboard.Id)
	assert.Nil(t, err, "Error creating dashboard")

	readDashboard, err := dashboardApi.GetById(createdDashboard.Id)

	assert.Nil(t, err, "Error getting dashboard by id")
	assert.NotNil(t, readDashboard, "Dashboard retrieved from get by id was null.")

	assert.Equal(t, request.Attributes.Title, readDashboard.Attributes.Title)
	assert.Equal(t, request.Attributes.PanelsJson, readDashboard.Attributes.PanelsJson)
	assert.Equal(t, request.Attributes.OptionsJson, readDashboard.Attributes.OptionsJson)
	assert.Equal(t, request.Attributes.UiStateJSON, readDashboard.Attributes.UiStateJSON)
	assert.Equal(t, request.Attributes.Version, readDashboard.Attributes.Version)
}

func Test_DashboardRead_Unknown_Search_Returns_404(t *testing.T) {
	client := DefaultTestKibanaClient()

	dashboardClient := client.Dashboard()
	_, err := dashboardClient.GetById(uuid.NewV4().String())

	assert.NotNil(t, err, "Expected to get a 404 error")
	httpErr, ok := err.(*HttpError)
	if !ok {
		t.Error("Expected http error")
	}

	assert.Equal(t, 404, httpErr.Code)
}

func Test_DashboardUpdate(t *testing.T) {
	client := DefaultTestKibanaClient()
	dashboardApi := client.Dashboard()

	request, err := NewDashboardRequestBuilder().
		WithTitle("China errors").
		WithDescription("This dashboard shows errors from china").
		WithPanelsJson("[{\"size_x\":6,\"size_y\":3,\"panelIndex\":1,\"type\":\"visualization\",\"id\":\"bc8a1970-175b-11e8-accb-65182aaf9591\",\"col\":1,\"row\":1},{\"size_x\":6,\"size_y\":3,\"panelIndex\":2,\"type\":\"search\",\"id\":\"aca8b340-175b-11e8-accb-65182aaf9591\",\"col\":7,\"row\":1,\"columns\":[\"_source\"],\"sort\":[\"@timestamp\",\"desc\"]}]").
		WithOptionsJson("{\"darkTheme\":false}").
		WithUiStateJson("{\"P-1\":{\"vis\":{\"defaultColors\":{\"0 - 50\":\"rgb(0,104,55)\",\"50 - 75\":\"rgb(255,255,190)\",\"75 - 100\":\"rgb(165,0,38)\"}}}}").
		Build()

	assert.Nil(t, err)

	createdDashboard, err := dashboardApi.Create(request)
	defer dashboardApi.Delete(createdDashboard.Id)
	assert.Nil(t, err, "Error creating dashboard")

	createdDashboard.Attributes.Title = "China errors updated"
	updatedDashboard, err := dashboardApi.Update(createdDashboard.Id, &UpdateDashboardRequest{Attributes: createdDashboard.Attributes})
	assert.Nil(t, err)
	assert.Equal(t, "China errors updated", updatedDashboard.Attributes.Title)
}
