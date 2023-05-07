package htmlLetter

import (
	"API_MBundestag/dataLogic"
	"API_MBundestag/database"
	"API_MBundestag/help"
	"API_MBundestag/help/generics"
	"API_MBundestag/htmlHandler"
	"API_MBundestag/htmlHandler/htmlBasics"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"net/url"
)

type LetterCreatePageStruct struct {
	Message         string
	Names           []string
	SelectedAccount string
	Accounts        database.AccountList
	Letter          database.Letter
	ModMail         bool
}

type ModMailCreatePageStruct struct {
	LetterCreatePageStruct
}

func getEmtpyLetterCreateStruct(modMail bool, acc *database.Account) *LetterCreatePageStruct {
	val := LetterCreatePageStruct{}
	htmlHandler.FillOwnAccounts(&val, acc)
	htmlHandler.FillAllNotSuspendedNames(&val)
	val.ModMail = modMail
	val.Letter = database.Letter{Info: database.LetterInfo{NoSigning: true}}
	return &val
}

func getEmtpyLetterCreateStructWithMessage(modMail bool, acc *database.Account, message string) *LetterCreatePageStruct {
	val := LetterCreatePageStruct{}
	htmlHandler.FillOwnAccounts(&val, acc)
	htmlHandler.FillAllNotSuspendedNames(&val)
	val.Message = message + "\n" + val.Message
	val.ModMail = modMail
	val.Letter = database.Letter{Info: database.LetterInfo{NoSigning: true}}
	return &val
}

// GetCreateLetterPage handles gin context requests for the creation of a letter
func GetCreateLetterPage(c *gin.Context) {
	acc, b := dataLogic.CheckUserPrivileged(c, database.HeadAdmin, database.Admin, database.MediaAdmin, database.User)
	if !b {
		htmlBasics.MakeErrorPage(c, &acc, generics.NotAuthorizedToView)
		return
	}
	letterStruct := getEmtpyLetterCreateStruct(false, &acc)
	letterStruct.Letter = database.Letter{Info: database.LetterInfo{NoSigning: true}}
	htmlHandler.MakeSite(letterStruct, c, &acc)
}

// GetCreateModMailPage handles gin context requests for the creation of a mod mail
func GetCreateModMailPage(c *gin.Context) {
	acc, b := dataLogic.CheckUserPrivileged(c, database.HeadAdmin, database.Admin)
	if !b {
		htmlBasics.MakeErrorPage(c, &acc, generics.NotAuthorizedToView)
		return
	}
	modMailStruct := ModMailCreatePageStruct{*getEmtpyLetterCreateStruct(true, &acc)}
	modMailStruct.Letter = database.Letter{Info: database.LetterInfo{NoSigning: true}}
	htmlHandler.MakeSite(&modMailStruct, c, &acc)
}

func PostCreateLetterPage(c *gin.Context) {
	acc, b := dataLogic.CheckUserPrivileged(c, database.HeadAdmin, database.Admin, database.MediaAdmin, database.User)
	if !b {
		htmlBasics.MakeErrorPage(c, &acc, generics.NotAuthorizedToView)
		return
	}

	letterStruct, err := validateLetterCreate(c, &acc, false)
	if err == nil {
		c.Redirect(http.StatusFound, "/letter?uuid="+url.QueryEscape(letterStruct.Letter.UUID)+"&usr="+url.QueryEscape(letterStruct.Letter.Author))
		return
	}

	modMailStruct := ModMailCreatePageStruct{*letterStruct}
	htmlHandler.MakeSite(&modMailStruct, c, &acc)
}

func PostCreateModMailPage(c *gin.Context) {
	acc, b := dataLogic.CheckUserPrivileged(c, database.HeadAdmin, database.Admin)
	if !b {
		htmlBasics.MakeErrorPage(c, &acc, generics.NotAuthorizedToView)
		return
	}

	letterStruct, err := validateLetterCreate(c, &acc, true)
	if err == nil {
		c.Redirect(http.StatusFound, "/letter?uuid="+url.QueryEscape(letterStruct.Letter.UUID)+"&usr="+url.QueryEscape(acc.DisplayName))
		return
	}

	htmlHandler.MakeSite(letterStruct, c, &acc)
}

var ErrorInLetter = htmlHandler.ValidationErrors{Info: "ErrorInLetter"}

func validateLetterCreate(c *gin.Context, acc *database.Account, modMail bool) (letterStruct *LetterCreatePageStruct, errReturn error) {
	//Set up the function
	errReturn = ErrorInLetter
	letterStruct = getEmtpyLetterCreateStruct(modMail, acc)
	//extract all infos from the context
	letterStruct.Letter = getLetterFromContext(c)
	//letterStruct.Letter.Info.ModMessage = modMail

	letterStruct.SelectedAccount = generics.GetText(c, "selectedAccount")

	writer := &database.Account{}
	switch true {
	case letterStruct.checkAuthor():
	case generics.CheckTitelAndContentEmptyLayer(letterStruct, &letterStruct.Letter):
	case generics.CheckLengthContentLayer(letterStruct, &letterStruct.Letter, generics.LetterContentLimit):
	case generics.CheckLengthTitleLayer(letterStruct, &letterStruct.Letter, generics.LetterTitleLimit):
	case letterStruct.checkWriter(writer, acc):
	//case gen.CheckAccountList(letterStruct, &letterStruct.Letter.Info.PeopleInvitedToSign):
	default:
		letterStruct.setSigning(writer)
		errReturn = letterStruct.finishLetter(c, writer)
	}
	return
}

func getLetterFromContext(c *gin.Context) database.Letter {
	return database.Letter{
		Title:   generics.GetText(c, "title"),
		Author:  generics.GetText(c, "author"),
		Content: generics.GetText(c, "content"),
		Info: database.LetterInfo{
			AllHaveToAgree: generics.GetBool(c, "allHaveToSign"),
			NoSigning:      generics.GetBool(c, "noSigning"),
			//PeopleInvitedToSign: generics.GetStringArray(c, "user"),
		},
	}
}

func (s *LetterCreatePageStruct) checkAuthor() bool {
	if s.ModMail && (s.Letter.Author == "") {
		s.Message = generics.AuthorEmptyError + "\n" + s.Message
		return true
	}
	return false
}

func (s *LetterCreatePageStruct) checkWriter(writer *database.Account, acc *database.Account) bool {
	if s.ModMail {
		return false
	}
	return generics.CheckWriter(s, writer, acc)
}

func (s *LetterCreatePageStruct) setSigning(writer *database.Account) {
	/*if s.Letter.Info.NoSigning {
		s.Letter.Info.PeopleNotYetSigned = []string{}
		s.Letter.Info.Signed = []string{}
		s.Letter.Info.Rejected = []string{}
	} else {
		s.Letter.Info.Signed = []string{}
		s.Letter.Info.PeopleNotYetSigned = make([]string, len(s.Letter.Info.PeopleInvitedToSign))
		copy(s.Letter.Info.PeopleNotYetSigned, s.Letter.Info.PeopleInvitedToSign)
		//only sign the letter if it is a personal one, never a mod message
		if !s.ModMail {
			s.Letter.Info.Signed = []string{writer.DisplayName}
			s.Letter.Info.PeopleNotYetSigned = helper.RemoveFirstStringOccurrenceFromArray(s.Letter.Info.PeopleNotYetSigned, writer.DisplayName)
		}
		s.Letter.Info.Rejected = []string{}
	}
	if helper.GetPositionOfString(s.Letter.Info.PeopleInvitedToSign, writer.DisplayName) == -1 && writer.DisplayName != "" && !s.ModMail {
		s.Letter.Info.PeopleInvitedToSign = append(s.Letter.Info.PeopleInvitedToSign, writer.DisplayName)
	}*/
}

func (s *LetterCreatePageStruct) finishLetter(c *gin.Context, writer *database.Account) (err error) {
	//set other parameter needed
	s.Letter.HTMLContent = help.CreateHTML(s.Letter.Content)
	s.Letter.UUID = uuid.New().String()
	//set the letter parameter, either for the actual author or for the author the moderation created
	if s.ModMail {
		s.Letter.Flair = generics.GetText(c, "flair")
	} else {
		s.Letter.Author = writer.DisplayName
		s.Letter.Flair = writer.Flair
	}
	err = s.Letter.CreateMe()

	if err != nil {
		s.Message = generics.LetterCouldNotBeCreated + "\n" + s.Message
		return
	}

	err = nil
	return
}
