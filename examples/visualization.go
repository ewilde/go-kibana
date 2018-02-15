package examples

import "github.com/ewilde/go-kibana"

func createVisualization(search *kibana.Search) (*kibana.Visualization, error) {
	client := kibana.NewClient(kibana.NewDefaultConfig())
	client.Config.KibanaVersion = kibana.DefaultKibanaVersion6

	request, _ := kibana.NewVisualizationRequestBuilder().
		WithTitle("Geography filter on china with errors").
		WithDescription("Gauge visualization based on a saved search").
		WithVisualizationState(`{"title":"Chinese search","type":"gauge","params":{"type":"gauge","addTooltip":true,"addLegend":true,"gauge":{"verticalSplit":false,"extendRange":true,"percentageMode":false,"gaugeType":"Arc","gaugeStyle":"Full","backStyle":"Full","orientation":"vertical","colorSchema":"Green to Red","gaugeColorMode":"Labels","colorsRange":[{"from":0,"to":50},{"from":50,"to":75},{"from":75,"to":100}],"invertColors":false,"labels":{"show":true,"color":"black"},"scale":{"show":true,"labels":false,"color":"#333"},"type":"meter","style":{"bgWidth":0.9,"width":0.9,"mask":false,"bgMask":false,"maskBars":50,"bgFill":"#eee","bgColor":false,"subText":"","fontSize":60,"labelColor":true}}},"aggs":[{"id":"1","enabled":true,"type":"count","schema":"metric","params":{}}]}`).
		WithSavedSearchId(search.Id).
		Build()

	return client.Visualization().Create(request)
}
