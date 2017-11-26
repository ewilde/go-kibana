package containers

import (
	"errors"
	"fmt"
	"gopkg.in/ory-am/dockertest.v3"
	"log"
	"net/http"
	"github.com/parnurzeal/gorequest"
)

type elasticSearchContainer struct {
	Name        string
	pool        *dockertest.Pool
	resource    *dockertest.Resource
	HostAddress string
}

func NewElasticSearchContainer(pool *dockertest.Pool, version string) *elasticSearchContainer {

	envVars := []string{
		"discovery.type=single-node",
	}

	options := &dockertest.RunOptions{
		Repository: "docker.elastic.co/elasticsearch/elasticsearch",
		Tag:        version,
		Env:        envVars,
	}

	resource, err := pool.RunWithOptions(options)
	elasticSearchAddress := fmt.Sprintf("http://localhost:%v", resource.GetPort("9300/tcp"))

	if err := pool.Retry(func() error {
		client := gorequest.New()
		response, body, err := client.Get("http://localhost:9300").End()
		if err != nil {
			return err[0]
		}

		if response.StatusCode >= 300 {
			return errors.New(fmt.Sprintf("Status: %d, %s", response.StatusCode ,body))
		}

		return nil
	}); err != nil {
		gokong.AppLogger  .Fatalf("Could not connect to elastic search: %s", err)
	}

	if err != nil {
		log.Fatalf("Could not connect to elastic search: %s", err)
	}

	return &elasticSearchContainer{
		Name:        getContainerName(resource),
		pool:        pool,
		resource:    resource,
		HostAddress: kongAddress,
	}
}

func (kong *elasticSearchContainer) Stop() error {
	return kong.pool.Purge(kong.resource)
}
