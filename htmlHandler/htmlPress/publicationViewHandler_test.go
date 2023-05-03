package htmlPress

import (
	"API_MBundestag/database_old"
	"API_MBundestag/help/generics"
	"API_MBundestag/htmlHandler"
	"API_MBundestag/htmlHandler/htmlBasics"
	wr "API_MBundestag/htmlWrapper"
	"database/sql"
	"github.com/stretchr/testify/assert"
	"html/template"
	"testing"
	"time"
)

func TestPublicationViewHandler(t *testing.T) {
	helper.TestBlocker.Lock()
	defer helper.TestBlocker.Unlock()
	Setup()
	htmlBasics.Setup()

	database.TestNewsDB()
	database.TestAccountDB()

	t.Run("setupTestPublicationViewHandler", setupTestPublicationViewHandler)
	t.Run("testGetPublicationViewPage", testGetPublicationViewPage)
	t.Run("testPostPublicationViewPage", testPostPublicationViewPage)
}

func testPostPublicationViewPage(t *testing.T) {
	w, ctx := htmlHandler.GetEmptyContext(t)
	PostPublicationViewPage(ctx)
	assert.Equal(t, "error "+generics.NotAuthorizedToView, w.Body.String())

	acc := database.Account{}
	err := acc.GetByUserName("admin")
	assert.Nil(t, err)

	w, ctx = htmlHandler.GetContextSetForUserWithForm(t, acc, map[string]string{"uuid": "test"})
	PostPublicationViewPage(ctx)
	assert.Equal(t, "error "+generics.NewsIsAlreadyPublished, w.Body.String())

	w, ctx = htmlHandler.GetContextSetForUserWithForm(t, acc, map[string]string{"uuid": "test3"})
	PostPublicationViewPage(ctx)
	assert.Equal(t, "/publication?uuid=test3", w.Header().Get("Location"))
}

func testGetPublicationViewPage(t *testing.T) {
	acc := database.Account{
		DisplayName: "admin",
		Username:    "admin",
		Role:        database.MediaAdmin,
	}
	err := acc.CreateMe()
	assert.Nil(t, err)

	art := database.Article{Publication: database.EternatityPublicationName, UUID: "test", Headline: "lol"}
	err = art.CreateMe()
	assert.Nil(t, err)
	art = database.Article{Publication: database.EternatityPublicationName, UUID: "test2", Headline: "lol2"}
	err = art.CreateMe()
	assert.Nil(t, err)

	w, ctx := htmlHandler.GetEmptyContext(t)
	ctx.Request.URL.RawQuery = "uuid=" + database.EternatityPublicationName
	GetPublicationViewPage(ctx)
	assert.Equal(t, "error "+generics.PublicationDoesNotExistsOrNotAllowedToView, w.Body.String())

	w, ctx = htmlHandler.GetContextSetForUser(t, acc)
	ctx.Request.URL.RawQuery = "uuid=" + database.EternatityPublicationName
	GetPublicationViewPage(ctx)
	assert.Equal(t, "viewPublication "+database.EternatityPublicationName+"  lol  lol2  "+generics.FormatHiddenNormalNews+generics.FormatTimeForArticle, w.Body.String())

	pub := database.Publication{
		UUID:         "test",
		PublishTime:  time.Now(),
		Publicated:   true,
		BreakingNews: false,
	}
	err = pub.CreateMe()
	assert.Nil(t, err)

	w, ctx = htmlHandler.GetEmptyContext(t)
	ctx.Request.URL.RawQuery = "uuid=test"
	GetPublicationViewPage(ctx)
	assert.Equal(t, "viewPublication test  "+generics.FormatNormalNews+generics.FormatTimeForArticle, w.Body.String())

	pub = database.Publication{
		UUID:         "test2",
		PublishTime:  time.Now(),
		Publicated:   true,
		BreakingNews: true,
	}
	err = pub.CreateMe()
	assert.Nil(t, err)

	w, ctx = htmlHandler.GetEmptyContext(t)
	ctx.Request.URL.RawQuery = "uuid=test2"
	GetPublicationViewPage(ctx)
	assert.Equal(t, "viewPublication test2  "+generics.FormatBreakingNews+generics.FormatTimeForArticle, w.Body.String())

	pub = database.Publication{
		UUID:         "test3",
		PublishTime:  time.Now(),
		Publicated:   false,
		BreakingNews: true,
	}
	err = pub.CreateMe()
	assert.Nil(t, err)

	w, ctx = htmlHandler.GetContextSetForUser(t, acc)
	ctx.Request.URL.RawQuery = "uuid=test3"
	GetPublicationViewPage(ctx)
	assert.Equal(t, "viewPublication test3  "+generics.FormatHiddenBreakingNews+generics.FormatTimeForArticle, w.Body.String())
}

func setupTestPublicationViewHandler(t *testing.T) {
	temp, err := template.New("layout").Funcs(wr.DefaultFunctions).Parse("{{.Template}} {{.Page.Error}}")
	assert.Nil(t, err)
	temp2, err := template.New("layout").Funcs(wr.DefaultFunctions).Parse("{{.Template}} {{.Page.Publication.UUID}} {{range $i, $a := .Page.Articles}} {{$a.Headline}} {{end}} {{.Page.FormatNewspaper}}{{.Page.FormatArticle}}")
	assert.Nil(t, err)
	htmlHandler.Template = &wr.Templates{
		Extension: "",
		Dir:       "",
		Templates: map[string]*template.Template{
			"error":           temp,
			"viewPublication": temp2,
		},
	}
}

func TestValidatePublishNewspaper(t *testing.T) {
	helper.TestBlocker.Lock()
	defer helper.TestBlocker.Unlock()
	Setup()
	htmlBasics.Setup()

	database.TestNewsDB()

	t.Run("testFailGetPublication", testFailGetPublication)
	t.Run("testFailPublicationAlreadyPublicated", testFailPublicationAlreadyPublicated)
	t.Run("testPublishNormalNewspaper", testPublishNormalNewspaper)
	t.Run("testPublishBreakingNews", testPublishBreakingNews)
}

func testPublishBreakingNews(t *testing.T) {
	pub := database.Publication{
		UUID:         "test2",
		CreateTime:   time.Time{},
		PublishTime:  time.Now(),
		Publicated:   false,
		BreakingNews: true,
	}
	err := pub.CreateMe()
	assert.Nil(t, err)

	_, ctx := htmlHandler.GetContextWithForm(t, map[string]string{"uuid": "test2"})
	bo := getPublication(ctx, &pub)
	assert.False(t, bo)
	id := ""
	bo = publishIfNotNormalNewspaper(&id, &pub)
	assert.False(t, bo)
	assert.Equal(t, "test2", id)

	err = pub.GetByID("test2")
	assert.Nil(t, err)
	assert.Equal(t, true, pub.Publicated)
}

func testPublishNormalNewspaper(t *testing.T) {
	article := database.Article{
		UUID:        "test",
		Publication: database.EternatityPublicationName,
		Author:      "a",
		Flair:       "a",
		Headline:    "a",
		Subtitle:    sql.NullString{Valid: true, String: "test"},
		Content:     "llols",
		HTMLContent: "asdv",
	}
	err := article.CreateMe()
	assert.Nil(t, err)

	_, ctx := htmlHandler.GetContextWithForm(t, map[string]string{"uuid": database.EternatityPublicationName})
	pub := database.Publication{}
	bo := getPublication(ctx, &pub)
	assert.False(t, bo)
	id := ""
	bo = publishIfNotNormalNewspaper(&id, &pub)
	assert.False(t, bo)
	bo = createNewPublicationIfNormalNewspaper(&id, &pub)
	assert.False(t, bo)
	assert.NotEqual(t, "", id)

	err = pub.GetByID(id)
	assert.Nil(t, err)
	assert.Equal(t, true, pub.Publicated)
	assert.Equal(t, false, pub.BreakingNews)

	err = article.GetByID("test")
	assert.Nil(t, err)
	assert.Equal(t, pub.UUID, article.Publication)
	assert.Equal(t, "a", article.Headline)
}

func testFailPublicationAlreadyPublicated(t *testing.T) {
	pub := database.Publication{
		UUID:         "test",
		CreateTime:   time.Time{},
		PublishTime:  time.Now(),
		Publicated:   true,
		BreakingNews: false,
	}
	err := pub.CreateMe()
	assert.Nil(t, err)

	_, ctx := htmlHandler.GetContextWithForm(t, map[string]string{"uuid": "test"})
	bo := getPublication(ctx, &pub)
	assert.False(t, bo)
	assert.True(t, pub.Publicated)

}

func testFailGetPublication(t *testing.T) {
	pub := &database.Publication{}
	_, ctx := htmlHandler.GetEmptyContext(t)
	bo := getPublication(ctx, pub)
	assert.True(t, bo)
}
