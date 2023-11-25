package examples

import (
	"github.com/irfannurhakim/go-kibana"
)

func createSearch() (*kibana.Search, error) {
	client := kibana.NewClient(kibana.NewDefaultConfig())
	client.Config.KibanaVersion = kibana.DefaultKibanaVersion6

	requestSearch, _ := client.Search().NewSearchSource().
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
