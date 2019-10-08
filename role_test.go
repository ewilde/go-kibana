package kibana

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func skipIfNotXpackSecurity(t *testing.T) {
	_, useXpackSecurity := os.LookupEnv("USE_XPACK_SECURITY")
	if !useXpackSecurity {
		t.Skip("Skipping testing as we don't have xpack security")
	}
}
func Test_RoleGet(t *testing.T) {
	skipIfNotXpackSecurity(t)
	client := DefaultTestKibanaClient()
	roleAPI := client.Role()

	role, err := roleAPI.GetByID("kibana_user")

	assert.Nil(t, err)
	assert.NotNil(t, role, "Role retrieved from get by id was null.")
}

func Test_RolePutEmpty(t *testing.T) {
	skipIfNotXpackSecurity(t)
	client := DefaultTestKibanaClient()
	roleAPI := client.Role()
	basicRole := &Role{
		Name:     "empty",
		Metadata: make(map[string]interface{}),
		ElasticSearch: &RoleElasticSearch{
			Cluster: []string{},
			Indices: make([]interface{}, 0),
			RunAs:   []string{},
		},
		Kibana: []*RoleKibana{},
	}
	err := roleAPI.CreateOrUpdate(basicRole)

	assert.Nil(t, err, "Error creating role")
}

func Test_RolePutBasic(t *testing.T) {
	skipIfNotXpackSecurity(t)
	client := DefaultTestKibanaClient()
	roleAPI := client.Role()
	basicRoleKibana := RoleKibana{
		Base:    []string{},
		Feature: make(map[string][]string),
		Spaces:  []string{"database"},
	}
	basicRoleKibana.Feature["dashboard"] = []string{"read"}
	roleMetadata := make(map[string]interface{})
	roleMetadata["version"] = 1
	basicRole := &Role{
		Name:     "basic",
		Metadata: roleMetadata,
		ElasticSearch: &RoleElasticSearch{
			Cluster: []string{},
			Indices: make([]interface{}, 0),
			RunAs:   []string{},
		},
		Kibana: []*RoleKibana{&basicRoleKibana},
	}
	err := roleAPI.CreateOrUpdate(basicRole)

	assert.Nil(t, err, "Error creating role")

	role, err := roleAPI.GetByID("basic")

	assert.Nil(t, err)
	assert.NotNil(t, role, "Role retrieved from get by id was null.")
	assert.Equal(t, role.Metadata["version"], 1.0)
	assert.Equal(t, role.Kibana[0].Spaces[0], "database")

}

func Test_RoleDeleteBasic(t *testing.T) {
	skipIfNotXpackSecurity(t)
	client := DefaultTestKibanaClient()
	roleAPI := client.Role()
	basicRoleKibana := RoleKibana{
		Base:    []string{},
		Feature: make(map[string][]string),
		Spaces:  []string{"database"},
	}
	basicRoleKibana.Feature["dashboard"] = []string{"read"}
	roleMetadata := make(map[string]interface{})
	roleMetadata["version"] = 1
	basicRole := &Role{
		Name:     "basic",
		Metadata: roleMetadata,
		ElasticSearch: &RoleElasticSearch{
			Cluster: []string{},
			Indices: make([]interface{}, 0),
			RunAs:   []string{},
		},
		Kibana: []*RoleKibana{&basicRoleKibana},
	}
	err := roleAPI.CreateOrUpdate(basicRole)
	assert.Nil(t, err, "Error creating role")
	err = roleAPI.Delete("basic")
	assert.Nil(t, err, "Error deleting role")
}
