package containers

import (
	"gopkg.in/ory-am/dockertest.v3"
	"log"
	"os"
)

type TestContext struct {
	containers []container
	KibanaUri  string
}

func StartKibana() *TestContext {
	log.SetOutput(os.Stdout)

	var err error
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	elasticSearch := NewElasticSearchContainer(pool)
	kibana := NewKibanaContainer(pool, elasticSearch)
	return &TestContext{containers: []container{elasticSearch, kibana}, KibanaUri: kibana.Uri}
}

func StopKibana(testContext *TestContext) {

	for _, container := range testContext.containers {
		err := container.Stop()
		if err != nil {
			log.Printf("Could not stop container: %v \n", err)
		}
	}

}
