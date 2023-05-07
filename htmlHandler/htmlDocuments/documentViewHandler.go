package htmlDocuments

import (
	"API_MBundestag/dataLogic"
	"API_MBundestag/database"
	gen "API_MBundestag/help/generics"
	"API_MBundestag/htmlHandler"
	"API_MBundestag/htmlHandler/htmlBasics"
	"github.com/gin-gonic/gin"
	"html/template"
)

type PostViewStruct struct {
	database.Document
	CanAddTag bool
	BackgroundInfo
}

type DiscussionViewStruct struct {
	database.Document
	Commentable bool
	BackgroundInfo
}

type VoteViewStruct struct {
	database.Document
	Voteable bool
	BackgroundInfo
}

type BackgroundInfo struct {
	Admin           bool
	FormatString    string
	Message         string
	TagText         string
	TagColor        string
	Content         string
	SelectedAccount string
	Preview         template.HTML
	Accounts        database.AccountList
}

func GetDocumentViewPage(c *gin.Context) {
	doc := database.Document{}
	acc, legiable := dataLogic.CheckUserPrivileged(c, database.HeadAdmin, database.Admin, database.MediaAdmin, database.User)
	err := doc.GetByIDOnlyWithAccount(c.Query("uuid"), acc.ID)

	if err != nil {
		htmlBasics.MakeErrorPage(c, &acc, gen.DocumentDoesNotExists)
		return
	}

	makeDocumentToPage(doc, legiable, c, &acc, BackgroundInfo{
		Admin:        dataLogic.CheckIfHasRole(&acc, database.HeadAdmin, database.Admin),
		FormatString: gen.LongTimeString,
		TagColor:     "#FFFFFF",
		Message:      "",
	})
}

func makeDocumentToPage(doc database.Document, legiable bool, c *gin.Context, acc *database.Account, b BackgroundInfo) {
	switch doc.Type {
	case database.LegislativeText:
		canEdit := b.Admin || checkIfAdminInOrg(doc.Organisation, acc)
		htmlHandler.MakeSite(&PostViewStruct{Document: doc, CanAddTag: canEdit, BackgroundInfo: b}, c, acc)
	case database.RunningDiscussion:
		htmlHandler.FillOwnAccounts(&b, acc)
		htmlHandler.MakeSite(&DiscussionViewStruct{Document: doc, Commentable: legiable, BackgroundInfo: b}, c, acc)
	case database.FinishedDiscussion:
		htmlHandler.MakeSite(&DiscussionViewStruct{Document: doc, Commentable: false, BackgroundInfo: b}, c, acc)
	case database.RunningVote:
		htmlHandler.MakeSite(&VoteViewStruct{Document: doc, Voteable: legiable, BackgroundInfo: b}, c, acc)
	case database.FinishedVote:
		htmlHandler.MakeSite(&VoteViewStruct{Document: doc, Voteable: false, BackgroundInfo: b}, c, acc)
	}
}

func checkIfAdminInOrg(organisation string, acc *database.Account) bool {
	/*org := database.Organisation{}
	err := org.GetByNameAndOnlyWithAccount(organisation, acc.DisplayName)
	if err != nil {
		return false
	}
	for _, name := range org.Info.Admins {
		if name == acc.DisplayName {
			return true
		}
		admin := database.Account{}
		err = admin.GetByDisplayName(name)
		if err != nil {
			return false
		}
		if admin.Linked.Valid && admin.Linked.Int64 == acc.ID {
			return true
		}
	}*/
	return false
}

func PostDocumentViewPage(c *gin.Context) {
	doc := database.Document{}
	acc, _ := dataLogic.CheckUserPrivileged(c)
	err := doc.GetByIDOnlyWithAccount(c.Query("uuid"), acc.ID)

	if err != nil {
		htmlBasics.MakeErrorPage(c, &acc, gen.DocumentDoesNotExists)
		return
	}

	switch c.Query("type") {
	case "addTag":
		addTagToDocument(c, &acc, doc)
	case "hideTag":
		hideTagDocument(c, &acc, doc)
	case "comment":
		commentOnDiscussion(c, &acc, doc)
	case "hideComment":
		hideCommentFromDiscussion(c, &acc, doc)
	case "blockDocument":
		blockDocument(c, &acc, doc)
	default:
		htmlBasics.MakeErrorPage(c, &acc, gen.TypeDoesNotExist)
	}
}

var CanNotBlockDocument = "Du kannst keine Dokumente verstecken"
var ErrorWhileBlockingDocument = "Es ist ein Fehler beim blockieren aufgetreten"
var DocumentSuccessfulHidden = "Dokument wurde erfolgreich blockiert"
var DocumentSuccessfulDehidden = "Dokument wurde erfolgreich wieder freigeschaltet"

func blockDocument(c *gin.Context, acc *database.Account, doc database.Document) {
	if !(dataLogic.CheckIfHasRole(acc, database.HeadAdmin, database.Admin)) {
		htmlBasics.MakeErrorPage(c, acc, CanNotBlockDocument)
	}

	b := BackgroundInfo{
		Admin:        true,
		FormatString: gen.LongTimeString,
		TagColor:     "#FFFFFF",
		Message:      "",
	}

	hide := !doc.Blocked
	doc.Blocked = hide
	err := doc.SaveChanges()
	if err != nil {
		b.Message = ErrorWhileBlockingDocument
	}

	if hide {
		b.Message = DocumentSuccessfulHidden
	} else {
		b.Message = DocumentSuccessfulDehidden
	}

	makeDocumentToPage(doc, true, c, acc, b)
}
