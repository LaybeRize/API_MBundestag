package htmlPress

import (
	"API_MBundestag/dataLogic"
	"API_MBundestag/database"
	"API_MBundestag/help"
	"API_MBundestag/help/generics"
	"API_MBundestag/htmlHandler"
	"API_MBundestag/htmlHandler/htmlBasics"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

type RejectArticleStruct struct {
	Article           database.Article
	ArticleFormat     string
	ModMessageContent string
	Message           string
}

func GetRejectArticlePage(c *gin.Context) {
	acc, b := dataLogic.CheckUserPrivileged(c, database.HeadAdmin, database.Admin, database.MediaAdmin)
	if !b {
		htmlBasics.MakeErrorPage(c, &acc, generics.NotAuthorizedToView)
		return
	}

	rejectStruct := &RejectArticleStruct{
		Article:       database.Article{},
		ArticleFormat: generics.FormatTimeForArticle,
	}
	err := rejectStruct.Article.GetByID(c.Query("uuid"))
	if err != nil {
		rejectStruct.Message = generics.CanNotFindArticle
	}

	htmlHandler.MakeSite(rejectStruct, c, &acc)
}

func PostRejectArticlePage(c *gin.Context) {
	acc, b := dataLogic.CheckUserPrivileged(c, database.HeadAdmin, database.Admin, database.MediaAdmin)
	if !b {
		htmlBasics.MakeErrorPage(c, &acc, generics.NotAuthorizedToView)
		return
	}

	rejectStruct := &RejectArticleStruct{
		Article:           database.Article{},
		ArticleFormat:     generics.FormatTimeForArticle,
		ModMessageContent: generics.GetText(c, "content"),
	}

	success := rejectStruct.validateRejection(c)
	if success && rejectStruct.Article.Publication == database.EternatityPublicationName {
		c.Redirect(http.StatusFound, "/publication?uuid="+database.EternatityPublicationName)
		return
	}
	if success {
		c.Redirect(http.StatusFound, "/newspaper-approval")
		return
	}
	htmlHandler.MakeSite(rejectStruct, c, &acc)
}

func (rejectStruct *RejectArticleStruct) validateRejection(c *gin.Context) bool {

	pub := &database.Publication{}
	switch false {
	case rejectStruct.getArticle(c):
	case rejectStruct.getPublcation(pub):
	case rejectStruct.createRejectionLetter(c):
	case rejectStruct.deleteArticle():
	case rejectStruct.deletePublicationIfNeeded(pub):
	default:
		return true
	}
	return false
}

func (rejectStruct *RejectArticleStruct) getArticle(c *gin.Context) bool {
	err := rejectStruct.Article.GetByID(generics.GetText(c, "uuid"))
	if err != nil {
		rejectStruct.Message = generics.CanNotFindArticle
		return false
	}
	return true
}

func (rejectStruct *RejectArticleStruct) getPublcation(pub *database.Publication) bool {
	err := pub.GetByID(rejectStruct.Article.Publication)
	if err != nil {
		rejectStruct.Message = generics.CanNotFindPublicationForArticle
		return false
	}
	if pub.Publicated {
		rejectStruct.Message = generics.ArticleAlreadyPublished
		return false
	}
	return true
}

func (rejectStruct *RejectArticleStruct) createRejectionLetter(c *gin.Context) bool {
	letter := database.Letter{
		UUID:    uuid.New().String(),
		Author:  generics.AuthorQualityCheck,
		Title:   fmt.Sprintf(generics.LetterRejectTitle, rejectStruct.Article.Headline),
		Content: fmt.Sprintf(generics.LetterRejectText, rejectStruct.Article.Subtitle.String, rejectStruct.Article.Content, generics.GetText(c, "content")),
		Info: database.LetterInfo{
			//ModMessage:          true,
			AllHaveToAgree: false,
			NoSigning:      true,
			//PeopleInvitedToSign: []string{rejectStruct.Article.Author},
			PeopleNotYetSigned: []string{},
			Signed:             []string{},
			Rejected:           []string{},
		},
	}
	letter.HTMLContent = help.CreateHTML(letter.Content)
	err := letter.CreateMe()
	if err != nil {
		rejectStruct.Message = generics.RejectionCouldNotBeCreated
		return false
	}
	return true
}

func (rejectStruct *RejectArticleStruct) deleteArticle() bool {
	err := rejectStruct.Article.DeleteMe()
	if err != nil {
		rejectStruct.Message = generics.CouldNotDeleteArticle
		return false
	}
	return true
}

func (rejectStruct *RejectArticleStruct) deletePublicationIfNeeded(pub *database.Publication) bool {
	if pub.UUID != database.EternatityPublicationName {
		err := pub.DeleteMe()
		if err != nil {
			rejectStruct.Message = generics.CouldNotDeletePublication
			return false
		}
	}
	return true
}
