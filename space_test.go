package kibana

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_SpaceGet(t *testing.T) {
	skipIfNotXpackSecurity(t)
	client := DefaultTestKibanaClient()
	spaceAPI := client.Space()

	space, err := spaceAPI.GetByID("default")

	assert.Nil(t, err)
	assert.NotNil(t, space, "Space retrieved from get by id was null.")
}

func Test_SpacePostEmpty(t *testing.T) {
	skipIfNotXpackSecurity(t)
	client := DefaultTestKibanaClient()
	spaceAPI := client.Space()
	basicSpace := &Space{
		Id:   "empty",
		Name: "empty",
	}
	err := spaceAPI.Create(basicSpace)

	assert.Nil(t, err, "Error creating space")
}

func Test_SpacePostBasic(t *testing.T) {
	skipIfNotXpackSecurity(t)
	client := DefaultTestKibanaClient()
	spaceAPI := client.Space()
	basicSpace := &Space{
		Id:               "basic",
		Name:             "Blue team",
		Description:      "Team managing water",
		Color:            "#0000ff",
		Initials:         "BL",
		DisabledFeatures: []string{"dev_tools"},
	}
	err := spaceAPI.Create(basicSpace)

	assert.Nil(t, err, "Error creating space")

	space, err := spaceAPI.GetByID("basic")

	assert.Nil(t, err)
	assert.NotNil(t, space, "Space retrieved from get by id was null.")
	assert.Equal(t, space.DisabledFeatures[0], "dev_tools")

}

func Test_SpaceDeleteBasic(t *testing.T) {
	skipIfNotXpackSecurity(t)
	client := DefaultTestKibanaClient()
	spaceAPI := client.Space()
	basicSpace := &Space{
		Id:   "blue",
		Name: "Blue team",
	}
	err := spaceAPI.Create(basicSpace)
	assert.Nil(t, err, "Error creating space")
	err = spaceAPI.Delete("blue")
	assert.Nil(t, err, "Error deleting space")
}

func Test_SpaceUpdateBasic(t *testing.T) {
	skipIfNotXpackSecurity(t)
	client := DefaultTestKibanaClient()
	spaceAPI := client.Space()
	basicSpace := &Space{
		Id:   "blue",
		Name: "Blue team",
	}
	err := spaceAPI.Create(basicSpace)
	assert.Nil(t, err, "Error creating space")
	basicSpace.Description = "The blue team is a nice team"

	err = spaceAPI.Update(basicSpace)
	assert.Nil(t, err, "Error updating space")
	err = spaceAPI.Delete("blue")
	assert.Nil(t, err, "Error deleting space")
}
