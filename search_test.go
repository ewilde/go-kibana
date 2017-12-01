package kibana

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_SearchCreate(t *testing.T) {
	client := NewClient(NewDefaultConfig())

	searchSource, err := NewSearchSourceBuilder().
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
		WithSearchSource(searchSource).
		Build()

	assert.Nil(t, err)

	result, err := NewClient(NewDefaultConfig()).Search().Create(request)

	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, request.Attributes.Title, result.Attributes.Title)
	assert.Equal(t, request.Attributes.Columns, result.Attributes.Columns)
	assert.Equal(t, request.Attributes.Sort, result.Attributes.Sort)
}
