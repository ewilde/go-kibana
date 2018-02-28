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

	request, _, err := createSearchRequest(client, t)
	assert.Nil(t, err)

	searchClient := client.Search()
	searchFromAccount1, err := searchClient.Create(request)
	assert.Nil(t, err)

	// 2. swap over to the second account (this is the action of our test)
	err = client.ChangeAccount(os.Getenv("LOGZ_IO_ACCOUNT_ID_2"))
	if !assert.Nil(t, err) {
		t.Fatal()
	}

	// 3. Now assert this search does not exist when logged into this account
	_, err = searchClient.GetById(searchFromAccount1.Id)
	if !assert.NotNil(t, err) {
		t.Fatal()
	}

	httpErr, ok := err.(*HttpError)
	if !ok {
		t.Error("Expected http error")
	}

	if !assert.Equal(t, 404, httpErr.Code) {
		t.Fatal()
	}

	// 4. Swap back to the first account and clean up
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
