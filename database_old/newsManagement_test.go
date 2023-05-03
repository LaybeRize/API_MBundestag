package database

import (
	"API_MBundestag/help"
	"database/sql"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func compareArticals(t *testing.T, expected Article, input Article) {
	assert.Equal(t, expected.UUID, input.UUID)
	assert.Equal(t, expected.Publication, input.Publication)
	assert.Equal(t, expected.Written.Format("2006-01-02T15:04:05"), input.Written.Format("2006-01-02T15:04:05"))
	assert.Equal(t, expected.Author, input.Author)
	assert.Equal(t, expected.Flair, input.Flair)
	assert.Equal(t, expected.Headline, input.Headline)
	assert.Equal(t, expected.Subtitle, input.Subtitle)
	assert.Equal(t, expected.Content, input.Content)
}

func TestNews(t *testing.T) {
	helper.TestBlocker.Lock()
	defer helper.TestBlocker.Unlock()
	TestNewsDB()

	t.Run("testCreatePublication", testCreatePublication)
	t.Run("testGetPublicationByID", testGetPublicationByID)
	t.Run("testGetChangePublication", testGetChangePublication)
	t.Run("testCreateArticle", testCreateArticle)
	t.Run("testGetArticleByID", testGetArticleByID)
	t.Run("testCompareSingleArticle", testCompareSingleArticle)
	t.Run("testOnlyAssociatedArticles", testOnlyAssociatedArticles)
	t.Run("testCompareCorrectArticleOrder", testCompareCorrectArticleOrder)
	t.Run("testChangeArticle", testChangeArticle)
	t.Run("testDeleteArticle", testDeleteArticle)
	t.Run("testUnpublishedList", testUnpublishedList)
	t.Run("testDeletePublication", testDeletePublication)
	t.Run("createNeededPublications", createNeededPublications)
	t.Run("testPublicationBeforeDate", testPublicationBeforeDate)
	t.Run("testPublicationAfterUUID", testPublicationAfterUUID)
	t.Run("testPublicationBeforeUUID", testPublicationBeforeUUID)
}

func testPublicationBeforeUUID(t *testing.T) {
	list := PublicationList{}
	err, ex := list.GetPublicationBefore("2", 2)
	assert.Nil(t, err)
	assert.True(t, ex)
	assert.Equal(t, 2, len(list))
	assert.Equal(t, "2022-12-14", list[0].PublishTime.Format("2006-01-02"))
	assert.Equal(t, "2022-12-13", list[1].PublishTime.Format("2006-01-02"))
	err, ex = list.GetPublicationBefore("4", 2)
	assert.Nil(t, err)
	assert.True(t, ex)
	assert.Equal(t, 1, len(list))
	assert.Equal(t, "2022-12-15", list[0].PublishTime.Format("2006-01-02"))
	err, ex = list.GetPublicationBefore("", 2)
	assert.Nil(t, err)
	assert.False(t, ex)
	assert.Equal(t, 0, len(list))
}

func testPublicationAfterUUID(t *testing.T) {
	list := PublicationList{}
	err, ex := list.GetPublicationAfter("4", 2)
	assert.Nil(t, err)
	assert.True(t, ex)
	assert.Equal(t, 2, len(list))
	assert.Equal(t, "2022-12-13", list[0].PublishTime.Format("2006-01-02"))
	assert.Equal(t, "2022-12-12", list[1].PublishTime.Format("2006-01-02"))
	err, ex = list.GetPublicationAfter("2", 2)
	assert.True(t, ex)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(list))
	assert.Equal(t, "2022-12-11", list[0].PublishTime.Format("2006-01-02"))
	err, ex = list.GetPublicationAfter("", 2)
	assert.False(t, ex)
	assert.Nil(t, err)
	assert.Equal(t, 2, len(list))
	assert.Equal(t, "2022-12-15", list[0].PublishTime.Format("2006-01-02"))
	assert.Equal(t, "2022-12-14", list[1].PublishTime.Format("2006-01-02"))
}

func testPublicationBeforeDate(t *testing.T) {
	ti, _ := time.Parse("2006-01-02", "2022-12-13")
	list := PublicationList{}
	err := list.GetPublicationBeforeDate(ti, 4)
	assert.Nil(t, err)
	assert.Equal(t, 2, len(list))
	assert.Equal(t, "2022-12-12", list[0].PublishTime.Format("2006-01-02"))
	assert.Equal(t, "2022-12-11", list[1].PublishTime.Format("2006-01-02"))
	err = list.GetPublicationBeforeDate(ti, 1)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(list))
	assert.Equal(t, "2022-12-12", list[0].PublishTime.Format("2006-01-02"))
}

func createNeededPublications(t *testing.T) {
	ti, _ := time.Parse("2006-01-02", "2022-12-11")
	pubGeneral := Publication{
		UUID:         "1",
		PublishTime:  ti,
		Publicated:   true,
		BreakingNews: false,
	}
	err := pubGeneral.CreateMe()
	assert.Nil(t, err)

	pubGeneral.UUID = "2"
	ti, _ = time.Parse("2006-01-02", "2022-12-12")
	pubGeneral.PublishTime = ti
	err = pubGeneral.CreateMe()
	assert.Nil(t, err)

	pubGeneral.UUID = "3"
	ti, _ = time.Parse("2006-01-02", "2022-12-13")
	pubGeneral.PublishTime = ti
	err = pubGeneral.CreateMe()
	assert.Nil(t, err)

	pubGeneral.UUID = "4"
	ti, _ = time.Parse("2006-01-02", "2022-12-14")
	pubGeneral.PublishTime = ti
	err = pubGeneral.CreateMe()
	assert.Nil(t, err)

	pubGeneral.UUID = "5"
	ti, _ = time.Parse("2006-01-02", "2022-12-15")
	pubGeneral.PublishTime = ti
	err = pubGeneral.CreateMe()
	assert.Nil(t, err)
}

func testDeletePublication(t *testing.T) {
	pubTest := Publication{}
	err := pubTest.GetByID("lol")
	assert.Nil(t, err)
	err = pubTest.DeleteMe()
	assert.Nil(t, err)
	err = pubTest.GetByID("lol")
	assert.Equal(t, sql.ErrNoRows, err)
}

func testUnpublishedList(t *testing.T) {
	pub = Publication{
		UUID:         "lol",
		CreateTime:   time.Now(),
		PublishTime:  time.Now(),
		Publicated:   false,
		BreakingNews: false,
	}
	err := pub.CreateMe()
	assert.Nil(t, err)

	list := PublicationList{}
	err = list.GetOnlyUnpublicated()
	assert.Nil(t, err)
	assert.Equal(t, 3, len(list))

	assert.Equal(t, EternatityPublicationName, list[0].UUID)
	assert.Equal(t, "test", list[1].UUID)
	assert.Equal(t, "lol", list[2].UUID)
}

func testDeleteArticle(t *testing.T) {
	art2 := Article{}
	err := art2.GetByID("newUUID")
	assert.Nil(t, err)
	err = art2.DeleteMe()
	assert.Nil(t, err)
	err = art2.GetByID("newUUID")
	assert.Equal(t, sql.ErrNoRows, err)
}

func testChangeArticle(t *testing.T) {
	art2 := Article{}
	err := art2.GetByID("newUUID")
	assert.Nil(t, err)
	art2.Subtitle.String = "asdbkasd"
	art2.Written = time.Now()
	err = art2.SaveChanges()
	assert.Nil(t, err)
	res := Article{}
	err = res.GetByID("newUUID")
	assert.Nil(t, err)
	compareArticals(t, art2, res)
}

func testCompareCorrectArticleOrder(t *testing.T) {
	art2 := Article{
		UUID:        "trhice",
		Publication: "test",
		Author:      "bazinga",
		Flair:       "lol",
		Headline:    "headline",
		Subtitle: sql.NullString{
			String: "subtitle",
			Valid:  true,
		},
		Content: "content goes here",
	}
	err := art2.CreateMe()
	assert.Nil(t, err)

	list := ArticleList{}
	err = list.GetAllArticlesToPublication("test")
	assert.Nil(t, err)
	assert.Equal(t, 2, len(list))
	assert.Equal(t, "test", list[0].UUID)
	assert.Equal(t, "trhice", list[1].UUID)
}

func testOnlyAssociatedArticles(t *testing.T) {
	art2 := Article{
		UUID:        "newUUID",
		Publication: "test2",
		Author:      "bazinga",
		Flair:       "lol",
		Headline:    "headline",
		Subtitle: sql.NullString{
			String: "subtitle",
			Valid:  true,
		},
		Content: "content goes here",
	}
	err := art2.CreateMe()
	assert.Nil(t, err)

	testCompareSingleArticle(t)
}

func testCompareSingleArticle(t *testing.T) {
	list := ArticleList{}
	err := list.GetAllArticlesToPublication("test")
	assert.Nil(t, err)

	assert.Equal(t, 1, len(list))
	res := list[0]
	compareArticals(t, art, res)
}

func testGetArticleByID(t *testing.T) {
	res := Article{}
	err := res.GetByID("test")
	assert.Nil(t, err)
	compareArticals(t, art, res)
}

var pub Publication
var art Article
var array ArticleList

func testCreateArticle(t *testing.T) {
	art = Article{
		UUID:        "test",
		Publication: "test",
		Written:     time.Time{},
		Author:      "bazinga",
		Flair:       "lol",
		Headline:    "headline",
		Subtitle: sql.NullString{
			String: "subtitle",
			Valid:  true,
		},
		Content: "content goes here",
	}
	err := art.CreateMe()
	art.Written = time.Now()
	array = append(array, art)
	assert.Nil(t, err)
}

func testGetChangePublication(t *testing.T) {
	res := Publication{}
	err := res.GetByID("test")
	assert.Nil(t, err)

	res.PublishTime, _ = time.Parse("2006-01-02T15:04:05", "2006-01-02T15:04:05")
	err = res.SaveChanges()
	assert.Nil(t, err)

	assert.Equal(t, pub.UUID, res.UUID)
	assert.Equal(t, pub.CreateTime.Format("2006-01-02T15:04:05"), res.CreateTime.Format("2006-01-02T15:04:05"))
	assert.NotEqual(t, pub.PublishTime.Format("2006-01-02T15:04:05"), res.PublishTime.Format("2006-01-02T15:04:05"))
	assert.Equal(t, pub.Publicated, res.Publicated)
	assert.Equal(t, pub.BreakingNews, res.BreakingNews)

	pub = res
	err = res.GetByID("test")
	assert.Nil(t, err)

	assert.Equal(t, pub.UUID, res.UUID)
	assert.Equal(t, pub.CreateTime.Format("2006-01-02T15:04:05"), res.CreateTime.Format("2006-01-02T15:04:05"))
	assert.Equal(t, pub.PublishTime.Format("2006-01-02T15:04:05"), res.PublishTime.Format("2006-01-02T15:04:05"))
	assert.Equal(t, pub.Publicated, res.Publicated)
	assert.Equal(t, pub.BreakingNews, res.BreakingNews)
}

func testGetPublicationByID(t *testing.T) {
	res := Publication{}
	err := res.GetByID("test")
	assert.Nil(t, err)
	assert.Equal(t, pub.UUID, res.UUID)
	assert.Equal(t, pub.CreateTime.Format("2006-01-02T15:04:05"), res.CreateTime.Format("2006-01-02T15:04:05"))
	assert.Equal(t, pub.PublishTime.Format("2006-01-02T15:04:05"), res.PublishTime.Format("2006-01-02T15:04:05"))
	assert.Equal(t, pub.Publicated, res.Publicated)
	assert.Equal(t, pub.BreakingNews, res.BreakingNews)
}

func testCreatePublication(t *testing.T) {
	pub = Publication{
		UUID:         "test",
		CreateTime:   time.Now(),
		PublishTime:  time.Now(),
		Publicated:   false,
		BreakingNews: false,
	}

	err := pub.CreateMe()
	assert.Nil(t, err)
}
