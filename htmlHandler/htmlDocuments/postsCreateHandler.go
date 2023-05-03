package htmlDocuments

import (
	"API_MBundestag/dataLogic"
	"API_MBundestag/database_old"
	"API_MBundestag/help/generics"
	"API_MBundestag/htmlHandler"
	"API_MBundestag/htmlHandler/htmlBasics"
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"net/url"
	"time"
)

type CreatePostPageStruct struct {
	Message              string
	SelectedAccount      string
	Accounts             database.AccountList
	SelectedOrganisation string
	Organisations        database.OrganisationList
	Content              string
	Title                string
	Subtitle             string
	Tag                  string
}

func getEmptyCreatePostStruct(acc *database.Account) *CreatePostPageStruct {
	res := CreatePostPageStruct{}
	htmlHandler.FillOwnAccounts(&res, acc)
	htmlHandler.FillOwnOrganisations(&res, acc)
	return &res
}

func GetPostsCreateHandler(c *gin.Context) {
	acc, b := dataLogic.CheckUserPrivileged(c, database.HeadAdmin, database.Admin, database.MediaAdmin, database.User)
	if !b {
		htmlBasics.MakeErrorPage(c, &acc, generics.NotAuthorizedToView)
		return
	}

	res := getEmptyCreatePostStruct(&acc)
	res.SelectedOrganisation = c.Query("org")
	res.SelectedAccount = c.Query("usr")
	htmlHandler.MakeSite(res, c, &acc)
}

func PostPostsCreateHandler(c *gin.Context) {
	acc, b := dataLogic.CheckUserPrivileged(c, database.HeadAdmin, database.Admin, database.MediaAdmin, database.User)
	if !b {
		htmlBasics.MakeErrorPage(c, &acc, generics.NotAuthorizedToView)
		return
	}

	res := getEmptyCreatePostStruct(&acc)
	res.fillStructFromContext(c)
	res.validateCreatePost(c, &acc)
}

func (post *CreatePostPageStruct) fillStructFromContext(c *gin.Context) {
	post.SelectedOrganisation = htmlHandler.GetText(c, "selectedOrganisation")
	post.SelectedAccount = htmlHandler.GetText(c, "selectedAccount")
	post.Content = htmlHandler.GetText(c, "content")
	post.Tag = htmlHandler.GetText(c, "tag")
	post.Title = htmlHandler.GetText(c, "title")
	post.Subtitle = htmlHandler.GetText(c, "subtitle")
}

func (post *CreatePostPageStruct) validateCreatePost(c *gin.Context, acc *database.Account) {
	writer := &database.Account{}
	orga := &database.Organisation{}
	id := ""
	switch true {
	case htmlHandler.CheckTitelAndContentEmpty(post):
	case htmlHandler.CheckLengthContent(post, generics.PostContentLimit):
	case htmlHandler.CheckLengthTitle(post, generics.PostTitleLimit):
	case htmlHandler.CheckLengthSubtitle(post, generics.PostSubtitleLimit):
	case htmlHandler.CheckLengthField(post, generics.PostTagLimit, "Tag", generics.TagTooLong):
	case htmlHandler.CheckWriter(post, writer, acc):
	case htmlHandler.CheckOrgExists(post, orga):
	case post.checkIfAllowedToPost(orga, writer):
	case post.checkOrgaNotSecret(orga):
	case post.tryCreationPost(&id, writer):
	default:
		c.Redirect(http.StatusFound, "/document?uuid="+url.QueryEscape(id))
		return
	}

	htmlHandler.MakeSite(post, c, acc)
}

func (post *CreatePostPageStruct) checkIfAllowedToPost(orga *database.Organisation, writer *database.Account) bool {
	if helper.GetPositionOfString(orga.Info.Admins, writer.DisplayName) == -1 && writer.Role != database.HeadAdmin {
		post.Message = generics.YouAreNotAllowedForOrganisation + "\n" + post.Message
		return true
	}
	return false
}

func (post *CreatePostPageStruct) checkOrgaNotSecret(orga *database.Organisation) bool {
	if orga.Status == database.Secret {
		post.Message = generics.SecretOrgsCanNotCreatePosts + "\n" + post.Message
		return true
	}
	return false
}

func (post *CreatePostPageStruct) tryCreationPost(id *string, writer *database.Account) bool {
	doc := database.Document{
		UUID:         uuid.New().String(),
		Organisation: post.SelectedOrganisation,
		Type:         database.LegislativeText,
		Author:       writer.DisplayName,
		Flair:        writer.Flair,
		Title:        post.Title,
		Subtitle:     sql.NullString{Valid: post.Subtitle != "", String: post.Subtitle},
		HTMLContent:  helper.CreateHTML(post.Content),
		Private:      false,
		Info: database.DocumentInfo{
			Post: []database.Posts{},
		},
	}

	if post.Tag != "" {
		doc.Info.Post = append(doc.Info.Post, database.Posts{
			UUID:      uuid.New().String(),
			Hidden:    false,
			Submitted: time.Now().UTC(),
			Info:      post.Tag,
			Color:     "#FFFFFF",
		})
	}

	err := doc.CreateMe()
	if err != nil {
		post.Message = generics.PostCreationFailed + "\n" + post.Message
		return true
	}
	*id = doc.UUID
	return false
}
