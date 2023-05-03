package htmlDocuments

import (
	"API_MBundestag/dataLogic"
	"API_MBundestag/database_old"
	"API_MBundestag/help/generics"
	"API_MBundestag/htmlHandler"
	gen "API_MBundestag/htmlHandler/generics"
	"API_MBundestag/htmlHandler/htmlBasics"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"net/url"
	"time"
)

type DiscussionCreateStruct struct {
	Info                 database.DocumentInfo
	Names                []string
	SelectedAccount      string
	Accounts             database.AccountList
	SelectedOrganisation string
	Organisations        database.OrganisationList
	Content              string
	Title                string
	Subtitle             string
	MakePrivate          bool
	FormatForTime        string
	Message              string
}

func getEmptyCreateDiscussionStruct(acc *database.Account) *DiscussionCreateStruct {
	res := DiscussionCreateStruct{}
	htmlHandler.FillOwnAccounts(&res, acc)
	htmlHandler.FillOwnOrganisations(&res, acc)
	htmlHandler.FillAllNotSuspendedNames(&res)
	res.Info.Finishing = time.Now().Add(time.Hour*24 + 10*time.Minute)
	res.FormatForTime = generics.TimeParseDiscussion
	return &res
}

func GetDiscussionCreatePage(c *gin.Context) {
	acc, b := dataLogic.CheckUserPrivileged(c, database.HeadAdmin, database.Admin, database.MediaAdmin, database.User)
	if !b {
		htmlBasics.MakeErrorPage(c, &acc, generics.NotAuthorizedToView)
		return
	}

	res := getEmptyCreateDiscussionStruct(&acc)
	res.SelectedOrganisation = c.Query("org")
	res.SelectedAccount = c.Query("usr")
	htmlHandler.MakeSite(res, c, &acc)
}

func PostDiscussionCreatePage(c *gin.Context) {
	acc, b := dataLogic.CheckUserPrivileged(c, database.HeadAdmin, database.Admin, database.MediaAdmin, database.User)
	if !b {
		htmlBasics.MakeErrorPage(c, &acc, generics.NotAuthorizedToView)
		return
	}

	res := getEmptyCreateDiscussionStruct(&acc)
	res.fillStructFromContext(c)

	res.makeDiscussion(c, &acc)
}

func (discuss *DiscussionCreateStruct) fillStructFromContext(c *gin.Context) {
	discuss.SelectedOrganisation = htmlHandler.GetText(c, "selectedOrganisation")
	discuss.SelectedAccount = htmlHandler.GetText(c, "selectedAccount")
	discuss.Content = htmlHandler.GetText(c, "content")
	discuss.Title = htmlHandler.GetText(c, "title")
	discuss.Subtitle = htmlHandler.GetText(c, "subtitle")
	discuss.MakePrivate = htmlHandler.GetBool(c, "private")
	discuss.FormatForTime = generics.TimeParseDiscussion
	discuss.Info = database.DocumentInfo{
		Poster:                    htmlHandler.GetStringArray(c, "poster"),
		Viewer:                    htmlHandler.GetStringArray(c, "allowed"),
		AnyPosterAllowed:          htmlHandler.GetBool(c, "anyPoster"),
		OrganisationPosterAllowed: htmlHandler.GetBool(c, "orgPoster"),
		Finishing:                 time.Now().Add(time.Hour*24 + 10*time.Minute),
	}
}

var maxTimeDiscussionInDays = 14
var TimeNotValid = "Der angegebene Zeitstempel ist nicht valide"
var TimeNotInAllowedInterval = "Das Ende ist entweder weniger als einen Tag oder mehr als %d Tage entfernt"
var OnPrivateDiscussionNotEveryoneCanComment = "Bei einer privaten Diskussion dürfen nicht alle Personen kommentieren (bitte mal nachdenken)"
var SecretOrgsCanNotCreatePublicDiscussion = "Eine geheime Organisation darf keine öffentliche Diskussion erstellen"
var NonAdminsCanNotAddPersonToSecretDiscussions = "Ein normales Mitglied darf nur alle Organsationsmitglieder oder niemanden zu einer Diskussion in einer geheimen Organisation zulassen"
var PublicOrganisationCanNotPublishPrivateDiscussion = "Eine öffentliche Organisation darf keine private Diskussion veröffentlichen"
var NonAdminsAreNotAllowedToLetEveryoneComment = "Normale Nutzer können nicht allen Personen erlauben zu kommentieren"
var DiscussionCreationFailed = "Es ist ein Fehler beim erstellen der Diskussion aufgetreten"

func (discuss *DiscussionCreateStruct) makeDiscussion(c *gin.Context, acc *database.Account) {
	writer := &database.Account{}
	orga := &database.Organisation{}
	id := ""
	switch true {
	case discuss.parseFinishingTime(c):
	case discuss.checkIfTimeInWindow():
	case htmlHandler.CheckWriter(discuss, writer, acc):
	case htmlHandler.CheckOrgExists(discuss, orga):
	case discuss.organisationCheck(orga, writer):
	case discuss.checkIfBoolsAreCorrect(orga, writer):
	case htmlHandler.CheckTitelAndContentEmpty(discuss):
	case htmlHandler.CheckLengthContent(discuss, generics.PostContentLimit):
	case htmlHandler.CheckLengthTitle(discuss, generics.PostTitleLimit):
	case htmlHandler.CheckLengthSubtitle(discuss, generics.PostSubtitleLimit):
	case gen.CheckAccountList(discuss, &discuss.Info.Poster):
	case gen.CheckAccountList(discuss, &discuss.Info.Viewer):
	case discuss.addViewer(discuss.Info.Poster):
	case discuss.addViewer(discuss.Info.Viewer):
	case discuss.createDocument(&id, writer.Flair):
	default:
		c.Redirect(http.StatusFound, "/document?uuid="+url.QueryEscape(id))
		return
	}

	htmlHandler.MakeSite(discuss, c, acc)
}

func (discuss *DiscussionCreateStruct) organisationCheck(orga *database.Organisation, writer *database.Account) bool {
	isAdmin := helper.GetPositionOfString(orga.Info.Admins, writer.DisplayName) != -1 || writer.Role == database.HeadAdmin
	isUser := helper.GetPositionOfString(orga.Info.User, writer.DisplayName) != -1
	if !isAdmin && !isUser {
		discuss.Message = generics.YouAreNotAllowedForOrganisation + "\n" + discuss.Message
		return true
	}
	return false
}

func (discuss *DiscussionCreateStruct) parseFinishingTime(c *gin.Context) bool {
	t, err := time.ParseInLocation(generics.TimeParseDiscussion, htmlHandler.GetText(c, "until"), time.Now().Location())
	if err == nil {
		discuss.Info.Finishing = t
	} else {
		discuss.Message = TimeNotValid + "\n" + discuss.Message
		return true
	}
	return false
}

func (discuss *DiscussionCreateStruct) checkIfTimeInWindow() bool {
	if !(discuss.Info.Finishing.After(time.Now().Add(time.Hour*24)) &&
		discuss.Info.Finishing.Before(time.Now().Add(time.Hour*24*time.Duration(maxTimeDiscussionInDays)))) {
		discuss.Message = fmt.Sprintf(TimeNotInAllowedInterval, maxTimeDiscussionInDays) + "\n" + discuss.Message
		return true
	}
	return false
}

func (discuss *DiscussionCreateStruct) checkIfBoolsAreCorrect(orga *database.Organisation, writer *database.Account) bool {
	isSecret := orga.Status == database.Secret
	isPublic := orga.Status == database.Public
	isAdmin := helper.GetPositionOfString(orga.Info.Admins, writer.DisplayName) != -1 || writer.Role == database.HeadAdmin
	//make sure that in secrete Organisation only private discussions are posted
	if isSecret && !discuss.MakePrivate {
		discuss.Message = SecretOrgsCanNotCreatePublicDiscussion + "\n" + discuss.Message
		return true
	}
	//make sure that not all people are allowed to a private discussion (it kinda defeats the point)
	if discuss.MakePrivate && discuss.Info.AnyPosterAllowed {
		discuss.Message = OnPrivateDiscussionNotEveryoneCanComment + "\n" + discuss.Message
		return true
	}
	//make sure that non admins of secret organisations do not add people
	if isSecret && !isAdmin && (len(discuss.Info.Viewer) != 0 || len(discuss.Info.Poster) != 0) {
		discuss.Message = NonAdminsCanNotAddPersonToSecretDiscussions + "\n" + discuss.Message
		return true
	}
	//public organisations are not allowed to publicate private discussions
	if isPublic && discuss.MakePrivate {
		discuss.Message = PublicOrganisationCanNotPublishPrivateDiscussion + "\n" + discuss.Message
		return true
	}
	//if you are not admin of a public/private organisation you are not allowed to include everyone in a discussion
	if !isAdmin && discuss.Info.AnyPosterAllowed {
		discuss.Message = NonAdminsAreNotAllowedToLetEveryoneComment + "\n" + discuss.Message
		return true
	}
	return false
}

func (discuss *DiscussionCreateStruct) addViewer(array []string) bool {
	infoRef := &discuss.Info
	acc := database.Account{}
	for _, str := range array {
		err := acc.GetByDisplayName(str)
		if acc.Role == database.PressAccount {
			err = acc.GetByID(acc.Linked.Int64)
		}
		if err != nil {
			discuss.Message = generics.ViewerError + "\n" + discuss.Message
			return true
		}
		infoRef.Allowed = append(infoRef.Allowed, acc.DisplayName)
	}
	infoRef.Allowed = helper.RemoveDuplicates(infoRef.Allowed)
	return false
}

func (discuss *DiscussionCreateStruct) createDocument(id *string, flair string) bool {

	discuss.Info.Finishing = discuss.Info.Finishing.UTC()

	doc := database.Document{
		UUID:         uuid.New().String(),
		Organisation: discuss.SelectedOrganisation,
		Type:         database.Discussion,
		Author:       discuss.SelectedAccount,
		Flair:        flair,
		Title:        discuss.Title,
		Subtitle:     sql.NullString{Valid: discuss.Subtitle != "", String: discuss.Subtitle},
		HTMLContent:  helper.CreateHTML(discuss.Content),
		Private:      discuss.MakePrivate,
		Info:         discuss.Info,
	}

	err := doc.CreateMe()
	if err != nil {
		discuss.Info.Finishing = discuss.Info.Finishing.In(time.Now().Location())
		discuss.Message = DiscussionCreationFailed + "\n" + discuss.Message
		return true
	}
	*id = doc.UUID
	return false
}
