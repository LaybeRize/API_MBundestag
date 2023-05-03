package dataLogic

import (
	"API_MBundestag/database"
	"API_MBundestag/help/generics"
	"database/sql"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"testing"
)

func TestTitleLogic(t *testing.T) {
	database.TestSetup()

	t.Run("testSetupAccountsForTitleLogic", testSetupAccountsForTitleLogic)
	t.Run("testCreateTitle", testCreateTitle)
	t.Run("testGetTitle", testGetTitle)
	t.Run("testChangeTitle", testChangeTitle)
	t.Run("testDeleteTitle", testDeleteTitle)
}

func testDeleteTitle(t *testing.T) {
	title := Title{}
	err := title.GetMe("test_titleLogic2")
	assert.Nil(t, err)

	var msg generics.Message
	var positve bool
	title.DeleteMe(&msg, &positve)
	assert.True(t, positve)
	assert.Equal(t, SuccessDeletedTitle+"\n", msg)

	acc := database.Account{}
	err = acc.GetByDisplayName("testTitleLogic")
	assert.Nil(t, err)
	assert.Equal(t, "", acc.Flair)
	err = acc.GetByDisplayName("testTitleLogic2")
	assert.Nil(t, err)
	assert.Equal(t, "", acc.Flair)
	err = acc.GetByDisplayName("testTitleLogic3")
	assert.Nil(t, err)
	assert.Equal(t, "", acc.Flair)
}

func testChangeTitle(t *testing.T) {
	title := Title{}
	err := title.GetMe("test_titleLogic")
	assert.Nil(t, err)

	title.Name = "test_titleLogic2"
	title.Flair = "test_tl2"
	title.SubGroup = "test_asjdlalsd"
	title.Holder = []string{"testTitleLogic2", "testTitleLogic3"}

	var msg generics.Message
	var positve bool
	title.ChangeMe(&msg, &positve)
	assert.True(t, positve)
	assert.Equal(t, SuccessChangedTitle+"\n", msg)

	titleDB := database.Title{}
	err = titleDB.GetByName("test_titleLogic")
	assert.Equal(t, gorm.ErrRecordNotFound, err)
	err = titleDB.GetByName("test_titleLogic2")
	assert.Nil(t, err)

	assert.Equal(t, title.Name, titleDB.Name)
	assert.Equal(t, title.Flair, titleDB.Flair.String)
	assert.Equal(t, title.MainGroup, titleDB.MainGroup)
	assert.Equal(t, title.SubGroup, titleDB.SubGroup)
	assert.Equal(t, len(title.Holder), len(titleDB.Holder))
	assert.Equal(t, "testTitleLogic2", titleDB.Holder[0].DisplayName)
	assert.Equal(t, "testTitleLogic3", titleDB.Holder[1].DisplayName)

	acc := database.Account{}
	err = acc.GetByDisplayName("testTitleLogic")
	assert.Nil(t, err)
	assert.Equal(t, "", acc.Flair)
	err = acc.GetByDisplayName("testTitleLogic2")
	assert.Nil(t, err)
	assert.Equal(t, "test_tl2", acc.Flair)
	err = acc.GetByDisplayName("testTitleLogic3")
	assert.Nil(t, err)
	assert.Equal(t, "test_tl2", acc.Flair)
}

func testGetTitle(t *testing.T) {
	expected := Title{
		OldName:   "test_titleLogic",
		Name:      "test_titleLogic",
		Flair:     "test_tl",
		MainGroup: "test_titleLogic",
		SubGroup:  "test_titleLogic",
		Holder:    []string{"testTitleLogic", "testTitleLogic3"},
	}

	title := Title{}
	err := title.GetMe("test_titleLogic")
	assert.Nil(t, err)
	assert.Equal(t, expected, title)
}

func testCreateTitle(t *testing.T) {
	title := Title{
		Name:      "test_titleLogic",
		MainGroup: "test_titleLogic",
		Flair:     "test_tl",
		SubGroup:  "test_titleLogic",
		Holder:    []string{"testTitleLogic", "testTitleLogic3"},
	}
	var msg generics.Message
	var positve bool
	title.CreateMe(&msg, &positve)
	assert.True(t, positve)
	assert.Equal(t, SuccessCreatedTitle+"\n", msg)

	titleDB := database.Title{}
	err := titleDB.GetByName("test_titleLogic")
	assert.Nil(t, err)
	assert.Equal(t, title.Name, titleDB.Name)
	assert.Equal(t, title.Flair, titleDB.Flair.String)
	assert.Equal(t, title.MainGroup, titleDB.MainGroup)
	assert.Equal(t, title.SubGroup, titleDB.SubGroup)
	assert.Equal(t, len(title.Holder), len(titleDB.Holder))
	assert.Equal(t, "testTitleLogic", titleDB.Holder[0].DisplayName)
	assert.Equal(t, "testTitleLogic3", titleDB.Holder[1].DisplayName)

	acc := database.Account{}
	err = acc.GetByDisplayName("testTitleLogic")
	assert.Nil(t, err)
	assert.Equal(t, "test_tl", acc.Flair)
	err = acc.GetByDisplayName("testTitleLogic2")
	assert.Nil(t, err)
	assert.Equal(t, "", acc.Flair)
	err = acc.GetByDisplayName("testTitleLogic3")
	assert.Nil(t, err)
	assert.Equal(t, "test_tl", acc.Flair)
}

func testSetupAccountsForTitleLogic(t *testing.T) {
	acc := database.Account{
		DisplayName: "testTitleLogic",
		Username:    "testTitleLogic",
		Password:    "test",
		Role:        database.User,
	}
	err := acc.CreateMe()
	assert.Nil(t, err)
	id := acc.ID
	acc = database.Account{
		DisplayName: "testTitleLogic2",
		Username:    "testTitleLogic2",
		Password:    "XXXX",
		Role:        database.User,
		Linked:      sql.NullInt64{Valid: true, Int64: id},
	}
	err = acc.CreateMe()
	assert.Nil(t, err)
	acc = database.Account{
		DisplayName: "testTitleLogic3",
		Username:    "testTitleLogic3",
		Password:    "XXXX",
		Role:        database.User,
	}
	err = acc.CreateMe()
	assert.Nil(t, err)
}
