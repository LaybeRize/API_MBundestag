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

type VoteCreateStruct struct {
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
	AmountVotes          int
	EmptyVote            CreateSingleVote
	Votes                []CreateSingleVote
}

type CreateSingleVote struct {
	Question               string
	Type                   database.VoteType
	Number                 int
	ShowNumbersWhileVoting bool
	ShowNamesWhileVoting   bool
	ShowNamesAfterVoting   bool
	Options                []string
}

func getEmptyCreateVoteStruct(acc *database.Account) *VoteCreateStruct {
	res := VoteCreateStruct{}
	htmlHandler.FillOwnAccounts(&res, acc)
	htmlHandler.FillOwnOrganisations(&res, acc)
	htmlHandler.FillAllNotSuspendedNames(&res)
	res.Info.Finishing = time.Now().Add(time.Hour*24 + 10*time.Minute)
	res.FormatForTime = generics.TimeParseDiscussion
	res.AmountVotes = 0
	res.EmptyVote = CreateSingleVote{
		Type:    database.SingleVote,
		Number:  10,
		Options: []string{},
	}
	res.Votes = []CreateSingleVote{}
	return &res
}

func GetVoteCreatePage(c *gin.Context) {
	acc, b := dataLogic.CheckUserPrivileged(c, database.HeadAdmin, database.Admin, database.MediaAdmin, database.User)
	if !b {
		htmlBasics.MakeErrorPage(c, &acc, generics.NotAuthorizedToView)
		return
	}

	res := getEmptyCreateVoteStruct(&acc)
	res.SelectedOrganisation = c.Query("org")
	res.SelectedAccount = c.Query("usr")
	htmlHandler.MakeSite(res, c, &acc)
}

func PostVoteCreatePage(c *gin.Context) {
	acc, b := dataLogic.CheckUserPrivileged(c, database.HeadAdmin, database.Admin, database.MediaAdmin, database.User)
	if !b {
		htmlBasics.MakeErrorPage(c, &acc, generics.NotAuthorizedToView)
		return
	}

	res := getEmptyCreateVoteStruct(&acc)
	res.fillVoteStructFromContext(c)
	res.makeVote(c, &acc)
}

func (voteStruct *VoteCreateStruct) fillVoteStructFromContext(c *gin.Context) {
	voteStruct.SelectedOrganisation = generics.GetText(c, "selectedOrganisation")
	voteStruct.SelectedAccount = generics.GetText(c, "selectedAccount")
	voteStruct.Content = generics.GetText(c, "content")
	voteStruct.Title = generics.GetText(c, "title")
	voteStruct.Subtitle = generics.GetText(c, "subtitle")
	voteStruct.MakePrivate = generics.GetBool(c, "private")
	voteStruct.FormatForTime = generics.TimeParseDiscussion
	voteStruct.Info = database.DocumentInfo{
		Poster:                    generics.GetStringArray(c, "poster"),
		Viewer:                    generics.GetStringArray(c, "allowed"),
		AnyPosterAllowed:          generics.GetBool(c, "anyPoster"),
		OrganisationPosterAllowed: generics.GetBool(c, "orgPoster"),
		Finishing:                 time.Now().Add(time.Hour*24 + 10*time.Minute),
	}
	voteStruct.Votes = *extractVotesFromContext(c)
	voteStruct.AmountVotes = len(voteStruct.Votes)
}

func extractVotesFromContext(c *gin.Context) *[]CreateSingleVote {
	array := []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"}
	var voteArray []CreateSingleVote
	voteArray = []CreateSingleVote{}
	for _, str := range array {
		v := CreateSingleVote{
			Question:               generics.GetText(c, "question"+str),
			Type:                   database.VoteType(generics.GetText(c, "selectVoteType"+str)),
			Number:                 generics.GetNumber(c, "maxValue"+str, 10, 2, 50),
			ShowNumbersWhileVoting: generics.GetBool(c, "showNumsW"+str),
			ShowNamesWhileVoting:   generics.GetBool(c, "showNamesW"+str),
			ShowNamesAfterVoting:   generics.GetBool(c, "showNamesA"+str),
			Options:                generics.GetStringArray(c, "option"+str),
		}
		if len(v.Options) != 0 || v.Question != "" || v.Number != 10 {
			voteArray = append(voteArray, v)
		}
	}
	return &voteArray
}

var maxTimeVoteInDays = 21
var SecretOrgsCanNotCreatePublicVote = "Eine geheime Organisation kann keine öffentlichen Abstimmungen veröffentlichen"
var OnPrivateVoteNotEveryoneCanVote = "An einer privaten Abstimmung kann nicht jeder teilnehmen"
var NonAdminsCanNotAddPersonToSecretVotes = "Nutzer können keine Abstimmung mit speziellen Personen in einer geheimen Organisation erstellen"
var PublicOrganisationCanNotPublishPrivateVotes = "Öffentliche Organisationen können keine privaten Abstimmungen erstellen"
var NonAdminsAreNotAllowedToLetEveryoneVote = "Nutzer dürfen nicht alle zu einer Abstimmung zulassen"
var CanNotCreateVoteWithoutVotes = "Du kannst keine Abstimmung erstellen in der nichts abzustimmen ist"
var CanNotCreateVoteWithoutVotees = "Du kannst keine Abstimmung erstellen, bei der keiner Abstimmen darf"
var OptionsMissingOnVote = "Kein Optionen angegeben bei Abstimmung %d"
var QuestionMissingOnVote = "Keine Frage bei Abstimmung %d angegeben"
var WrongVoteMethod = "In Abstimmung %d ist keine valide Abstimmungsmethode ausgewählt"

func (voteStruct *VoteCreateStruct) makeVote(c *gin.Context, acc *database.Account) {
	writer := &database.Account{}
	orga := &database.Organisation{}
	id := ""
	switch true {
	case voteStruct.parseFinishingTime(c):
	case voteStruct.checkIfTimeInWindow():
	case generics.CheckWriter(voteStruct, writer, acc):
	case generics.CheckOrgExists(voteStruct, orga):
	case voteStruct.organisationCheck(orga, writer):
	case voteStruct.checkIfBoolsAreCorrect(orga, writer):
	case generics.CheckTitelAndContentEmpty(voteStruct):
	case generics.CheckLengthContent(voteStruct, generics.PostContentLimit):
	case generics.CheckLengthTitle(voteStruct, generics.PostTitleLimit):
	case generics.CheckLengthSubtitle(voteStruct, generics.PostSubtitleLimit):
	case gen.CheckAccountList(voteStruct, &voteStruct.Info.Poster):
	case gen.CheckAccountList(voteStruct, &voteStruct.Info.Viewer):
	case voteStruct.addViewer(voteStruct.Info.Poster):
	case voteStruct.addViewer(voteStruct.Info.Viewer):
	case voteStruct.checkVotes():
	case voteStruct.createDocument(&id, writer.Flair):
	default:
		c.Redirect(http.StatusFound, "/document?uuid="+url.QueryEscape(id))
		return
	}

	htmlHandler.MakeSite(voteStruct, c, acc)
}

func (voteStruct *VoteCreateStruct) parseFinishingTime(c *gin.Context) bool {
	t, err := time.ParseInLocation(generics.TimeParseDiscussion, generics.GetText(c, "until"), time.Now().Location())
	if err == nil {
		voteStruct.Info.Finishing = t
	} else {
		voteStruct.Message = TimeNotValid + "\n" + voteStruct.Message
		return true
	}
	return false
}

func (voteStruct *VoteCreateStruct) checkIfTimeInWindow() bool {
	if !(voteStruct.Info.Finishing.After(time.Now().Add(time.Hour*24)) &&
		voteStruct.Info.Finishing.Before(time.Now().Add(time.Hour*24*time.Duration(maxTimeVoteInDays)))) {
		voteStruct.Message = fmt.Sprintf(TimeNotInAllowedInterval, maxTimeVoteInDays) + "\n" + voteStruct.Message
		return true
	}
	return false
}

func (voteStruct *VoteCreateStruct) organisationCheck(orga *database.Organisation, writer *database.Account) bool {
	isAdmin := helper.GetPositionOfString(orga.Info.Admins, writer.DisplayName) != -1 || writer.Role == database.HeadAdmin
	isUser := helper.GetPositionOfString(orga.Info.User, writer.DisplayName) != -1
	if !isAdmin && !isUser {
		voteStruct.Message = generics.YouAreNotAllowedForOrganisation + "\n" + voteStruct.Message
		return true
	}
	return false
}

func (voteStruct *VoteCreateStruct) checkIfBoolsAreCorrect(orga *database.Organisation, writer *database.Account) bool {
	isSecret := orga.Status == database.Secret
	isPublic := orga.Status == database.Public
	isAdmin := helper.GetPositionOfString(orga.Info.Admins, writer.DisplayName) != -1 || writer.Role == database.HeadAdmin
	//make sure that in secrete Organisation only private discussions are posted
	if isSecret && !voteStruct.MakePrivate {
		voteStruct.Message = SecretOrgsCanNotCreatePublicVote + "\n" + voteStruct.Message
		return true
	}
	//make sure that not all people are allowed to a private discussion (it kinda defeats the point)
	if voteStruct.MakePrivate && voteStruct.Info.AnyPosterAllowed {
		voteStruct.Message = OnPrivateVoteNotEveryoneCanVote + "\n" + voteStruct.Message
		return true
	}
	//make sure that non admins of secret organisations do not add people
	if isSecret && !isAdmin && (len(voteStruct.Info.Viewer) != 0 || len(voteStruct.Info.Poster) != 0) {
		voteStruct.Message = NonAdminsCanNotAddPersonToSecretVotes + "\n" + voteStruct.Message
		return true
	}
	//public organisations are not allowed to publicate private discussions
	if isPublic && voteStruct.MakePrivate {
		voteStruct.Message = PublicOrganisationCanNotPublishPrivateVotes + "\n" + voteStruct.Message
		return true
	}
	//if you are not admin of a public/private organisation you are not allowed to include everyone in a discussion
	if !isAdmin && voteStruct.Info.AnyPosterAllowed {
		voteStruct.Message = NonAdminsAreNotAllowedToLetEveryoneVote + "\n" + voteStruct.Message
		return true
	}
	return false
}

func (voteStruct *VoteCreateStruct) addViewer(array []string) bool {
	infoRef := &voteStruct.Info
	acc := database.Account{}
	for _, str := range array {
		err := acc.GetByDisplayName(str)
		if acc.Role == database.PressAccount {
			err = acc.GetByID(acc.Linked.Int64)
		}
		if err != nil {
			voteStruct.Message = generics.ViewerError + "\n" + voteStruct.Message
			return true
		}
		infoRef.Allowed = append(infoRef.Allowed, acc.DisplayName)
	}
	infoRef.Allowed = helper.RemoveDuplicates(infoRef.Allowed)
	return false
}

func (voteStruct *VoteCreateStruct) checkVotes() bool {
	if len(voteStruct.Info.Poster) == 0 && !voteStruct.Info.OrganisationPosterAllowed && !voteStruct.Info.AnyPosterAllowed {
		voteStruct.Message = CanNotCreateVoteWithoutVotees + "\n" + voteStruct.Message
		return true
	}
	if len(voteStruct.Votes) == 0 {
		voteStruct.Message = CanNotCreateVoteWithoutVotes + "\n" + voteStruct.Message
		return true
	}
	for pos, v := range voteStruct.Votes {
		if len(v.Options) == 0 {
			voteStruct.Message = fmt.Sprintf(OptionsMissingOnVote, pos+1) + "\n" + voteStruct.Message
			return true
		}
		if v.Question == "" {
			voteStruct.Message = fmt.Sprintf(QuestionMissingOnVote, pos+1) + "\n" + voteStruct.Message
			return true
		}
		if _, ok := database.VoteTranslation[v.Type]; !ok {
			voteStruct.Message = fmt.Sprintf(WrongVoteMethod, pos+1) + "\n" + voteStruct.Message
			return true
		}
	}
	return false
}

var ErrorWhileCreatingVotes = "Es ist ein Fehler beim erstellen der Abstimmungen aufgetreten. Bitte wende dich an einen Moderator"
var ErrorWhileCreatingVoteItself = "Es ist ein Fehler beim erstellen des Abstimmungsdokument aufgetreten. Bitte wende dich an eine Moderator"

func (voteStruct *VoteCreateStruct) createDocument(id *string, flair string) bool {
	*id = uuid.New().String()
	voteIds := []string{}
	for _, v := range voteStruct.Votes {
		newVote := database.Votes{
			UUID:                   uuid.New().String(),
			Parent:                 *id,
			Question:               v.Question,
			ShowNumbersWhileVoting: v.ShowNumbersWhileVoting,
			ShowNamesWhileVoting:   v.ShowNamesWhileVoting,
			ShowNamesAfterVoting:   v.ShowNamesAfterVoting,
			Finished:               false,
			Info: database.VoteInfo{
				Allowed: voteStruct.Info.Allowed,
				Results: map[string]database.Results{},
				Summary: database.Summary{
					Sums:         map[string]int{},
					RankedMap:    map[string]map[string]int{}, //first string is the person second one is the option
					Person:       map[string]string{},
					InvalidVotes: []string{},
					CSV:          "",
				},
				VoteMethod:  v.Type,
				MaxPosition: v.Number,
				Options:     v.Options,
			},
		}
		err := newVote.CreateMe()
		if err != nil {
			voteStruct.Message = ErrorWhileCreatingVotes + "\n" + voteStruct.Message
			return true
		}
		voteIds = append(voteIds, newVote.UUID)
	}

	voteStruct.Info.Votes = voteIds
	voteStruct.Info.Finishing = voteStruct.Info.Finishing.UTC()

	doc := database.Document{
		UUID:         *id,
		Organisation: voteStruct.SelectedOrganisation,
		Type:         database.UnfinishedVote,
		Author:       voteStruct.SelectedAccount,
		Flair:        flair,
		Title:        voteStruct.Title,
		Subtitle:     sql.NullString{Valid: voteStruct.Subtitle != "", String: voteStruct.Subtitle},
		HTMLContent:  helper.CreateHTML(voteStruct.Content),
		Private:      voteStruct.MakePrivate,
		Info:         voteStruct.Info,
	}

	err := doc.CreateMe()
	if err != nil {
		voteStruct.Info.Finishing = voteStruct.Info.Finishing.In(time.Now().Location())
		voteStruct.Message = ErrorWhileCreatingVoteItself + "\n" + voteStruct.Message
		return true
	}
	return false
}
