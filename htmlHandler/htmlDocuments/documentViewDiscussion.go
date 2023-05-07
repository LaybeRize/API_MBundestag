package htmlDocuments

import (
	"API_MBundestag/dataLogic"
	"API_MBundestag/database"
	"API_MBundestag/help/generics"
	"API_MBundestag/htmlHandler"
	"API_MBundestag/htmlHandler/htmlBasics"
	"github.com/gin-gonic/gin"
	"time"
)

var CanNotHideComments = "Dir fehlt die Berechtigung Kommentare zu verstecken"
var CommentUUIDDoesNotExists = "Der Kommentar mit der UUID existiert nicht"
var CommentChangeNotSuccessfulSaved = "Kommentare Status konnte nicht geändert werden"
var CommentSuccessfulHidden = "Kommentar erfolgreich versteckt"
var CommentSuccessfulDehidden = "Kommentare wurde erfolgreich wiederhergestellt"
var DocumentIsNotADiscussion = "Dokument ist keine Diskussion"

func hideCommentFromDiscussion(c *gin.Context, acc *database.Account, doc database.Document) {
	if !(dataLogic.CheckIfHasRole(acc, database.HeadAdmin, database.Admin)) {
		htmlBasics.MakeErrorPage(c, acc, CanNotHideComments)
	}

	b := BackgroundInfo{
		Admin:        true,
		FormatString: generics.LongTimeString,
	}

	legiable := true
	hide := false
	switch true {
	case b.checkIfDiscussion(&doc):
	case b.checkIfCommentExists(&doc, c.Query("comment"), &hide, &legiable):
	case b.trySaving(&doc):
	case hide:
		b.Message = CommentSuccessfulHidden + "\n" + b.Message
	default:
		b.Message = CommentSuccessfulDehidden + "\n" + b.Message
	}

	makeDocumentToPage(doc, legiable, c, acc, b)
}

func (b *BackgroundInfo) checkIfDiscussion(doc *database.Document) bool {
	if doc.Type != database.RunningDiscussion && doc.Type != database.FinishedDiscussion {
		b.Message = DocumentIsNotADiscussion + "\n" + b.Message
		return true
	}
	return false
}

func (b *BackgroundInfo) checkIfCommentExists(doc *database.Document, comment string, hide *bool, legible *bool) bool {
	exists := false
	*legible = doc.Type == database.RunningDiscussion
	for i, tag := range doc.Info.Discussion {
		if tag.UUID == comment {
			*hide = !doc.Info.Discussion[i].Hidden
			doc.Info.Discussion[i].Hidden = *hide
			exists = true
		}
	}

	if !exists {
		b.Message = CommentUUIDDoesNotExists + "\n" + b.Message
		return true
	}
	return false
}

func (b *BackgroundInfo) trySaving(doc *database.Document) bool {
	err := doc.SaveChanges()
	if err != nil {
		b.Message = CommentChangeNotSuccessfulSaved + "\n" + b.Message
		return true
	}
	return false
}

var CanNotCommentOnDiscussion = "Du kannst ohne Account nicht kommentieren"
var CanNotCommentOnNonDiscussions = "Du kannst keinen Kommentar unter ein einem Dokument das keine Diskussion ist, oder dessen Diskussionszeitraum abgelaufen ist, machen"
var DiscussionCloseError = "Diskussion konnte nicht korrekt beendet werden"
var DiscussionCommentLength = 4000
var YouCanNotCommentWithThisAccount = "Mit dem Account darf nicht kommentiert werden"
var ProblemWithOrganisationAccess = "Es ist ein Fehler beim Zugriff auf die Datenbank aufgetreten"
var CommentCanNotBeEmpty = "Kommentare dürfen nicht leer sein"
var ErrorWhileCreatingComment = "Es ist ein Fehler beim erstellen des Kommentars aufgetreten"
var CommentCreated = "Kommentar wurde erfolgreich erstellt"

func commentOnDiscussion(c *gin.Context, acc *database.Account, doc database.Document) {
	if !(dataLogic.CheckIfHasRole(acc, database.HeadAdmin, database.Admin, database.MediaAdmin, database.User)) {
		htmlBasics.MakeErrorPage(c, acc, CanNotCommentOnDiscussion)
	}

	b := BackgroundInfo{
		Admin:           dataLogic.CheckIfHasRole(acc, database.HeadAdmin, database.Admin),
		FormatString:    generics.LongTimeString,
		Content:         generics.GetText(c, "content"),
		SelectedAccount: generics.GetText(c, "selectedAccount"),
	}
	htmlHandler.FillOwnAccounts(&b, acc)

	legiable := false
	writer := &database.Account{}
	switch true {
	case b.checkIfRunningDiscussion(&doc, &legiable):
	case b.checkIfDiscussionHasRunOut(&doc, &legiable):
	case generics.CheckFieldNotEmpty(&b, "Content", CommentCanNotBeEmpty):
	case generics.CheckLengthContent(&b, DiscussionCommentLength):
	case generics.CheckWriter(&b, writer, acc):
	case b.checkIfAllowedComment(&doc, writer):
	case b.trySavingComment(&doc, writer):
	default:
		b.Message = CommentCreated + "\n" + b.Message
		b.Content = ""
	}

	makeDocumentToPage(doc, legiable, c, acc, b)
}

func (b *BackgroundInfo) checkIfRunningDiscussion(doc *database.Document, legiable *bool) bool {
	if doc.Type != database.RunningDiscussion {
		b.Message = CanNotCommentOnNonDiscussions + "\n" + b.Message
		return true
	}
	*legiable = true
	return false
}

func (b *BackgroundInfo) checkIfDiscussionHasRunOut(doc *database.Document, legibale *bool) bool {
	if doc.Info.Finishing.Before(time.Now().UTC()) {
		err := dataLogic.CloseDiscussionOrVote(doc.UUID)
		if err != nil {
			b.Message = DiscussionCloseError + "\n" + b.Message
		}
		*legibale = err != nil
		return true
	}
	return false
}

func (b *BackgroundInfo) checkIfAllowedComment(doc *database.Document, writer *database.Account) bool {
	/*if doc.Info.AnyPosterAllowed {
		return false
	}
	//check if user is added as poster if neither org-members or anyone is allowed to comment
	if !doc.Info.OrganisationPosterAllowed {
		if help.GetPositionOfString(doc.Info.Poster, writer.DisplayName) == -1 {
			b.Message = YouCanNotCommentWithThisAccount + "\n" + b.Message
			return true
		}
		return false
	}
	//if org members are allowed check if writer is member or part of the poster group
	org := database.Organisation{}
	err := org.GetByName(doc.Organisation)
	if err != nil {
		b.Message = ProblemWithOrganisationAccess + "\n" + b.Message
		return true
	}
	if helper.GetPositionOfString(doc.Info.Poster, writer.DisplayName) == -1 &&
		helper.GetPositionOfString(org.Info.User, writer.DisplayName) == -1 &&
		helper.GetPositionOfString(org.Info.Admins, writer.DisplayName) == -1 {
		b.Message = YouCanNotCommentWithThisAccount + "\n" + b.Message
		return true
	}*/
	return false
}

func (b *BackgroundInfo) trySavingComment(doc *database.Document, writer *database.Account) bool {
	/*doc.Info.Discussion = append(doc.Info.Discussion, database.Discussions{
		UUID:        uuid.New().String(),
		Hidden:      false,
		Written:     time.Now().UTC(),
		Author:      writer.DisplayName,
		Flair:       writer.Flair,
		HTMLContent: helper.CreateHTML(b.Content),
	})
	err := doc.SaveChanges()
	if err != nil {
		doc.Info.Discussion = doc.Info.Discussion[:len(doc.Info.Discussion)-1]
		b.Message = ErrorWhileCreatingComment + "\n" + b.Message
		return true
	}*/
	return false
}
