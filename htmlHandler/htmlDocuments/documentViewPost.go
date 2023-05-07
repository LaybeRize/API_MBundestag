package htmlDocuments

import (
	"API_MBundestag/dataLogic"
	"API_MBundestag/database"
	"API_MBundestag/help/generics"
	"API_MBundestag/htmlHandler/htmlBasics"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"regexp"
	"time"
)

var CanNotHideTag = "Du bist nicht berechtigt Tags zu verstecken/freizustellen"
var CanNotHideTagOnNonLegislativeText = "Eine Diskussion oder Abstimmung besitzt keine Tags"
var TagUUIDDoesNotExists = "Die Tag UUID existiert nicht"
var TagSuccessfulHidden = "Tag wurde erfolgreich versteckt"
var TagSuccessfulDehidden = "Tag wurde erfolgreich wiederhergestellt"

func hideTagDocument(c *gin.Context, acc *database.Account, doc database.Document) {
	if !(dataLogic.CheckIfHasRole(acc, database.HeadAdmin, database.Admin)) {
		htmlBasics.MakeErrorPage(c, acc, CanNotHideTag)
	}
	b := BackgroundInfo{
		Admin:        dataLogic.CheckIfHasRole(acc, database.HeadAdmin, database.Admin),
		FormatString: generics.LongTimeString,
		TagColor:     "#FFFFFF",
	}

	hide := false
	switch true {
	case b.checkIfLegislativeText(&doc, CanNotHideTagOnNonLegislativeText):
	case b.checkIfTagExists(&doc, c, &hide):
	case b.tryChanges(&doc):
	case hide:
		b.Message = TagSuccessfulHidden
	default:
		b.Message = TagSuccessfulDehidden
	}
	makeDocumentToPage(doc, true, c, acc, b)
}

func (b *BackgroundInfo) checkIfLegislativeText(doc *database.Document, message string) bool {
	if doc.Type != database.LegislativeText {
		b.Message = message
		return true
	}
	return false
}

func (b *BackgroundInfo) checkIfTagExists(doc *database.Document, c *gin.Context, hide *bool) bool {
	exists := false
	for i, tag := range doc.Info.Post {
		if tag.UUID == c.Query("tag") {
			*hide = !doc.Info.Post[i].Hidden
			doc.Info.Post[i].Hidden = *hide
			exists = true
		}
	}

	if !exists {
		b.Message = TagUUIDDoesNotExists
		return true
	}
	return false
}

func (b *BackgroundInfo) tryChanges(doc *database.Document) bool {
	err := doc.SaveChanges()
	if err != nil {
		b.Message = TagNotSuccessfulSaved
		return true
	}
	return false
}

var CanNotAddTag = "Du hast keine Berechtigung diesem Dokument ein Tag anzufügen"
var CanNotAddTagToNonLegislativeText = "Einer Diskussion oder Abstimmung kann kein Tag hinzugefügt werden"
var TagEmpty = "Tag darf nicht leer sein"
var TagNotSuccessfulSaved = "Tag konnte nicht korrekt gespeichert werden"
var SuccessfulAddedTag = "Tag wurde erfolgreich angehängt"

func addTagToDocument(c *gin.Context, acc *database.Account, doc database.Document) {
	if !(checkIfAdminInOrg(doc.Organisation, acc) || dataLogic.CheckIfHasRole(acc, database.HeadAdmin, database.Admin)) {
		htmlBasics.MakeErrorPage(c, acc, CanNotAddTag)
	}

	b := BackgroundInfo{
		Admin:        dataLogic.CheckIfHasRole(acc, database.HeadAdmin, database.Admin),
		FormatString: generics.LongTimeString,
		TagText:      generics.GetText(c, "tag"),
		TagColor:     generics.GetText(c, "color"),
	}
	if m, err := regexp.MatchString(`^#[0-9a-fA-F]{6}$`, b.TagColor); err != nil || !m {
		b.TagColor = "#FFFFFF"
	}

	switch true {
	case b.checkIfLegislativeText(&doc, CanNotAddTagToNonLegislativeText):
	case generics.CheckFieldNotEmpty(&b, "TagText", TagEmpty):
	case generics.CheckLengthField(&b, generics.PostTagLimit, "TagText", generics.TagTooLong):
	case b.tryAddingNewTag(&doc):
	default:
		b.TagText = ""
		b.Message = SuccessfulAddedTag
	}

	makeDocumentToPage(doc, true, c, acc, b)
}

func (b *BackgroundInfo) tryAddingNewTag(doc *database.Document) bool {
	doc.Info.Post = append(doc.Info.Post, database.Posts{
		UUID:      uuid.New().String(),
		Hidden:    false,
		Submitted: time.Now().UTC(),
		Info:      b.TagText,
		Color:     b.TagColor,
	})

	err := doc.SaveChanges()
	if err != nil {
		b.Message = TagNotSuccessfulSaved
		return true
	}
	return false
}
