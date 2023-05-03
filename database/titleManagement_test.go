package database

import (
	"database/sql"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"testing"
)

func TestTitleManagement(t *testing.T) {
	TestSetup()
	t.Run("testSetupAccountsTitle", testSetupAccountsTitle)
	t.Run("testTitleCycle", testTitleCycle)
	t.Run("testSetupTitles", testSetupTitles)
	t.Run("testMultipleTitles", testMultipleTitles)
	t.Run("testTitleMeFromDelete", testTitleMeFromDelete)
}

func testTitleMeFromDelete(t *testing.T) {
	acc := Account{}
	err := acc.GetByDisplayName("title_test1")
	assert.Nil(t, err)
	err = DeleteMeFromTitles(acc.ID)
	assert.Nil(t, err)
	title := Title{}
	err = title.GetByName("test_title1")
	assert.Nil(t, err)
	assert.Equal(t, 0, len(title.Holder))
	err = title.GetByName("test_title3")
	assert.Nil(t, err)
	assert.Equal(t, 1, len(title.Holder))
	assert.Equal(t, "title_test2", title.Holder[0].DisplayName)
}

func testMultipleTitles(t *testing.T) {
	list := TitleList{}
	err := list.GetAll()
	assert.Nil(t, err)
	counter := 0
	for _, title := range list {
		switch title.Name {
		case "test_title1", "test_title2", "test_title3":
			counter++
		}
	}
	assert.Equal(t, 3, counter)
	acc := Account{}
	err = acc.GetByDisplayName("title_test1")
	assert.Nil(t, err)
	err = list.GetAllForUserID(acc.ID)
	assert.Equal(t, "test_title1", list[0].Name)
	assert.Equal(t, "test_title3", list[1].Name)
	assert.Equal(t, 2, len(list))
}

func testSetupTitles(t *testing.T) {
	acc := Account{}
	err := acc.GetByDisplayName("title_test1")
	assert.Nil(t, err)
	title := Title{
		Name:      "test_title1",
		MainGroup: "test",
		SubGroup:  "test",
		Flair:     sql.NullString{Valid: true, String: "test"},
		Holder:    []Account{acc},
	}
	err = title.CreateMe()
	assert.Nil(t, err)
	title.Name = "test_title2"
	title.Holder = []Account{}
	err = title.CreateMe()
	assert.Nil(t, err)
	acc2 := Account{}
	err = acc2.GetByDisplayName("title_test2")
	assert.Nil(t, err)
	title.Name = "test_title3"
	title.Holder = []Account{acc, acc2}
	err = title.CreateMe()
	assert.Nil(t, err)
}

func testSetupAccountsTitle(t *testing.T) {
	forTestCreateAccount(t, "title_test1", Account{Role: User})
	forTestCreateAccount(t, "title_test2", Account{Role: Admin})
}

func testTitleCycle(t *testing.T) {
	acc := Account{}
	err := acc.GetByDisplayName("title_test1")
	assert.Nil(t, err)
	title := Title{
		Name:      "test_title",
		MainGroup: "test",
		SubGroup:  "test",
		Flair:     sql.NullString{Valid: true, String: "test"},
		Holder:    []Account{acc},
	}
	err = title.CreateMe()
	assert.Nil(t, err)
	second := Title{}
	err = second.GetByName("test_title")
	assert.Nil(t, err)
	assert.Equal(t, title, second)
	title.Holder = []Account{}
	title.MainGroup = "bazing"
	err = title.SaveChanges()
	assert.Nil(t, err)
	err = title.UpdateHolder()
	assert.Nil(t, err)
	err = second.GetByName("test_title")
	assert.Nil(t, err)
	assert.Equal(t, title, second)
	title.Name = "test_title_change"
	err = title.ChangeTitleName("test_title")
	assert.Nil(t, err)
	err = second.GetByName("test_title")
	assert.Equal(t, gorm.ErrRecordNotFound, err)
	err = second.GetByName("test_title_change")
	assert.Nil(t, err)
	assert.Equal(t, title, second)
	err = title.DeleteMe()
	assert.Nil(t, err)
	err = second.GetByName("test_title_change")
	assert.Equal(t, gorm.ErrRecordNotFound, err)
}
