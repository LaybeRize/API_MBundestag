package htmlPress

import (
	"API_MBundestag/database_old"
	gen "API_MBundestag/help/generics"
	"API_MBundestag/htmlHandler"
	"API_MBundestag/htmlHandler/htmlBasics"
	wr "API_MBundestag/htmlWrapper"
	"github.com/stretchr/testify/assert"
	"html/template"
	"testing"
	"time"
)

func TestNewpaperListPageRequests(t *testing.T) {
	helper.TestBlocker.Lock()
	defer helper.TestBlocker.Unlock()
	Setup()
	htmlBasics.Setup()

	database.TestAccountDB()
	database.TestNewsDB()

	t.Run("setupValidateNewsPages", setupValidateNewsPages)
	t.Run("additionSetupForNewsRequests", additionSetupForNewsRequests)
	t.Run("testUnpublishedPage", testUnpublishedPage)
	t.Run("testPublishedPage", testPublishedPage)
}

func testPublishedPage(t *testing.T) {
	w, ctx := htmlHandler.GetEmptyContext(t)
	GetNewsPaperListPage(ctx)
	assert.Equal(t, "newspaperList true false false   20  test5 true false test4 true true test3 true false test2 true true test1 true true", w.Body.String())

	w, ctx = htmlHandler.GetEmptyContext(t)
	ctx.Request.URL.RawQuery = "amount=2"
	GetNewsPaperListPage(ctx)
	assert.Equal(t, "newspaperList true true false test4  2  test5 true false test4 true true", w.Body.String())
}

func testUnpublishedPage(t *testing.T) {
	w, ctx := htmlHandler.GetEmptyContext(t)
	GetNewsPaperHiddenListPage(ctx)
	assert.Equal(t, "error "+gen.NotAuthorizedToView, w.Body.String())

	acc := database.Account{}
	err := acc.GetByUserName("admin")
	assert.Nil(t, err)
	w, ctx = htmlHandler.GetContextSetForUser(t, acc)
	GetNewsPaperHiddenListPage(ctx)
	assert.Equal(t, "newspaperList false false false   0  "+database.EternatityPublicationName+" false false hidden1 false true hidden2 false true", w.Body.String())
}

func additionSetupForNewsRequests(t *testing.T) {
	acc := database.Account{
		DisplayName: "admin",
		Username:    "admin",
		Role:        database.Admin,
	}
	err := acc.CreateMe()
	assert.Nil(t, err)

	pub := database.Publication{
		UUID:         "hidden1",
		PublishTime:  time.Now(),
		Publicated:   false,
		BreakingNews: true,
	}
	err = pub.CreateMe()
	assert.Nil(t, err)
	time.Sleep(10)
	pub.PublishTime, pub.UUID = time.Now(), "hidden2"
	err = pub.CreateMe()
	assert.Nil(t, err)
	time.Sleep(10)

	temp, err := template.New("layout").Funcs(wr.DefaultFunctions).Parse("{{.Template}} {{.Page.Error}}")
	assert.Nil(t, err)
	temp2, err := template.New("layout").Funcs(wr.DefaultFunctions).Parse("{{.Template}} {{.Page.Search}} {{.Page.HasNext}} {{.Page.HasBefore}} {{.Page.NextUUID}} {{.Page.BeforeUUID}} {{.Page.Amount}} {{range $i, $pub := .Page.PubList}} {{$pub.UUID}} {{$pub.Publicated}} {{$pub.BreakingNews}}{{end}}")
	assert.Nil(t, err)
	htmlHandler.Template = &wr.Templates{
		Extension: "",
		Dir:       "",
		Templates: map[string]*template.Template{
			"error":         temp,
			"newspaperList": temp2,
		},
	}
}

func TestValidationOfNewsPages(t *testing.T) {
	helper.TestBlocker.Lock()
	defer helper.TestBlocker.Unlock()
	Setup()
	htmlBasics.Setup()

	database.TestNewsDB()

	t.Run("setupValidateNewsPages", setupValidateNewsPages)
	t.Run("testBeforeNewsPage", testBeforeNewsPage)
	t.Run("testAfterNewsPage", testAfterNewsPage)
}

func testAfterNewsPage(t *testing.T) {
	_, ctx := htmlHandler.GetEmptyContext(t)
	ctx.Request.URL.RawQuery = "uuid=test4"
	res := &NewspaperListViewStruct{}
	err := res.validateNewsPaperReadNextPage(ctx, 10)
	assert.Nil(t, err)
	assert.Equal(t, true, res.HasBefore)
	assert.Equal(t, false, res.HasNext)
	assert.Equal(t, "test3", res.BeforeUUID)
	assert.Equal(t, "", res.NextUUID)
	assert.Equal(t, 3, len(res.PubList))
	assert.Equal(t, "test3", res.PubList[0].UUID)
	assert.Equal(t, "test2", res.PubList[1].UUID)
	assert.Equal(t, "test1", res.PubList[2].UUID)

	_, ctx = htmlHandler.GetEmptyContext(t)
	res = &NewspaperListViewStruct{}
	err = res.validateNewsPaperReadNextPage(ctx, 10)
	assert.Nil(t, err)
	assert.Equal(t, false, res.HasBefore)
	assert.Equal(t, false, res.HasNext)
	assert.Equal(t, "", res.BeforeUUID)
	assert.Equal(t, "", res.NextUUID)
	assert.Equal(t, 5, len(res.PubList))
	assert.Equal(t, "test5", res.PubList[0].UUID)
	assert.Equal(t, "test1", res.PubList[4].UUID)

}

func testBeforeNewsPage(t *testing.T) {
	_, ctx := htmlHandler.GetEmptyContext(t)
	ctx.Request.URL.RawQuery = "uuid=test4"
	res := &NewspaperListViewStruct{}
	err := res.validateNewsPaperReadPageBefore(ctx, 10)
	assert.Nil(t, err)
	assert.Equal(t, false, res.HasBefore)
	assert.Equal(t, true, res.HasNext)
	assert.Equal(t, "", res.BeforeUUID)
	assert.Equal(t, "test5", res.NextUUID)
	assert.Equal(t, 1, len(res.PubList))
	assert.Equal(t, "test5", res.PubList[0].UUID)

	_, ctx = htmlHandler.GetEmptyContext(t)
	res = &NewspaperListViewStruct{}
	err = res.validateNewsPaperReadPageBefore(ctx, 10)
	assert.Nil(t, err)
	assert.Equal(t, false, res.HasBefore)
	assert.Equal(t, false, res.HasNext)
	assert.Equal(t, "", res.BeforeUUID)
	assert.Equal(t, "", res.NextUUID)
	assert.Equal(t, 0, len(res.PubList))
}

func setupValidateNewsPages(t *testing.T) {
	pub := database.Publication{
		UUID:         "test1",
		PublishTime:  time.Now(),
		Publicated:   true,
		BreakingNews: true,
	}
	err := pub.CreateMe()
	assert.Nil(t, err)
	time.Sleep(10)
	pub.PublishTime, pub.UUID = time.Now(), "test2"
	err = pub.CreateMe()
	assert.Nil(t, err)
	time.Sleep(10)
	pub.BreakingNews = false
	pub.PublishTime, pub.UUID = time.Now(), "test3"
	err = pub.CreateMe()
	assert.Nil(t, err)
	time.Sleep(10)
	pub.BreakingNews = true
	pub.PublishTime, pub.UUID = time.Now(), "test4"
	err = pub.CreateMe()
	assert.Nil(t, err)
	time.Sleep(10)
	pub.BreakingNews = false
	pub.PublishTime, pub.UUID = time.Now(), "test5"
	err = pub.CreateMe()
	assert.Nil(t, err)
}
