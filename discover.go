package kibana

import "github.com/parnurzeal/gorequest"

type DiscoverClient struct {
	config *Config
	client *gorequest.SuperAgent
}

type DiscoverRequest struct {
	Attributes DiscoverRequestAttributes `json:"attributes"`
}

type DiscoverRequestAttributes struct {
	Title                 string                        `json:"title"`
	Description           string                        `json:"description"`
	Columns               []string                      `json:"columns"`
	Sort                  []string                      `json:"sort"`
	Version               int                           `json:"version"`
	KibanaSavedObjectMeta DiscoverKibanaSavedObjectMeta `json:"kibanaSavedObjectMeta"`
}

type DiscoverKibanaSavedObjectMeta struct {
	searchSourceJSON string `json:"searchSourceJSON"`
}
