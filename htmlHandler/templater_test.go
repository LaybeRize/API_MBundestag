package htmlHandler

import (
	"API_MBundestag/dataLogic"
	"API_MBundestag/database"
	wr "API_MBundestag/htmlWrapper"
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"html/template"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestErrorValidator(t *testing.T) {
	temp := ValidationErrors{Info: "test"}
	assert.Equal(t, "test", temp.Error())
}

func TestMiddleHardwareForTests(t *testing.T) {
	w, ctx := GetEmptyContext(t)
	MiddleHardwareForTests(ctx)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "/", w.Header().Get("Location"))
}

func TestTemplating(t *testing.T) {
	temp, err := template.New("layout").Funcs(wr.DefaultFunctions).Parse("{{.Title}} {{.Page.Info}}")
	assert.Nil(t, err)
	Template = &wr.Templates{
		Extension: "",
		Dir:       "",
		Templates: map[string]*template.Template{
			"start": temp,
		},
	}

	t.Run("correctIdentity", correctIdentity)
	t.Run("correctTemplating", correctTemplating)
}

func correctTemplating(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	type StartPageStruct struct {
		LoggedIn bool
		Account  database.Account
		Info     string
	}
	PageIdentityMap[Identity(StartPageStruct{})] = BasicStruct{
		Title:    "Startseite",
		Site:     "start",
		Template: "start",
	}
	MakeSite(&StartPageStruct{Info: "adjknlödf"}, ctx, &database.Account{})
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "Startseite adjknlödf", w.Body.String())
}

func correctIdentity(t *testing.T) {
	type StartPageStruct struct {
		LoggedIn bool
		Account  database.Account
		Info     string
	}
	assert.Equal(t, "htmlHandler.StartPageStruct", Identity(StartPageStruct{}))
}

func TestReflection(t *testing.T) {
	database.TestSetup()

	t.Run("createAccountsReflection", createAccountsReflection)
	t.Run("createOrganisationsReflection", createOrganisationsReflection)
	createTitlesReflection()
	t.Run("reflectMainAndSubGroupsForTitles", reflectMainAndSubGroupsForTitles)
	t.Run("reflectNamesForTitles", reflectNamesForTitles)
	t.Run("reflectMainAndSubGroupsForOrganisations", reflectMainAndSubGroupsForOrganisations)
	t.Run("reflectNamesForOrganisations", reflectNamesForOrganisations)
	t.Run("reflectOnlyNotSuspendedNames", reflectOnlyNotSuspendedNames)
	t.Run("reflectAllNamesAndDisplaynames", reflectAllNamesAndDisplaynames)
	t.Run("reflectAllOwnAccounts", reflectAllOwnAccounts)
}

func reflectAllOwnAccounts(t *testing.T) {
	type Reflection struct {
		Accounts database.AccountList
		Message  string
	}
	ref := Reflection{}
	acc := database.Account{}
	err := acc.GetByDisplayName("a")
	assert.Nil(t, err)
	FillOwnAccounts(&ref, &acc)
	assert.Equal(t, 2, len(ref.Accounts))
	assert.Equal(t, "a", ref.Accounts[0].DisplayName)
	assert.Equal(t, "b", ref.Accounts[1].DisplayName)
}

func reflectAllNamesAndDisplaynames(t *testing.T) {
	type Reflection struct {
		Names   database.NameList
		Message string
	}
	ref := Reflection{}
	FillUserAndDisplayNames(&ref)
	assert.Equal(t, database.NameList{{"a", "a"},
		{"b", "b"},
		{"c", "c"},
	}, ref.Names)
	assert.Equal(t, "", ref.Message)
}

func reflectOnlyNotSuspendedNames(t *testing.T) {
	type Reflection struct {
		Names   []string
		Message string
	}
	ref := Reflection{}
	FillAllNotSuspendedNames(&ref)
	assert.Equal(t, []string{"a", "b"}, ref.Names)
	assert.Equal(t, "", ref.Message)
}

func reflectNamesForOrganisations(t *testing.T) {
	type Reflection struct {
		OrgNames []string
		Message  string
	}
	ref := Reflection{}
	FillOrganisationNames(&ref)
	assert.Equal(t, []string{"a", "b", "c"}, ref.OrgNames)
	assert.Equal(t, "", ref.Message)
}

func reflectMainAndSubGroupsForOrganisations(t *testing.T) {
	type Reflection struct {
		ExistingMainGroup []string
		ExistingSubGroup  []string
		Message           string
	}
	ref := Reflection{}
	FillOrganisationGroups(&ref)
	assert.Equal(t, []string{"a", "b", "c"}, ref.ExistingMainGroup)
	assert.Equal(t, []string{"a", "b"}, ref.ExistingSubGroup)
	assert.Equal(t, "", ref.Message)
}

func reflectNamesForTitles(t *testing.T) {
	type Reflection struct {
		TitleNames []string
	}
	ref := Reflection{}
	FillTitleNames(&ref)
	assert.Equal(t, dataLogic.GetTitleNames(), ref.TitleNames)
}

func reflectMainAndSubGroupsForTitles(t *testing.T) {
	type Reflection struct {
		ExistingMainGroup []string
		ExistingSubGroup  []string
	}
	ref := Reflection{}
	FillTitleGroups(&ref)
	assert.Equal(t, dataLogic.GetMainGroupNames(), ref.ExistingMainGroup)
	assert.Equal(t, dataLogic.GetSubGroupNames(), ref.ExistingSubGroup)
}

func createTitlesReflection() {
	dataLogic.MainGroupNames = []string{"asdjhb", "bhdasjd", "ascs"}
	dataLogic.SubGroupNames = []string{"asdbj", "hgsad", "resdf", "dgfs", "jhsawdsc"}
	dataLogic.TitleNames = []string{"yvcbds", "zhgdfse", "tsasdcg", "asfv"}
}

func createOrganisationsReflection(t *testing.T) {
	org := database.Organisation{
		Name:      "a",
		MainGroup: "a",
		SubGroup:  "a",
		Flair:     sql.NullString{},
		Status:    database.Public,
	}
	err := org.CreateMe()
	assert.Nil(t, err)
	org.Name, org.MainGroup, org.SubGroup = "b", "b", "b"
	err = org.CreateMe()
	assert.Nil(t, err)
	org.Name, org.MainGroup, org.SubGroup = "c", "c", "b"
	err = org.CreateMe()
	assert.Nil(t, err)
}

func createAccountsReflection(t *testing.T) {
	acc := database.Account{
		DisplayName: "a",
		Flair:       "a",
		Username:    "a",
		Password:    "a",
		Role:        database.User,
		Linked:      sql.NullInt64{},
	}
	err := acc.CreateMe()
	assert.Nil(t, err)
	acc.Role = database.PressAccount
	acc.Linked = sql.NullInt64{Valid: true, Int64: 1}
	acc.DisplayName, acc.Username = "b", "b"
	err = acc.CreateMe()
	assert.Nil(t, err)
	acc.DisplayName, acc.Username = "c", "c"
	err = acc.CreateMe()
	assert.Nil(t, err)
	acc.Suspended = true
	err = acc.SaveChanges()
	assert.Nil(t, err)
}
