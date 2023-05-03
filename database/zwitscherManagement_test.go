package database

import (
	"database/sql"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestZwitscherManagement(t *testing.T) {
	TestSetup()
	t.Run("testCreateZwitscher", testCreateZwitscher)
	t.Run("testChangeZwitscher", testChangeZwitscher)
	t.Run("testCreateRelatedTweets", testCreateRelatedTweets)
	t.Run("testZwitscherQueryResult", testZwitscherQueryResult)
	t.Run("testArrayQueryZwitscher", testArrayQueryZwitscher)
}

func testArrayQueryZwitscher(t *testing.T) {
	counter := 0
	list := ZwitscherList{}
	err := list.GetLatested(20, false)
	assert.Nil(t, err)
	for _, z := range list {
		switch z.UUID {
		case "zwitscher_test1", "zwitscher_test2", "zwitscher_test3", "zwitscher_test4", "zwitscher_test5":
			counter++
		}
	}
	assert.Equal(t, 2, counter)
}

func testZwitscherQueryResult(t *testing.T) {
	z := Zwitscher{}
	err := z.GetByUUID("zwitscher_test3")
	assert.Nil(t, err)
	assert.NotNil(t, z.Parent)
	assert.Equal(t, 2, len(z.Children))
	assert.Equal(t, "zwitscher_test2", z.Parent.UUID)
	assert.Equal(t, "zwitscher_test4", z.Children[0].UUID)
	assert.Equal(t, "zwitscher_test5", z.Children[1].UUID)
}

func testCreateRelatedTweets(t *testing.T) {
	z := Zwitscher{
		UUID:           "zwitscher_test2",
		Written:        testGetTime(t, "3000-01-09"),
		Blocked:        false,
		Author:         "tset",
		Flair:          "test",
		HTMLContent:    "bazinga",
		ConnectedTo:    sql.NullString{},
		AmountChildren: 12,
	}
	err := z.CreateMe()
	assert.Nil(t, err)
	z.UUID = "zwitscher_test3"
	z.ConnectedTo = sql.NullString{Valid: true, String: "zwitscher_test2"}
	err = z.CreateMe()
	assert.Nil(t, err)
	z.UUID = "zwitscher_test4"
	z.ConnectedTo = sql.NullString{Valid: true, String: "zwitscher_test3"}
	err = z.CreateMe()
	assert.Nil(t, err)
	z.UUID = "zwitscher_test5"
	z.ConnectedTo = sql.NullString{Valid: true, String: "zwitscher_test3"}
	z.Written = testGetTime(t, "3000-01-10")
	err = z.CreateMe()
	assert.Nil(t, err)
}

func testChangeZwitscher(t *testing.T) {
	z := Zwitscher{}
	err := z.GetByUUID("zwitscher_test1")
	assert.Nil(t, err)
	z.Flair = "vbasdasd"
	z.Author = "basdasd"
	err = z.SaveChanges()
	assert.Nil(t, err)
	second := Zwitscher{}
	err = second.GetByUUID("zwitscher_test1")
	assert.Nil(t, err)
	assert.Equal(t, z, second)
}

func testCreateZwitscher(t *testing.T) {
	z := Zwitscher{
		UUID:           "zwitscher_test1",
		Written:        testGetTime(t, "3000-01-10"),
		Blocked:        false,
		Author:         "tset",
		Flair:          "test",
		HTMLContent:    "bazinga",
		ConnectedTo:    sql.NullString{},
		AmountChildren: 12,
		Children:       []Zwitscher{},
	}
	err := z.CreateMe()
	assert.Nil(t, err)
	second := Zwitscher{}
	err = second.GetByUUID("zwitscher_test1")
	assert.Nil(t, err)
	second.Written = second.Written.UTC()
	assert.Equal(t, z, second)
}
