package htmlPress

import (
	"API_MBundestag/dataLogic"
	"API_MBundestag/database_old"
	"API_MBundestag/help/generics"
	"API_MBundestag/htmlHandler"
	"API_MBundestag/htmlHandler/htmlBasics"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"net/url"
	"time"
)

type PublicationViewStruct struct {
	Publication     database.Publication
	Articles        database.ArticleList
	FormatNewspaper string
	FormatArticle   string
}

func GetPublicationViewPage(c *gin.Context) {
	acc, b := dataLogic.CheckUserPrivileged(c, database.HeadAdmin, database.Admin, database.MediaAdmin)

	pub := database.Publication{}
	err := pub.GetByID(c.Query("uuid"))
	if err != nil || (!pub.Publicated && !b) {
		htmlBasics.MakeErrorPage(c, &acc, generics.PublicationDoesNotExistsOrNotAllowedToView)
		return
	}

	articleList := database.ArticleList{}
	err = articleList.GetAllArticlesToPublication(pub.UUID)
	if err != nil {
		htmlBasics.MakeErrorPage(c, &acc, generics.ErrorWhileLoadingArticles)
		return
	}

	pubStruct := &PublicationViewStruct{
		Publication:   pub,
		Articles:      articleList,
		FormatArticle: generics.FormatTimeForArticle,
	}

	switch true {
	case pub.Publicated && pub.BreakingNews:
		pubStruct.FormatNewspaper = generics.FormatBreakingNews
	case pub.Publicated && !pub.BreakingNews:
		pubStruct.FormatNewspaper = generics.FormatNormalNews
	case !pub.Publicated && pub.BreakingNews:
		pubStruct.FormatNewspaper = generics.FormatHiddenBreakingNews
	case !pub.Publicated && !pub.BreakingNews:
		pubStruct.FormatNewspaper = generics.FormatHiddenNormalNews
	}

	htmlHandler.MakeSite(pubStruct, c, &acc)
}

func PostPublicationViewPage(c *gin.Context) {
	acc, b := dataLogic.CheckUserPrivileged(c, database.HeadAdmin, database.Admin, database.MediaAdmin)
	if !b {
		htmlBasics.MakeErrorPage(c, &acc, generics.NotAuthorizedToView)
		return
	}

	id := ""
	pub := &database.Publication{}
	switch true {
	case getPublication(c, pub):
		htmlBasics.MakeErrorPage(c, &acc, generics.ErrorWhilePublishingNews)
	case pub.Publicated:
		htmlBasics.MakeErrorPage(c, &acc, generics.NewsIsAlreadyPublished)
	case publishIfNotNormalNewspaper(&id, pub):
		htmlBasics.MakeErrorPage(c, &acc, generics.ErrorWhilePublishingNews)
	case createNewPublicationIfNormalNewspaper(&id, pub):
		htmlBasics.MakeErrorPage(c, &acc, generics.ErrorWhilePublishingNews)
	default:
		c.Redirect(http.StatusFound, "/publication?uuid="+url.QueryEscape(id))
	}
}

func createNewPublicationIfNormalNewspaper(id *string, pub *database.Publication) bool {
	if *id != "" {
		return false
	}

	pub = &database.Publication{
		UUID:         uuid.New().String(),
		PublishTime:  time.Now().UTC(),
		Publicated:   true,
		BreakingNews: false,
	}
	err := pub.CreateMe()
	if err != nil {
		return true
	}

	*id = pub.UUID

	list := database.ArticleList{}
	err = list.GetAllArticlesToPublication(database.EternatityPublicationName)
	if err != nil {
		return true
	}

	for _, art := range list {
		art.Publication = pub.UUID
		err = art.SaveChanges()
		if err != nil {
			return true
		}
	}

	return false
}

func publishIfNotNormalNewspaper(id *string, pub *database.Publication) bool {
	if pub.UUID == database.EternatityPublicationName {
		return false
	}

	*id = pub.UUID
	pub.Publicated = true
	pub.PublishTime = time.Now().UTC()
	err := pub.SaveChanges()
	if err != nil {
		return true
	}
	return false
}

func getPublication(c *gin.Context, pub *database.Publication) bool {
	err := pub.GetByID(generics.GetText(c, "uuid"))
	if err != nil {
		return true
	}
	return false
}
