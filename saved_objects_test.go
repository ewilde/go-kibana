package kibana

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_SavedObjectsGetByType(t *testing.T) {
	client := DefaultTestKibanaClient()

	result, err := client.SavedObjects().GetByType(
		NewSavedObjectRequestBuilder().
			WithFields([]string{"title"}).
			WithType("index-pattern").
			WithPerPage(15).
			Build())

	require.NoError(t, err)

	require.NotNil(t, result)
	assert.Equal(t, 1, result.Page)
	assert.Equal(t, 15, result.PerPage)
	assert.Equal(t, 1, result.Total)

	assert.Len(t, result.SavedObjects, 1)
	assert.NotZero(t, len(result.SavedObjects[0].Id))
	assert.Equal(t, "index-pattern", result.SavedObjects[0].Type)
	assert.NotZero(t, result.SavedObjects[0].Version)
	assert.NotNil(t, result.SavedObjects[0].Attributes)
}

func Test_SavedObjectsGetByType_with_multiple_fields(t *testing.T) {
	client := DefaultTestKibanaClient()

	result, err := client.SavedObjects().GetByType(
		NewSavedObjectRequestBuilder().
			WithFields([]string{"title", "timeFieldName", "fields"}).
			WithType("index-pattern").
			WithPerPage(15).
			Build())

	require.NoError(t, err)

	require.NotNil(t, result)
	assert.Equal(t, 1, result.Page)
	assert.Equal(t, 15, result.PerPage)
	assert.Equal(t, 1, result.Total)

	assert.Len(t, result.SavedObjects, 1)
	assert.NotZero(t, result.SavedObjects[0].Id)
	assert.Equal(t, "index-pattern", result.SavedObjects[0].Type)
	assert.NotZero(t, result.SavedObjects[0].Version)
	assert.NotNil(t, result.SavedObjects[0].Attributes)
	assert.Equal(t, "@timestamp", result.SavedObjects[0].Attributes["timeFieldName"])
	assert.NotEmpty(t, result.SavedObjects[0].Attributes["fields"])
}
