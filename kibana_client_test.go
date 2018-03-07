package kibana

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func Test_NewClient(t *testing.T) {
	kibanaClient := DefaultTestKibanaClient()

	assert.NotNil(t, kibanaClient)
	assert.Equal(t, os.Getenv(EnvKibanaUri), kibanaClient.Config.KibanaBaseUri)
}

func Test_Change_account(t *testing.T) {
	if skip := testPreCheckForLogz(t); skip {
		t.Skip()
		return
	}

	// 1. create a saved search using the default account
	client := DefaultTestKibanaClient()
	searchClient := client.Search()

	request, _, err := createSearchRequest(searchClient, client.Config.DefaultIndexId, t)
	assert.Nil(t, err)

	searchFromAccount1, err := searchClient.Create(request)
	assert.Nil(t, err)

	// 2. swap over to the second account (this is the action of our test)
	err = client.ChangeAccount(os.Getenv("LOGZ_IO_ACCOUNT_ID_2"))
	if !assert.Nil(t, err) {
		t.Fatal()
	}

	// 3. Swap back to the first account and clean up
	client.ChangeAccount(os.Getenv("LOGZ_IO_ACCOUNT_ID_1"))
	err = searchClient.Delete(searchFromAccount1.Id)
	assert.Nil(t, err)
}

func TestMain(m *testing.M) {
	client := DefaultTestKibanaClient()

	if client.Config.KibanaType == KibanaTypeVanilla {
		RunTestsWithContainers(m, client)
	} else {
		RunTestsWithoutContainers(m)
	}
}
