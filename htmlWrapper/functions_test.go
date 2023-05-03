package htmlWrapper

import (
	"API_MBundestag/database"
	"github.com/stretchr/testify/assert"
	"html/template"
	"testing"
	"time"
)

func TestNotZero(t *testing.T) {
	assert.False(t, notZero(0))
	assert.True(t, notZero(1))
}

func TestLenStrNotZero(t *testing.T) {
	assert.False(t, lenStrNotZero([]string{}))
	assert.True(t, lenStrNotZero([]string{"test"}))
}

func TestUeq(t *testing.T) {
	assert.False(t, ueq("a", "a"))
	assert.True(t, ueq("a", "b"))
}

func TestShowPublishButton(t *testing.T) {
	assert.False(t, showPublishButton(true, database.ArticleList{database.Article{}}))
	assert.False(t, showPublishButton(true, database.ArticleList{}))
	assert.False(t, showPublishButton(false, database.ArticleList{}))
	assert.True(t, showPublishButton(false, database.ArticleList{database.Article{}}))
}

func TestEqRole(t *testing.T) {
	assert.False(t, eqRole(database.Admin, 12))
	assert.False(t, eqRole(database.Admin, "test"))
	assert.True(t, eqRole(database.Admin, string(database.Admin)))
	assert.False(t, eqRole(database.Admin, database.User))
	assert.True(t, eqRole(database.Admin, database.Admin))
}

func TestFulfillsClass(t *testing.T) {
	assert.False(t, fulfillsClass(5, database.NotLoggedIn))
	assert.True(t, fulfillsClass(6, database.NotLoggedIn))
	assert.True(t, fulfillsClass(5, database.User))
	assert.True(t, fulfillsClass(4, database.MediaAdmin))
	assert.True(t, fulfillsClass(3, database.Admin))
	assert.True(t, fulfillsClass(2, database.HeadAdmin))
	assert.False(t, fulfillsClass(1, database.HeadAdmin))
}

func TestOrgStatus(t *testing.T) {
	assert.False(t, orgStatus(string(database.Public), database.Secret))
	assert.True(t, orgStatus(string(database.Public), database.Public))
}

func TestOneOfValues(t *testing.T) {
	assert.False(t, oneOfValues())
	assert.False(t, oneOfValues("test"))
	assert.False(t, oneOfValues("test", "lol"))
	assert.True(t, oneOfValues("test", "lol", "test"))
}

func TestOneOfValuesInArray(t *testing.T) {
	assert.False(t, oneOfValuesInArray("test", []string{}))
	assert.False(t, oneOfValuesInArray("test", []string{"lol"}))
	assert.True(t, oneOfValuesInArray("test", []string{"test", "lol"}))
}

func TestYesNo(t *testing.T) {
	assert.Equal(t, "test", yesno("", "test", false))
	assert.Equal(t, "test", yesno("test", "", true))
}

func TestPlural(t *testing.T) {
	assert.Equal(t, "test", plural("", "test", 12))
	assert.Equal(t, "test", plural("test", "", 1))
}

func TestDateFormat(t *testing.T) {
	timeVal, err := time.Parse("2006-01-02T15:04:05", "2012-02-13T12:15:37")
	assert.Nil(t, err)
	assert.Equal(t, "Das ist der Tag 13.02.2012", dateFormat("Das ist der Tag 02.01.2006", timeVal))
}

func TestWithFlair(t *testing.T) {
	assert.Equal(t, "Test", withFlair("Test", ""))
	assert.Equal(t, "Test, Flair", withFlair("Test", "Flair"))
}

func TestJsonFunc(t *testing.T) {
	type TestStruct struct {
		Bazinga string
		Test    bool
		Val     int
	}
	val := TestStruct{
		Bazinga: "bruh",
		Test:    false,
		Val:     129,
	}
	assert.Equal(t, "{\"Bazinga\":\"bruh\",\"Test\":false,\"Val\":129}", jsonFunc(val))
}

func TestRoleTranslation(t *testing.T) {
	assert.Equal(t, database.RoleTranslation[database.User], roleTranslations(database.User))
}

func TestStatusTranslation(t *testing.T) {
	assert.Equal(t, database.StatusTranslation[database.Public], statusTranslations(database.Public))
}

func TestArrayOrEmpty(t *testing.T) {
	assert.Equal(t, "Ist leer", arrayOrEmpty("Ist leer", []string{}))
	assert.Equal(t, "a, b, c", arrayOrEmpty("Ist leer", []string{"a", "b", "c"}))
}

func TestTitle(t *testing.T) {
	assert.Equal(t, "Dies Ist Ein Titel", title("dies ist ein titel"))
}

func TestNoEscape(t *testing.T) {
	assert.Equal(t, template.HTML("testfalse<p>bazinga</p>"), noescape("test", false, "<p>bazinga</p>"))
}

func TestNoEscapeURL(t *testing.T) {
	assert.Equal(t, template.URL("absdr qlnwe2%23%$%"), noescapeurl("absdr qlnwe2%23%$%"))
}

func TestQueryEscape(t *testing.T) {
	assert.Equal(t, template.URL("absdr+qlnwe2%2523%25%24%25"), queryEscape("absdr qlnwe2%23%$%"))
}

func TestUserLoop(t *testing.T) {
	i := 0
	for val := range userLoop("test", []string{"a", "b"}) {
		switch i {
		case 0:
			assert.Equal(t, UserLoop{
				Div:    "style=\"display: none\" id=\"divClassestest\" hidden",
				Input:  "id=\"inputClassestest\"",
				Button: "id=\"buttonClassestest\"",
			}, val)
		case 1:
			assert.Equal(t, UserLoop{
				Div:    "",
				Input:  "type=\"text\" value=\"a\"",
				Button: "onclick=\"deleteDivFromSelf(this)\"",
			}, val)
		case 2:
			assert.Equal(t, UserLoop{
				Div:    "",
				Input:  "type=\"text\" value=\"b\"",
				Button: "onclick=\"deleteDivFromSelf(this)\"",
			}, val)
		case 3:
			assert.Fail(t, "Should not get to this")
		}

		i++
	}
}

func TestRoleLoop(t *testing.T) {
	i := 0
	for val := range roleLoop(database.Admin) {
		if database.Roles[i] == string(database.Admin) {
			assert.Equal(t, template.HTMLAttr("selected"), val.Attribute)
		} else {
			assert.Equal(t, template.HTMLAttr(""), val.Attribute)
		}
		assert.Equal(t, database.Roles[i], val.Loop)
		i++
	}
}

func TestStatusLoop(t *testing.T) {
	i := 0
	for val := range statusLoop(database.Public) {
		if database.Stati[i] == string(database.Public) {
			assert.Equal(t, template.HTMLAttr("selected"), val.Attribute)
		} else {
			assert.Equal(t, template.HTMLAttr(""), val.Attribute)
		}
		assert.Equal(t, database.Stati[i], val.Loop)
		i++
	}
}
