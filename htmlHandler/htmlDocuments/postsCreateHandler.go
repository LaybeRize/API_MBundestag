package htmlDocuments

import (
	"API_MBundestag/dataLogic"
	"API_MBundestag/database"
	"API_MBundestag/help"
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
	post.SelectedOrganisation = generics.GetText(c, "selectedOrganisation")
	post.SelectedAccount = generics.GetText(c, "selectedAccount")
	post.Content = generics.GetText(c, "content")
	post.Tag = generics.GetText(c, "tag")
	post.Title = generics.GetText(c, "title")
	post.Subtitle = generics.GetText(c, "subtitle")
}

func (post *CreatePostPageStruct) validateCreatePost(c *gin.Context, acc *database.Account) {
	writer := &database.Account{}
	orga := &database.Organisation{}
	id := ""
	switch true {
	case generics.CheckTitelAndContentEmpty(post):
	case generics.CheckLengthContent(post, generics.PostContentLimit):
	case generics.CheckLengthTitle(post, generics.PostTitleLimit):
	case generics.CheckLengthSubtitle(post, generics.PostSubtitleLimit):
	case generics.CheckLengthField(post, generics.PostTagLimit, "Tag", generics.TagTooLong):
	case generics.CheckWriter(post, writer, acc):
	case generics.CheckOrgExists(post, orga):
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
	/*if help.GetPositionOfString(orga.Info.Admins, writer.DisplayName) == -1 && writer.Role != database.HeadAdmin {
		post.Message = generics.YouAreNotAllowedForOrganisation + "\n" + post.Message
		return true
	}*/
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
		HTMLContent:  help.CreateHTML(post.Content),
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
