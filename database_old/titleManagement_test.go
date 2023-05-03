package database

import (
	"API_MBundestag/help"
	"database/sql"
	"github.com/stretchr/testify/assert"
	"testing"
)

var title Title

func TestTitles(t *testing.T) {
	helper.TestBlocker.Lock()
	defer helper.TestBlocker.Unlock()
	TestTitlesDB()

	t.Run("testCreateTitle", testCreateTitle)
	t.Run("testGetTitleByName", testGetTitleByName)
	t.Run("testEditTitle", testEditTitle)
	t.Run("testEditTitle", testOverwriteTitle)
	t.Run("testDeleteTitle", testDeleteTitle)
	t.Run("testTitleList", testTitleList)
	t.Run("testTitleListForUser", testTitleListForUser)
}

func testTitleListForUser(t *testing.T) {
	list := TitleList{}
	err := list.GetAllForDisplayName("dvcydsa")
	assert.Nil(t, err)
	assert.Equal(t, 1, len(list))
	assert.Equal(t, "test", list[0].Name)

	err = list.GetAllForDisplayName("other")
	assert.Nil(t, err)
	assert.Equal(t, 1, len(list))
	assert.Equal(t, "test2", list[0].Name)
}

func testTitleList(t *testing.T) {
	err := title.CreateMe()
	assert.Nil(t, err)
	title.Name = "test2"
	title.Flair.Valid = false
	title.Info.Names = []string{"other"}
	err = title.CreateMe()
	assert.Nil(t, err)
	list := TitleList{}
	err = list.GetAll()
	assert.Nil(t, err)
	assert.Equal(t, 2, len(list))
	assert.Equal(t, "test", list[0].Name)
	assert.Equal(t, "test2", list[1].Name)
}

func testDeleteTitle(t *testing.T) {
	err := title.DeleteMe()
	assert.Nil(t, err)
	err = title.GetByName("test")
	assert.Equal(t, sql.ErrNoRows, err)
}

func testOverwriteTitle(t *testing.T) {
	title.Name = "lol"
	err := title.ChangeTitleName("test")
	assert.Nil(t, err)
	test := Title{}
	err = test.GetByName("test")
	assert.Equal(t, sql.ErrNoRows, err)
	err = test.GetByName("lol")
	assert.Equal(t, title, test)
	title.Name = "test"
	err = title.ChangeTitleName("lol")
	assert.Nil(t, err)
}

func testEditTitle(t *testing.T) {
	title.MainGroup = "vcxaesr"
	title.SubGroup = "vbcxsdf"
	title.Info.Names = []string{"dvcydsa"}
	err := title.SaveChanges()
	assert.Nil(t, err)
	res := Title{}
	err = res.GetByName("test")
	assert.Nil(t, err)
	assert.Equal(t, title, res)
}

func testGetTitleByName(t *testing.T) {
	res := Title{}
	err := res.GetByName("test")
	assert.Nil(t, err)
	assert.Equal(t, title, res)
}

func testCreateTitle(t *testing.T) {
	title = Title{
		Name:      "test",
		MainGroup: "sdasd",
		SubGroup:  "bcxdfs",
		Flair: sql.NullString{
			String: "asd",
			Valid:  true,
		},
		Info: TitleInfo{Names: []string{"lol, bazinga"}},
	}
	err := title.CreateMe()
	assert.Nil(t, err)
}
