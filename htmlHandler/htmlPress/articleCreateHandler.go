package htmlPress

import (
	"API_MBundestag/dataLogic"
	"API_MBundestag/database_old"
	"API_MBundestag/help/generics"
	"API_MBundestag/htmlHandler"
	gen "API_MBundestag/htmlHandler/generics"
	"API_MBundestag/htmlHandler/htmlBasics"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"time"
)

type CreateArticleStruct struct {
	Accounts        database.AccountList
	SelectedAccount string
	Article         database.Article
	BreakingNews    bool
	Message         string
}

func getEmtpyCreateArticleStruct(acc *database.Account) *CreateArticleStruct {
	val := CreateArticleStruct{}
	htmlHandler.FillOwnAccounts(&val, acc)
	val.BreakingNews = false
	return &val
}

func GetCreateArticlePage(c *gin.Context) {
	acc, b := dataLogic.CheckUserPrivileged(c, database.HeadAdmin, database.Admin, database.MediaAdmin, database.User)
	if !b {
		htmlBasics.MakeErrorPage(c, &acc, generics.NotAuthorizedToView)
		return
	}
	htmlHandler.MakeSite(getEmtpyCreateArticleStruct(&acc), c, &acc)
}

func PostCreateArticlePage(c *gin.Context) {
	acc, b := dataLogic.CheckUserPrivileged(c, database.HeadAdmin, database.Admin, database.MediaAdmin, database.User)
	if !b {
		htmlBasics.MakeErrorPage(c, &acc, generics.NotAuthorizedToView)
		return
	}

	htmlHandler.MakeSite(validateCreateArticle(c, &acc), c, &acc)
}

func validateCreateArticle(c *gin.Context, acc *database.Account) (articleStruct *CreateArticleStruct) {
	articleStruct = getEmtpyCreateArticleStruct(acc)

	articleStruct.Article = database.Article{
		Headline: htmlHandler.GetText(c, "title"),
		Subtitle: gen.GetNullString(c, "subtitle"),
		Content:  htmlHandler.GetText(c, "content"),
	}
	articleStruct.SelectedAccount = htmlHandler.GetText(c, "selectedAccount")
	articleStruct.BreakingNews = htmlHandler.GetBool(c, "breakingNews")

	writer := &database.Account{}
	pub := &database.Publication{}
	switch true {
	case articleStruct.checkIfContentAndHeadlineEmpty():
	case htmlHandler.CheckLengthContentLayer(articleStruct, &articleStruct.Article, generics.ArticleContentLimit):
	case htmlHandler.CheckLengthFieldLayer(articleStruct, &articleStruct.Article, generics.ArticleTitleLimit, "Headline", generics.ArticleLimitError):
	case htmlHandler.CheckLengthSubtitleLayer(articleStruct, &articleStruct.Article, generics.ArticleSubtitleLimit):
	case htmlHandler.CheckWriter(articleStruct, writer, acc):
	case articleStruct.updateFields(writer):
	case articleStruct.makePublication(pub):
	default:
		articleStruct.finishArticleCreation(pub, acc)
	}
	return
}

func (s *CreateArticleStruct) checkIfContentAndHeadlineEmpty() bool {
	//Return instantly if the text or title field is empty
	if s.Article.Content == "" || s.Article.Headline == "" {
		s.Message = generics.TextOrHeadlineAreEmpty + "\n" + s.Message
		return true
	}
	return false
}

func (s *CreateArticleStruct) updateFields(writer *database.Account) bool {
	s.Article.Author = writer.DisplayName
	s.Article.Flair = writer.Flair
	s.Article.UUID = uuid.New().String()
	s.Article.HTMLContent = helper.CreateHTML(s.Article.Content)
	return false
}

func (s *CreateArticleStruct) makePublication(pub *database.Publication) bool {
	if s.BreakingNews {
		pub.PublishTime = time.Now().UTC()
		pub.BreakingNews = true
		pub.UUID = uuid.New().String()
		pub.Publicated = false
		err := pub.CreateMe()
		if err != nil {
			s.Message = generics.ErrorWhileCreatingArticle + "\n" + s.Message
			return true
		}
	} else {
		err := pub.GetByID(database.EternatityPublicationName)
		if err != nil {
			s.Message = generics.ErrorWhileCreatingArticle + "\n" + s.Message
			return true
		}
	}
	return false
}

func (s *CreateArticleStruct) finishArticleCreation(pub *database.Publication, acc *database.Account) {
	s.Article.Publication = pub.UUID
	err := s.Article.CreateMe()
	if err != nil {
		s.Message = generics.ErrorWhileCreatingArticle + "\n" + s.Message
		return
	}

	*s = *getEmtpyCreateArticleStruct(acc)
	s.Message = generics.SuccessfulCreateArticle + "\n" + s.Message
}
