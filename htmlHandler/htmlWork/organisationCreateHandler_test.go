package htmlWork

import (
	"API_MBundestag/database"
	"testing"
)

func TestOrganisationCreateHandler(t *testing.T) {
	database.TestSetup()

	t.Run("testSetupOrganisationCreatePage", testSetupOrganisationCreatePage)
	t.Run("testGetOrganisationCreatePage", testGetOrganisationCreatePage)
	t.Run("testPostOrganisationCreatePage", testPostOrganisationCreatePage)
}

func testPostOrganisationCreatePage(t *testing.T) {

}

func testGetOrganisationCreatePage(t *testing.T) {

}

func testSetupOrganisationCreatePage(t *testing.T) {

}
