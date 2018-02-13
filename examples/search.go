package examples

import (
	"github.com/ewilde/go-kibana"
)

func createSearch() (*kibana.Search, error) {
	client := kibana.NewClient(kibana.NewDefaultConfig())
	client.Config.KibanaVersion = "6.0.0"

	requestSearch, _ := kibana.NewSearchSourceBuilder().
		WithIndexId(client.Config.DefaultIndexId).
		WithFilter(&kibana.SearchFilter{
			Query: &kibana.SearchFilterQuery{
				Match: map[string]*kibana.SearchFilterQueryAttributes{
					"geo.src": {
						Query: "CN",
						Type:  "phrase",
					},
				},
			},
		}).
		Build()

	request, _ := kibana.NewSearchRequestBuilder().
		WithTitle("Geography filter on china with errors").
		WithDisplayColumns([]string{"_source"}).
		WithSortColumns([]string{"@timestamp"}, kibana.Descending).
		WithSearchSource(requestSearch).
		Build()

	return client.Search().Create(request)
}
