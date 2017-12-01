package kibana

import (
	"encoding/json"
	"github.com/parnurzeal/gorequest"
	"github.com/pkg/errors"
)

type DiscoverClient struct {
	config *Config
	client *gorequest.SuperAgent
}

type DiscoverRequest struct {
	Attributes *DiscoverRequestAttributes `json:"attributes"`
}

type SortOrder int

const (
	Ascending SortOrder = iota
	Descending
)

type DiscoverRequestAttributes struct {
	Title                 string                         `json:"title"`
	Description           string                         `json:"description"`
	Hits                  int                            `json:"hits"`
	Columns               []string                       `json:"columns"`
	Sort                  []string                       `json:"sort"`
	Version               int                            `json:"version"`
	KibanaSavedObjectMeta *DiscoverKibanaSavedObjectMeta `json:"kibanaSavedObjectMeta"`
}

type DiscoverKibanaSavedObjectMeta struct {
	searchSourceJSON string `json:"searchSourceJSON"`
}

type DiscoverSearchSource struct {
	IndexId      string         `json:"index"`
	HighlightAll bool           `json:"highlightAll"`
	Version      bool           `json:"version"`
	Query        *DiscoverQuery `json:"query"`
	FilterJson   []string       `json:"filter"`
}

type DiscoverQuery struct {
	Query    string `json:"query"`
	Language string `json:"language"`
}

type DiscoverRequestBuilder struct {
	title          string
	description    string
	displayColumns []string
	sortColumns    []string
	searchSource   *DiscoverSearchSource
}

type DiscoverSearchSourceBuilder struct {
	indexId      string
	highlightAll bool
	query        *DiscoverQuery
	filterJson   []string
}

func NewDiscoverSearchSourceBuilder() *DiscoverSearchSourceBuilder {
	return &DiscoverSearchSourceBuilder{}
}

func (builder *DiscoverSearchSourceBuilder) WithIndexId(indexId string) *DiscoverSearchSourceBuilder {
	builder.indexId = indexId
	return builder
}

func (builder *DiscoverSearchSourceBuilder) WithQuery(query *DiscoverQuery) *DiscoverSearchSourceBuilder {
	builder.query = query
	return builder
}

func (builder *DiscoverSearchSourceBuilder) WithFilterJson(filter string) *DiscoverSearchSourceBuilder {
	builder.filterJson = append(builder.filterJson)
	return builder
}

func (builder *DiscoverSearchSourceBuilder) Build() (*DiscoverSearchSource, error) {
	if builder.indexId == "" {
		return nil, errors.New("Index id is required to create a discover search source")
	}

	return &DiscoverSearchSource{
		IndexId:      builder.indexId,
		HighlightAll: builder.highlightAll,
		Version:      true,
		Query:        builder.query,
		FilterJson:   builder.filterJson,
	}, nil
}

func NewDiscoverRequestBuilder() *DiscoverRequestBuilder {
	return &DiscoverRequestBuilder{}
}

func (builder *DiscoverRequestBuilder) WithTitle(title string) *DiscoverRequestBuilder {
	builder.title = title
	return builder
}

func (builder *DiscoverRequestBuilder) WithDescription(description string) *DiscoverRequestBuilder {
	builder.description = description
	return builder
}

func (builder *DiscoverRequestBuilder) WithDisplayColumns(columns []string) *DiscoverRequestBuilder {
	builder.displayColumns = columns
	return builder
}

func (builder *DiscoverRequestBuilder) WithSortColumns(columns []string, order SortOrder) *DiscoverRequestBuilder {
	var sortOrder = ""
	if order == Descending {
		sortOrder = "desc"
	} else {
		sortOrder = "asc"
	}

	builder.sortColumns = append(columns, sortOrder)
	return builder
}

func (builder *DiscoverRequestBuilder) WithSearchSource(searchSource *DiscoverSearchSource) *DiscoverRequestBuilder {
	builder.searchSource = searchSource
	return builder
}

func (builder *DiscoverRequestBuilder) Build() (*DiscoverRequest, error) {
	searchSourceBytes, err := json.Marshal(builder.searchSource)
	if err != nil {
		return nil, err
	}

	request := &DiscoverRequest{
		Attributes: &DiscoverRequestAttributes{
			Title:       builder.title,
			Description: builder.description,
			Hits:        0,
			Columns:     builder.displayColumns,
			Sort:        builder.sortColumns,
			Version:     1,
			KibanaSavedObjectMeta: &DiscoverKibanaSavedObjectMeta{
				searchSourceJSON: string(searchSourceBytes),
			},
		},
	}
	return request, nil
}
