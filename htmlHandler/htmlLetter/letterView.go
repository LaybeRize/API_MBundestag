package htmlLetter

import (
	"API_MBundestag/dataLogic"
	"API_MBundestag/database_old"
	"API_MBundestag/help/generics"
	"API_MBundestag/htmlHandler"
	"API_MBundestag/htmlHandler/htmlBasics"
	"github.com/gin-gonic/gin"
)

type ViewSingleLetter struct {
	Letter       database.Letter
	Account      database.Account
	FormatString string
}

func GetViewSingleLetter(c *gin.Context) {
	acc, b := dataLogic.CheckUserPrivileged(c, database.HeadAdmin, database.Admin, database.MediaAdmin, database.User)
	if !b {
		htmlBasics.MakeErrorPage(c, &acc, generics.NotAuthorizedToView)
		return
	}

	//check viewer
	viewer := database.Account{}
	letter := database.Letter{}
	switch true {
	case getViewer(c, &viewer, &acc):
		return
	case getLetter(c, &letter, &viewer, &acc):
		return
	}
	//setup letter
	viewLetter := ViewSingleLetter{
		Letter:       letter,
		Account:      viewer,
		FormatString: generics.LongTimeString,
	}
	switch true {
	case generics.GetIfType(c, "sign"):
		viewLetter.signLetter(viewer)
	case generics.GetIfType(c, "reject"):
		viewLetter.rejectLetter(viewer)
	}

	htmlHandler.MakeSite(&viewLetter, c, &acc)
}

func getLetter(c *gin.Context, letter *database.Letter, viewer *database.Account, acc *database.Account) bool {
	err := letter.GetByID(c.Query("uuid"))
	//if the user is a moderator, and it's a modMail, just display it, no question asked
	b := (acc.Role == database.HeadAdmin || acc.Role == database.Admin) && letter.Info.ModMessage
	for i := 0; i < len(letter.Info.PeopleInvitedToSign) && !b; i++ {
		if letter.Info.PeopleInvitedToSign[i] == viewer.DisplayName {
			b = true
		}
	}
	//handel errors
	if !b || err != nil {
		htmlBasics.MakeErrorPage(c, acc, generics.LetterDoesntExistOrNotAccessable)
		return true
	}
	return false
}

func getViewer(c *gin.Context, viewer *database.Account, acc *database.Account) bool {
	err := viewer.GetByDisplayName(c.Query("usr"))
	if err != nil && !generics.GetIfEmptyQuery(c, "usr") {
		htmlBasics.MakeErrorPage(c, acc, generics.AccountForLetterViewError)
		return true
	}
	if generics.GetIfEmptyQuery(c, "usr") {
		viewer = acc
	}
	if viewer.DisplayName != acc.DisplayName && viewer.Linked.Int64 != acc.ID {
		htmlBasics.MakeErrorPage(c, acc, generics.AccountForLetterViewError)
		return true
	}
	return false
}

func (viewLetter *ViewSingleLetter) rejectLetter(viewer database.Account) {
	letter := viewLetter.Letter
	if helper.GetPositionOfString(letter.Info.PeopleNotYetSigned, viewer.DisplayName) == -1 {
		return
	}
	letter.Info.PeopleNotYetSigned = helper.RemoveFirstStringOccurrenceFromArray(letter.Info.PeopleNotYetSigned, viewer.DisplayName)
	letter.Info.Rejected = append(letter.Info.Rejected, viewer.DisplayName)
	err := letter.SaveChanges()
	if err == nil {
		viewLetter.Letter = letter
	}
}

func (viewLetter *ViewSingleLetter) signLetter(viewer database.Account) {
	letter := viewLetter.Letter
	if helper.GetPositionOfString(letter.Info.PeopleNotYetSigned, viewer.DisplayName) == -1 {
		return
	}
	letter.Info.PeopleNotYetSigned = helper.RemoveFirstStringOccurrenceFromArray(letter.Info.PeopleNotYetSigned, viewer.DisplayName)
	letter.Info.Signed = append(letter.Info.Signed, viewer.DisplayName)
	err := letter.SaveChanges()
	if err == nil {
		viewLetter.Letter = letter
	}
}

func getLetterWithoutAccount(uuid string) (viewStruct *ViewSingleLetter, err error) {
	viewStruct = &ViewSingleLetter{
		FormatString: generics.LongTimeString,
	}
	err = viewStruct.Letter.GetByID(uuid)
	return
}
