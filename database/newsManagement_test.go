package database

import (
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"testing"
)

func TestNewsManagement(t *testing.T) {
	TestSetup()
	t.Run("testPublicationCycle", testPublicationCycle)
	t.Run("testSetupPublicationList", testSetupPublicationList)
	t.Run("testListsPublication", testListsPublication)
	t.Run("testArticleCycle", testArticleCycle)
	t.Run("testSetupArticleList", testSetupArticleList)
	t.Run("testArticleList", testArticleList)
}

func testPublicationCycle(t *testing.T) {
	pub := Publication{
		UUID:         "pub_test",
		Publicated:   true,
		PublishTime:  testGetTime(t, "1900-01-01"),
		BreakingNews: false,
	}
	err := pub.CreateMe()
	assert.Nil(t, err)
	err = pub.GetByID("pub_test")
	assert.Nil(t, err)
	assert.Equal(t, "pub_test", pub.UUID)
	assert.Equal(t, false, pub.BreakingNews)
	pub.BreakingNews = true
	err = pub.SaveChanges()
	assert.Nil(t, err)
	err = pub.GetByID("pub_test")
	assert.Nil(t, err)
	assert.Equal(t, true, pub.BreakingNews)
	err = pub.DeleteMe()
	assert.Nil(t, err)
	err = pub.GetByID("pub_test")
	assert.Equal(t, gorm.ErrRecordNotFound, err)
}

func testSetupPublicationList(t *testing.T) {
	pub := Publication{
		UUID:         "pub_test",
		Publicated:   false,
		BreakingNews: false,
	}
	err := pub.CreateMe()
	assert.Nil(t, err)
	pub.Publicated = true
	pub.UUID, pub.PublishTime = "pub_public_test1", testGetTime(t, "1900-01-20")
	err = pub.CreateMe()
	assert.Nil(t, err)
	pub.UUID, pub.PublishTime = "pub_public_test2", testGetTime(t, "1900-01-19")
	err = pub.CreateMe()
	assert.Nil(t, err)
	pub.UUID, pub.PublishTime = "pub_public_test3", testGetTime(t, "1900-01-18")
	err = pub.CreateMe()
	assert.Nil(t, err)
	pub.UUID, pub.PublishTime = "pub_public_test4", testGetTime(t, "1900-01-17")
	err = pub.CreateMe()
	assert.Nil(t, err)
	pub.UUID, pub.PublishTime = "pub_public_test5", testGetTime(t, "1900-01-16")
	err = pub.CreateMe()
	assert.Nil(t, err)
}

func testListsPublication(t *testing.T) {
	list := PublicationList{}
	err := list.GetOnlyUnpublicated()
	assert.Nil(t, err)
	exists := false
	for _, p := range list {
		if p.UUID == "pub_test" {
			exists = true
		}
	}
	assert.True(t, exists)
	err = list.GetPublicationBeforeDate(testGetTime(t, "1900-01-18"), 12)
	assert.Nil(t, err)
	assert.Equal(t, 2, len(list))
	assert.Equal(t, "pub_public_test4", list[0].UUID)
	assert.Equal(t, "pub_public_test5", list[1].UUID)
	err, exists = list.GetPublicationAfter("pub_public_test3", 12)
	assert.Nil(t, err)
	assert.True(t, exists)
	assert.Equal(t, 2, len(list))
	assert.Equal(t, "pub_public_test4", list[0].UUID)
	assert.Equal(t, "pub_public_test5", list[1].UUID)
	err, exists = list.GetPublicationBefore("pub_public_test3", 2)
	assert.Nil(t, err)
	assert.True(t, exists)
	assert.Equal(t, 2, len(list))
	assert.Equal(t, "pub_public_test1", list[0].UUID)
	assert.Equal(t, "pub_public_test2", list[1].UUID)
}

func testArticleCycle(t *testing.T) {
	pub := Article{
		UUID:        "pub_test",
		Publication: "pub_public_test1",
		Headline:    "test",
	}
	err := pub.CreateMe()
	assert.Nil(t, err)
	err = pub.GetByID("pub_test")
	assert.Nil(t, err)
	assert.Equal(t, "pub_test", pub.UUID)
	assert.Equal(t, "test", pub.Headline)
	pub.Headline = "test_change"
	err = pub.SaveChanges()
	assert.Nil(t, err)
	err = pub.GetByID("pub_test")
	assert.Nil(t, err)
	assert.Equal(t, "test_change", pub.Headline)
	err = pub.DeleteMe()
	assert.Nil(t, err)
	err = pub.GetByID("pub_test")
	assert.Equal(t, gorm.ErrRecordNotFound, err)
}

func testSetupArticleList(t *testing.T) {
	pub := Article{
		UUID:        "pub_test",
		Publication: "pub_public_test1",
		Written:     testGetTime(t, "1900-01-01"),
	}
	err := pub.CreateMe()
	assert.Nil(t, err)
	pub = Article{
		UUID:        "pub_change_test1",
		Publication: "pub_public_test2",
		Written:     testGetTime(t, "1900-01-01"),
	}
	err = pub.CreateMe()
	assert.Nil(t, err)
	pub = Article{
		UUID:        "pub_change_test2",
		Publication: "pub_public_test2",
		Written:     testGetTime(t, "1900-01-02"),
	}
	err = pub.CreateMe()
	assert.Nil(t, err)
	pub = Article{
		UUID:        "pub_change_test3",
		Publication: "pub_public_test2",
		Written:     testGetTime(t, "1900-01-03"),
	}
	err = pub.CreateMe()
	assert.Nil(t, err)
	pub = Article{
		UUID:        "pub_change_test4",
		Publication: "pub_public_test2",
		Written:     testGetTime(t, "1900-01-04"),
	}
	err = pub.CreateMe()
	assert.Nil(t, err)
}

func testArticleList(t *testing.T) {
	list := ArticleList{}
	err := list.GetAllArticlesToPublication("pub_public_test1")
	assert.Nil(t, err)
	assert.Equal(t, 1, len(list))
	assert.Equal(t, "pub_test", list[0].UUID)
	err = list.GetAllArticlesToPublication("pub_public_test2")
	assert.Nil(t, err)
	assert.Equal(t, 4, len(list))
	assert.Equal(t, "pub_change_test1", list[0].UUID)
	assert.Equal(t, "pub_change_test2", list[1].UUID)
	assert.Equal(t, "pub_change_test3", list[2].UUID)
	assert.Equal(t, "pub_change_test4", list[3].UUID)

	pub := Publication{UUID: "pub_public_test2"}
	err = pub.UpdateAllArticles("pub_update_test")
	assert.Nil(t, err)
	err = list.GetAllArticlesToPublication("pub_update_test")
	assert.Nil(t, err)
	assert.Equal(t, 4, len(list))
	assert.Equal(t, "pub_change_test1", list[0].UUID)
	assert.Equal(t, "pub_change_test2", list[1].UUID)
	assert.Equal(t, "pub_change_test3", list[2].UUID)
	assert.Equal(t, "pub_change_test4", list[3].UUID)
}
