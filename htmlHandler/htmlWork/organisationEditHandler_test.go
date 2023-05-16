package htmlWork

import (
	"API_MBundestag/database"
	"API_MBundestag/help/generics"
	hHa "API_MBundestag/htmlHandler"
	"API_MBundestag/htmlHandler/htmlBasics"
	wr "API_MBundestag/htmlWrapper"
	"database/sql"
	"github.com/stretchr/testify/assert"
	"html/template"
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
	acc := database.Account{}
	err := acc.GetByUserName("testEditHandlerError")
	assert.Nil(t, err)
	w, ctx := hHa.GetContextSetForUser(t, acc)
	PostEditOrganisationPage(ctx)
	assert.Equal(t, "TestError|"+generics.NotAuthorizedToView, w.Body.String())

}

func testGetOrganisationEditHandler(t *testing.T) {
	acc := database.Account{}
	err := acc.GetByUserName("testEditHandlerError")
	assert.Nil(t, err)
	w, ctx := hHa.GetContextSetForUser(t, acc)
	GetEditOrganisationPage(ctx)
	assert.Equal(t, "TestError|"+generics.NotAuthorizedToView, w.Body.String())

}

func testSetupOrganisationEditHandler(t *testing.T) {
	database.TestSetup()
	Setup()
	htmlBasics.Setup()

	org := database.Organisation{
		Name:      "testEditHandler",
		MainGroup: "testEditHandler",
		SubGroup:  "testEditHandler",
		Flair:     sql.NullString{Valid: true, String: "testEditHandler"},
		Status:    database.Public,
	}
	err := org.CreateMe()
	assert.Nil(t, err)
	acc := database.Account{
		DisplayName: "testEditHandlerError",
		Username:    "testEditHandlerError",
		Role:        database.MediaAdmin,
	}
	err = acc.CreateMe()
	assert.Nil(t, err)
	acc = database.Account{
		DisplayName: "testEditHandler",
		Username:    "testEditHandler",
		Role:        database.Admin,
	}
	err = acc.CreateMe()
	assert.Nil(t, err)

	temp, err := template.New("layout").Funcs(wr.DefaultFunctions).Parse("TestEditOrganisation|{{.Page.Organisation}}|{{.Page.Message}}|{{.Page.Positiv}}")
	assert.Nil(t, err)
	hHa.SetTemplate(t, temp, "editOrganisation")
}
