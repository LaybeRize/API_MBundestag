package database

import (
	"API_MBundestag/help"
	"database/sql"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var zwitsch Zwitscher
var hiddenList ZwitscherList
var normalList ZwitscherList

func TestZwitschers(t *testing.T) {
	helper.TestBlocker.Lock()
	defer helper.TestBlocker.Unlock()
	TestZwitscherDB()

	t.Run("testCreateZwitscher", testCreateZwitscher)
	t.Run("testGetByUUID", testGetByUUID)
	t.Run("testChangeToBlocked", testChangeToBlocked)
	t.Run("createAFewZwitscher", createAFewZwitscher)
	t.Run("testGetLatest", testGetLatest)
	t.Run("testGetCommentsFor", testGetCommentsFor)
	t.Run("testGetByUser", testGetByUser)
}

func testGetByUser(t *testing.T) {
	list := ZwitscherList{}
	err := list.GetByUser("test", 5, false)
	assert.Nil(t, err)
	assert.Equal(t, 2, len(list))
	assert.Equal(t, list[0].UUID, normalList[0].UUID)
	assert.Equal(t, list[1].UUID, normalList[1].UUID)
	err = list.GetByUser("test", 1, false)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(list))
	err = list.GetByUser("test", 5, true)
	assert.Nil(t, err)
	assert.Equal(t, 3, len(list))
	assert.Equal(t, list[0].UUID, hiddenList[0].UUID)
	assert.Equal(t, list[1].UUID, hiddenList[1].UUID)
	assert.Equal(t, list[2].UUID, hiddenList[2].UUID)
}

func testGetCommentsFor(t *testing.T) {
	list := ZwitscherList{}
	err := list.GetCommentsFor("lol", false)
	assert.Nil(t, err)
	assert.Equal(t, ZwitscherList{}, list)
	z := Zwitscher{
		UUID:        uuid.New().String(),
		Author:      "bruh",
		Flair:       "test",
		HTMLContent: "test",
		ConnectedTo: sql.NullString{
			String: "lol",
			Valid:  true,
		},
	}
	err = z.CreateMe()
	assert.Nil(t, err)
	err = list.GetCommentsFor("lol", false)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(list))
	list[0].Written = time.Time{}
	assert.Equal(t, ZwitscherList{z}, list)
	z.UUID = uuid.New().String()
	z.Blocked = true
	err = z.CreateMe()
	assert.Nil(t, err)
	err = z.SaveChanges()
	assert.Nil(t, err)
	err = list.GetCommentsFor("lol", false)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(list))
	assert.NotEqual(t, z.UUID, list[0].UUID)
	err = list.GetCommentsFor("lol", true)
	assert.Nil(t, err)
	assert.Equal(t, 2, len(list))
	assert.Equal(t, z.UUID, list[0].UUID)
}

func testGetLatest(t *testing.T) {
	list := ZwitscherList{}
	err := list.GetLatested(5, false)
	assert.Nil(t, err)
	assert.Equal(t, 2, len(list))
	assert.Equal(t, list[0].UUID, normalList[0].UUID)
	assert.Equal(t, list[1].UUID, normalList[1].UUID)
	err = list.GetLatested(1, false)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(list))
	err = list.GetLatested(5, true)
	assert.Nil(t, err)
	assert.Equal(t, 3, len(list))
	assert.Equal(t, list[0].UUID, hiddenList[0].UUID)
	assert.Equal(t, list[1].UUID, hiddenList[1].UUID)
	assert.Equal(t, list[2].UUID, hiddenList[2].UUID)
}

func createAFewZwitscher(t *testing.T) {
	hiddenList = ZwitscherList{zwitsch}
	time.Sleep(100)
	zwitsch.UUID = uuid.New().String()
	err := zwitsch.CreateMe()
	assert.Nil(t, err)
	hiddenList = append(ZwitscherList{zwitsch}, hiddenList...)
	normalList = ZwitscherList{zwitsch}
	time.Sleep(100)
	zwitsch.UUID = uuid.New().String()
	err = zwitsch.CreateMe()
	assert.Nil(t, err)
	hiddenList = append(ZwitscherList{zwitsch}, hiddenList...)
	normalList = append(ZwitscherList{zwitsch}, normalList...)
}

func testChangeToBlocked(t *testing.T) {
	zwitsch.Blocked = true
	err := zwitsch.SaveChanges()
	assert.Nil(t, err)
	res := Zwitscher{}
	err = res.GetByID(zwitsch.UUID)
	res.Written = time.Time{}
	assert.Nil(t, err)
	assert.Equal(t, zwitsch, res)
}

func testGetByUUID(t *testing.T) {
	res := Zwitscher{}
	err := res.GetByID(zwitsch.UUID)
	assert.Nil(t, err)
	res.Written = time.Time{}
	assert.Equal(t, zwitsch, res)
}

func testCreateZwitscher(t *testing.T) {
	zwitsch = Zwitscher{
		UUID:        uuid.New().String(),
		Author:      "test",
		Flair:       "test",
		HTMLContent: "test",
	}
	err := zwitsch.CreateMe()
	assert.Nil(t, err)
}
