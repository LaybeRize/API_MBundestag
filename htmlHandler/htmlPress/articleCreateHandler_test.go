package htmlPress

import (
	"API_MBundestag/database_old"
	"API_MBundestag/help/generics"
	"API_MBundestag/htmlHandler"
	"API_MBundestag/htmlHandler/htmlBasics"
	wr "API_MBundestag/htmlWrapper"
	"github.com/stretchr/testify/assert"
	"html/template"
	"testing"
)

func TestCreateArticleHandler(t *testing.T) {
	helper.TestBlocker.Lock()
	defer helper.TestBlocker.Unlock()
	Setup()
	htmlBasics.Setup()

	database.TestAccountDB()
	database.TestNewsDB()

	t.Run("setupAccountsArticles", setupAccountsArticles)
	t.Run("createArticleHandlerPageSetup", createArticleHandlerPageSetup)
	t.Run("testGetCreateArticlePage", testGetCreateArticlePage)
	t.Run("testPostCreateArticlePagePreview", testPostCreateArticlePagePreview)
	t.Run("testPostCreateArticlePage", testPostCreateArticlePage)
}

func testPostCreateArticlePage(t *testing.T) {
	acc := database.Account{}
	err := acc.GetByUserName("test")
	assert.Nil(t, err)
	w, ctx := htmlHandler.GetContextSetForUserWithForm(t, acc, map[string]string{"title": "title", "subtitle": "testSub", "content": "test", "selectedAccount": "press"})
	PostCreateArticlePage(ctx)
	assert.Equal(t, "createArticle  false    "+generics.SuccessfulCreateArticle+"\n", w.Body.String())

}

func testPostCreateArticlePagePreview(t *testing.T) {
	w, ctx := htmlHandler.GetEmptyContext(t)
	PostCreateArticlePage(ctx)
	assert.Equal(t, "error "+generics.NotAuthorizedToView, w.Body.String())

	acc := database.Account{}
	err := acc.GetByUserName("test")
	assert.Nil(t, err)
	w, ctx = htmlHandler.GetContextSetForUserWithForm(t, acc, map[string]string{"title": "title", "subtitle": "testSub", "content": "test", "selectedAccount": "press"})
	ctx.Request.URL.RawQuery = "type=preview"
	PostCreateArticlePage(ctx)
	assert.Equal(t, "createArticle press false title testSub test <p "+helper.ReplacerMap["p"]+">test</p>\n"+generics.PreviewText+"\n", w.Body.String())

}

func testGetCreateArticlePage(t *testing.T) {
	w, ctx := htmlHandler.GetEmptyContext(t)
	GetCreateArticlePage(ctx)
	assert.Equal(t, "error "+generics.NotAuthorizedToView, w.Body.String())

	acc := database.Account{}
	err := acc.GetByUserName("test")
	assert.Nil(t, err)
	w, ctx = htmlHandler.GetContextSetForUserWithForm(t, acc, map[string]string{"title": "title", "subtitle": "testSub", "content": "test", "selectedAccount": "press"})
	GetCreateArticlePage(ctx)
	assert.Equal(t, "createArticle  false    ", w.Body.String())
}

func createArticleHandlerPageSetup(t *testing.T) {
	temp, err := template.New("layout").Funcs(wr.DefaultFunctions).Parse("{{.Template}} {{.Page.Error}}")
	assert.Nil(t, err)
	temp2, err := template.New("layout").Funcs(wr.DefaultFunctions).Parse("{{.Template}} {{.Page.SelectedAccount}} {{.Page.BreakingNews}} {{.Page.Article.Headline}} {{.Page.Article.Subtitle.String}} {{.Page.Article.Content}} {{.Page.Preview}}{{.Page.Message}}")
	assert.Nil(t, err)
	htmlHandler.Template = &wr.Templates{
		Extension: "",
		Dir:       "",
		Templates: map[string]*template.Template{
			"error":         temp,
			"createArticle": temp2,
		},
	}
	helper.UpdateAttributes()
}

func TestValidateArticleCreate(t *testing.T) {
	helper.TestBlocker.Lock()
	defer helper.TestBlocker.Unlock()
	Setup()
	htmlBasics.Setup()

	database.TestAccountDB()
	database.TestNewsDB()

	t.Run("setupAccountsArticles", setupAccountsArticles)
	t.Run("testTextOrHeadlineAreEmpty", testTextOrHeadlineAreEmpty)
	t.Run("testWriteAccountDoesNotExistInDatabaseAritcleCreate", testWriteAccountDoesNotExistInDatabaseAritcleCreate)
	t.Run("testThatAccountIsNotAllowedToWriteArticleCreate", testThatAccountIsNotAllowedToWriteArticleCreate)
	t.Run("createBreakingNews", createBreakingNews)
	t.Run("createNormalNews", createNormalNews)
}

func createNormalNews(t *testing.T) {
	_, ctx := htmlHandler.GetContextWithForm(t, map[string]string{"title": "title", "subtitle": "testSub", "content": "test", "selectedAccount": "press"})
	acc := database.Account{}
	err := acc.GetByUserName("test")
	assert.Nil(t, err)
	acc2 := database.Account{}
	err = acc2.GetByUserName("press")
	assert.Nil(t, err)

	res := validateCreateArticle(ctx, &acc)
	assert.Equal(t, CreateArticleStruct{
		Accounts: database.AccountList{acc, acc2},
		Message:  generics.SuccessfulCreateArticle + "\n",
	}, *res)

	arts := database.ArticleList{}
	err = arts.GetAllArticlesToPublication(database.EternatityPublicationName)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(arts))
	assert.Equal(t, "title", arts[0].Headline)
}

func createBreakingNews(t *testing.T) {
	_, ctx := htmlHandler.GetContextWithForm(t, map[string]string{"title": "title", "subtitle": "testSub", "content": "test", "selectedAccount": "press", "breakingNews": "true"})
	acc := database.Account{}
	err := acc.GetByUserName("test")
	assert.Nil(t, err)
	acc2 := database.Account{}
	err = acc2.GetByUserName("press")
	assert.Nil(t, err)

	res := validateCreateArticle(ctx, &acc)
	assert.Equal(t, CreateArticleStruct{
		Accounts: database.AccountList{acc, acc2},
		Message:  generics.SuccessfulCreateArticle + "\n",
	}, *res)

	pubs := database.PublicationList{}
	err = pubs.GetOnlyUnpublicated()
	assert.Nil(t, err)
	assert.Equal(t, 2, len(pubs))
	assert.Equal(t, database.EternatityPublicationName, pubs[0].UUID)
	arts := database.ArticleList{}
	err = arts.GetAllArticlesToPublication(pubs[1].UUID)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(arts))
	assert.Equal(t, "title", arts[0].Headline)
}

func testThatAccountIsNotAllowedToWriteArticleCreate(t *testing.T) {
	_, ctx := htmlHandler.GetContextWithForm(t, map[string]string{"title": "title", "content": "test", "selectedAccount": "press2"})
	acc := database.Account{}
	err := acc.GetByUserName("test")
	assert.Nil(t, err)
	acc2 := database.Account{}
	err = acc2.GetByUserName("press")
	assert.Nil(t, err)

	res := validateCreateArticle(ctx, &acc)
	assert.Equal(t, CreateArticleStruct{
		Accounts: database.AccountList{acc, acc2},
		Article: database.Article{
			Headline: "title",
			Content:  "test",
		},
		SelectedAccount: "press2",
		Message:         generics.AccountIsNotYours + "\n",
	}, *res)
}

func testWriteAccountDoesNotExistInDatabaseAritcleCreate(t *testing.T) {
	_, ctx := htmlHandler.GetContextWithForm(t, map[string]string{"title": "title", "content": "test"})
	acc := database.Account{}
	err := acc.GetByUserName("test")
	assert.Nil(t, err)
	acc2 := database.Account{}
	err = acc2.GetByUserName("press")
	assert.Nil(t, err)

	res := validateCreateArticle(ctx, &acc)
	assert.Equal(t, CreateArticleStruct{
		Accounts: database.AccountList{acc, acc2},
		Article: database.Article{
			Headline: "title",
			Content:  "test",
		},
		Message: generics.AccountDoesNotExists + "\n",
	}, *res)
}

func testTextOrHeadlineAreEmpty(t *testing.T) {
	_, ctx := htmlHandler.GetEmptyContext(t)
	acc := database.Account{}
	err := acc.GetByUserName("test")
	assert.Nil(t, err)
	acc2 := database.Account{}
	err = acc2.GetByUserName("press")
	assert.Nil(t, err)

	res := validateCreateArticle(ctx, &acc)
	assert.Equal(t, CreateArticleStruct{
		Accounts: database.AccountList{acc, acc2},
		Message:  generics.TextOrHeadlineAreEmpty + "\n",
	}, *res)
}

func TestGetEmtpyCreateArticleStruct(t *testing.T) {
	helper.TestBlocker.Lock()
	defer helper.TestBlocker.Unlock()
	Setup()
	htmlBasics.Setup()

	database.TestAccountDB()

	t.Run("setupAccountsArticles", setupAccountsArticles)
	t.Run("testGetEmtpyCreateArticleStruct", testGetEmtpyCreateArticleStruct)
}

func testGetEmtpyCreateArticleStruct(t *testing.T) {
	acc := database.Account{}
	err := acc.GetByDisplayName("press")
	assert.Nil(t, err)
	res := getEmtpyCreateArticleStruct(&database.Account{ID: 1})
	assert.Equal(t, CreateArticleStruct{
		Accounts: database.AccountList{database.Account{ID: 1}, acc},
	}, *res)

}

func setupAccountsArticles(t *testing.T) {
	acc := database.Account{
		DisplayName: "test",
		Username:    "test",
		Role:        database.Admin,
	}
	err := acc.CreateMe()
	assert.Nil(t, err)
	acc.DisplayName, acc.Username = "press", "press"
	acc.Role, acc.Linked.Valid, acc.Linked.Int64 = database.PressAccount, true, 1
	err = acc.CreateMe()
	assert.Nil(t, err)
	acc.DisplayName, acc.Username = "press2", "press2"
	acc.Suspended = true
	err = acc.CreateMe()
	assert.Nil(t, err)
	err = acc.SaveChanges()
	assert.Nil(t, err)
}
