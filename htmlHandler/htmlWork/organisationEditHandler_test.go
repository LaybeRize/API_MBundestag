package htmlWork

import (
	"API_MBundestag/database"
	"testing"
)

func TestOrganisationEditHandler(t *testing.T) {
	database.TestSetup()

	t.Run("testSetupOrganisationEditHandler", testSetupOrganisationEditHandler)
	t.Run("testGetOrganisationEditHandler", testGetOrganisationEditHandler)
	t.Run("testPostSearchOrganisationEditHandler", testPostSearchOrganisationEditHandler)
	t.Run("testPostOrganisationEditHandler", testPostOrganisationEditHandler)
}

func testPostOrganisationEditHandler(t *testing.T) {

}

func testPostSearchOrganisationEditHandler(t *testing.T) {

}

func testGetOrganisationEditHandler(t *testing.T) {

}

func testSetupOrganisationEditHandler(t *testing.T) {

}
