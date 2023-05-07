package htmlDocuments

import (
	"API_MBundestag/dataLogic"
	"API_MBundestag/database"
	"API_MBundestag/help/generics"
	"API_MBundestag/htmlHandler"
	"API_MBundestag/htmlHandler/htmlBasics"
	"github.com/gin-gonic/gin"
	"net/url"
)

type ViewDocumentListStruct struct {
	HasNext          bool
	HasBefore        bool
	NextUUID         string
	BeforeUUID       string
	Amount           int
	ExtraQueryString string
	DocumentList     database.DocumentList
	FormatString     string
}

var ErrorWhileLoadingDocuments = "Es ist ein Fehler beim laden der Dokumente aufgetreten"

func GetDocumentListView(c *gin.Context) {
	acc, admin := dataLogic.CheckUserPrivileged(c, database.HeadAdmin, database.Admin)

	i := htmlHandler.ExtractAmount(c, 1, 50, 20)
	documentStruct := &ViewDocumentListStruct{
		FormatString: generics.LongTimeString,
		Amount:       i,
	}
	var err error
	if generics.GetIfType(c, "before") {
		err = documentStruct.validateDocumentReadPageBefore(c, i, &acc, admin)

	} else {
		err = documentStruct.validateDocumentReadPageNext(c, i, &acc, admin)
	}
	if err != nil {
		htmlBasics.MakeErrorPage(c, &acc, ErrorWhileLoadingDocuments)
		return
	}
	htmlHandler.MakeSite(documentStruct, c, &acc)
}

func (listStruct *ViewDocumentListStruct) validateDocumentReadPageNext(c *gin.Context, i int, acc *database.Account, admin bool) error {
	m := getMapForDocument(c, acc, admin)
	err, exists := listStruct.DocumentList.GetDocumentsAfter(c.Query("uuid"), i+1, acc.ID, m, getTypesForDocument(c)...)
	if len(listStruct.DocumentList) == 0 {
		return err
	}
	if len(listStruct.DocumentList) == i+1 {
		listStruct.HasNext = true
		listStruct.NextUUID = listStruct.DocumentList[i-1].UUID
		listStruct.DocumentList = listStruct.DocumentList[:i]
	}
	if exists {
		listStruct.HasBefore = true
		listStruct.BeforeUUID = listStruct.DocumentList[0].UUID
	}
	listStruct.ExtraQueryString = getExtraQueryString(m)
	return err
}

func (listStruct *ViewDocumentListStruct) validateDocumentReadPageBefore(c *gin.Context, i int, acc *database.Account, admin bool) error {
	m := getMapForDocument(c, acc, admin)
	err, exists := listStruct.DocumentList.GetDocumentsBefore(c.Query("uuid"), i+1, acc.ID, m, getTypesForDocument(c)...)
	if len(listStruct.DocumentList) == 0 {
		return err
	}
	if len(listStruct.DocumentList) == i+1 {
		listStruct.HasBefore = true
		listStruct.DocumentList = listStruct.DocumentList[1:]
		listStruct.BeforeUUID = listStruct.DocumentList[0].UUID
	}
	if exists {
		listStruct.HasNext = true
		listStruct.NextUUID = listStruct.DocumentList[len(listStruct.DocumentList)-1].UUID
	}
	listStruct.ExtraQueryString = getExtraQueryString(m)
	return err
}
func getTypesForDocument(c *gin.Context) (arr []database.DocumentType) {
	arr = []database.DocumentType{}
	switch c.Query("types") {
	case "post":
		arr = []database.DocumentType{database.LegislativeText}
	case "discussion":
		arr = []database.DocumentType{database.RunningDiscussion, database.FinishedDiscussion}
	case "vote":
		arr = []database.DocumentType{database.FinishedVote, database.RunningVote}
	case "post,discussion":
		arr = []database.DocumentType{database.RunningDiscussion, database.FinishedDiscussion,
			database.LegislativeText}
	case "post,vote":
		arr = []database.DocumentType{database.FinishedVote, database.RunningVote,
			database.LegislativeText}
	case "discussion,vote":
		arr = []database.DocumentType{database.RunningDiscussion, database.FinishedDiscussion,
			database.FinishedVote, database.RunningVote}
	default:
	}
	return
}

func getMapForDocument(c *gin.Context, acc *database.Account, admin bool) map[string]string {
	adminStr := "false"
	if admin {
		adminStr = "true"
	}
	blocked := c.Query("blocked")
	if !admin {
		blocked = "false"
	}
	return map[string]string{
		"organisation": c.Query("organisation"),
		"author":       c.Query("author"),
		"title":        c.Query("title"),
		"displayname":  acc.DisplayName,
		"admin":        adminStr,
		"blocked":      blocked,
	}
}

func getExtraQueryString(m map[string]string) (result string) {
	result = ""
	if m["organisation"] != "" {
		result += "&organisation=" + url.QueryEscape(m["organisation"])
	}
	if m["author"] != "" {
		result += "&author=" + url.QueryEscape(m["author"])
	}
	if m["title"] != "" {
		result += "&title=" + url.QueryEscape(m["title"])
	}
	if m["blocked"] == "true" {
		result += "&blocked=true"
	}
	return
}
