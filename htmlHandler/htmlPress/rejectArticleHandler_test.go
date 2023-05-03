package htmlPress

import (
	"API_MBundestag/database_old"
	"API_MBundestag/help/generics"
	"API_MBundestag/htmlHandler"
	"API_MBundestag/htmlHandler/htmlBasics"
	wr "API_MBundestag/htmlWrapper"
	"database/sql"
	"fmt"
	"github.com/stretchr/testify/assert"
	"html/template"
	"strings"
	"testing"
	"time"
)

func TestRejectArticleHandler(t *testing.T) {
	helper.TestBlocker.Lock()
	defer helper.TestBlocker.Unlock()
	Setup()
	htmlBasics.Setup()

	database.TestNewsDB()
	database.TestLettersDB()
	database.TestAccountDB()

	t.Run("setupRejectArticlePage", setupRejectArticlePage)
	t.Run("testGetRejectArticlePage", testGetRejectArticlePage)
	t.Run("testPostRejectArticlePagePreview", testPostRejectArticlePagePreview)
	t.Run("testPostRejectArticlePage", testPostRejectArticlePage)
}

func testPostRejectArticlePage(t *testing.T) {
	pub := database.Publication{
		UUID:         "test",
		Publicated:   false,
		BreakingNews: true,
	}
	err := pub.CreateMe()
	assert.Nil(t, err)

	acc := database.Account{}
	err = acc.GetByUserName("admin")
	assert.Nil(t, err)
	w, ctx := htmlHandler.GetContextSetForUserWithForm(t, acc, map[string]string{"content": "test", "uuid": "test"})
	PostRejectArticlePage(ctx)
	assert.Equal(t, "/newspaper-approval", w.Header().Get("Location"))

	list := database.LetterList{}
	var ex bool
	err, ex = list.GetPublicationAfter("", 10, "test", false)
	assert.False(t, ex)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(list))

	art := database.Article{
		UUID:        "test",
		Publication: database.EternatityPublicationName,
		Author:      "test",
		Flair:       "test",
		Headline:    "test",
		Content:     "test",
	}
	err = art.CreateMe()
	assert.Nil(t, err)

	w, ctx = htmlHandler.GetContextSetForUserWithForm(t, acc, map[string]string{"content": "test", "uuid": "test"})
	PostRejectArticlePage(ctx)
	assert.Equal(t, true, strings.HasPrefix(w.Header().Get("Location"), "/publication?uuid="))

	articles := database.ArticleList{}
	err = articles.GetAllArticlesToPublication(database.EternatityPublicationName)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(articles))

	err, ex = list.GetPublicationAfter("", 10, "test", false)
	assert.False(t, ex)
	assert.Nil(t, err)
	assert.Equal(t, 2, len(list))
}

func testPostRejectArticlePagePreview(t *testing.T) {
	w, ctx := htmlHandler.GetEmptyContext(t)
	PostRejectArticlePage(ctx)
	assert.Equal(t, "error "+generics.NotAuthorizedToView, w.Body.String())

	acc := database.Account{}
	err := acc.GetByUserName("admin")
	assert.Nil(t, err)
	w, ctx = htmlHandler.GetContextSetForUserWithForm(t, acc, map[string]string{"content": "test"})
	ctx.Request.URL.RawQuery = "type=preview"
	PostRejectArticlePage(ctx)
	assert.Equal(t, "rejectArticle   test <p "+helper.ReplacerMap["p"]+">test</p>\n "+generics.CanNotFindArticle+generics.FormatTimeForArticle, w.Body.String())
}

func testGetRejectArticlePage(t *testing.T) {
	art := database.Article{
		UUID:        "test",
		Publication: "test",
		Author:      "test",
		Flair:       "test",
		Headline:    "test",
		Content:     "test",
	}
	err := art.CreateMe()
	assert.Nil(t, err)

	w, ctx := htmlHandler.GetEmptyContext(t)
	GetRejectArticlePage(ctx)
	assert.Equal(t, "error "+generics.NotAuthorizedToView, w.Body.String())

	acc := database.Account{
		DisplayName: "admin",
		Username:    "admin",
		Role:        database.MediaAdmin,
	}
	err = acc.CreateMe()
	assert.Nil(t, err)

	w, ctx = htmlHandler.GetContextSetForUser(t, acc)
	GetRejectArticlePage(ctx)
	assert.Equal(t, "rejectArticle     "+generics.CanNotFindArticle+generics.FormatTimeForArticle, w.Body.String())

	w, ctx = htmlHandler.GetContextSetForUser(t, acc)
	ctx.Request.URL.RawQuery = "uuid=test"
	GetRejectArticlePage(ctx)
	assert.Equal(t, "rejectArticle test test   "+generics.FormatTimeForArticle, w.Body.String())

}

func setupRejectArticlePage(t *testing.T) {
	temp, err := template.New("layout").Funcs(wr.DefaultFunctions).Parse("{{.Template}} {{.Page.Error}}")
	assert.Nil(t, err)
	temp2, err := template.New("layout").Funcs(wr.DefaultFunctions).Parse("{{.Template}} {{.Page.Article.Headline}} {{.Page.Article.Content}} {{.Page.ModMessageContent}} {{.Page.Preview}} {{.Page.Message}}{{.Page.ArticleFormat}}")
	assert.Nil(t, err)
	htmlHandler.Template = &wr.Templates{
		Extension: "",
		Dir:       "",
		Templates: map[string]*template.Template{
			"error":         temp,
			"rejectArticle": temp2,
		},
	}
	helper.UpdateAttributes()
}

func TestValidateRejection(t *testing.T) {
	helper.TestBlocker.Lock()
	defer helper.TestBlocker.Unlock()
	Setup()
	htmlBasics.Setup()

	database.TestNewsDB()
	database.TestLettersDB()

	t.Run("setupRejectArticle", setupRejectArticle)
	t.Run("testCanNotFindArticle", testCanNotFindArticle)
	t.Run("testCanNotFindPublicationForArticle", testCanNotFindPublicationForArticle)
	t.Run("testArticleAlreadyPublished", testArticleAlreadyPublished)
	t.Run("testRejectNormalArticle", testRejectNormalArticle)
	t.Run("testRejectBreakingNewsArticle", testRejectBreakingNewsArticle)
}

func testRejectBreakingNewsArticle(t *testing.T) {
	rejectStruct := RejectArticleStruct{}
	_, ctx := htmlHandler.GetContextWithForm(t, map[string]string{"uuid": "test2", "content": "test"})
	boolean := rejectStruct.validateRejection(ctx)
	assert.Equal(t, true, boolean)

	pub := database.Publication{}
	err := pub.GetByID("test")
	assert.Equal(t, sql.ErrNoRows, err)

	letter := database.LetterList{}
	var ex bool
	err, ex = letter.GetPublicationAfter("", 10, "test2", false)
	assert.False(t, ex)
	assert.Equal(t, 1, len(letter))
	assert.Equal(t, fmt.Sprintf(generics.LetterRejectText, "", "test", "test"), letter[0].Content)

}

func testRejectNormalArticle(t *testing.T) {
	rejectStruct := RejectArticleStruct{}
	_, ctx := htmlHandler.GetContextWithForm(t, map[string]string{"uuid": "test3", "content": "test"})
	boolean := rejectStruct.validateRejection(ctx)
	assert.Equal(t, true, boolean)

	list := database.ArticleList{}
	err := list.GetAllArticlesToPublication(database.EternatityPublicationName)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(list))

	letter := database.LetterList{}
	var ex bool
	err, ex = letter.GetPublicationAfter("", 10, "test3", false)
	assert.False(t, ex)
	assert.Equal(t, 1, len(letter))
	assert.Equal(t, fmt.Sprintf(generics.LetterRejectText, "", "test", "test"), letter[0].Content)
}

func testArticleAlreadyPublished(t *testing.T) {
	pub := database.Publication{
		UUID:         "fail",
		PublishTime:  time.Now(),
		Publicated:   true,
		BreakingNews: false,
	}
	err := pub.CreateMe()
	assert.Nil(t, err)

	rejectStruct := RejectArticleStruct{}
	_, ctx := htmlHandler.GetContextWithForm(t, map[string]string{"uuid": "test"})
	boolean := rejectStruct.validateRejection(ctx)
	assert.Equal(t, false, boolean)
	assert.Equal(t, generics.ArticleAlreadyPublished, rejectStruct.Message)
}

func testCanNotFindPublicationForArticle(t *testing.T) {
	rejectStruct := RejectArticleStruct{}
	_, ctx := htmlHandler.GetContextWithForm(t, map[string]string{"uuid": "test"})
	boolean := rejectStruct.validateRejection(ctx)
	assert.Equal(t, false, boolean)
	assert.Equal(t, generics.CanNotFindPublicationForArticle, rejectStruct.Message)
}

func testCanNotFindArticle(t *testing.T) {
	rejectStruct := RejectArticleStruct{}
	_, ctx := htmlHandler.GetEmptyContext(t)
	boolean := rejectStruct.validateRejection(ctx)
	assert.Equal(t, false, boolean)
	assert.Equal(t, generics.CanNotFindArticle, rejectStruct.Message)
}

func setupRejectArticle(t *testing.T) {
	art := database.Article{
		UUID:        "test",
		Publication: "fail",
		Author:      "test",
		Headline:    "test",
		Content:     "test",
	}
	err := art.CreateMe()
	assert.Nil(t, err)
	art.UUID, art.Publication, art.Author = "test2", "test", "test2"
	err = art.CreateMe()
	art.UUID, art.Publication, art.Author = "test3", database.EternatityPublicationName, "test3"
	err = art.CreateMe()
	pub := database.Publication{
		UUID:         "test",
		PublishTime:  time.Now(),
		Publicated:   false,
		BreakingNews: true,
	}
	err = pub.CreateMe()
	assert.Nil(t, err)
}
