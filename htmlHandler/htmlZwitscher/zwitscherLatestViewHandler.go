package htmlZwitscher

import (
	"API_MBundestag/dataLogic"
	"API_MBundestag/database_old"
	"API_MBundestag/help/generics"
	"API_MBundestag/htmlHandler"
	"API_MBundestag/htmlHandler/htmlBasics"
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ZwitscherListViewStruct struct {
	Zwitscher        database.ZwitscherList
	CanZwitscher     bool
	SelectedAccount  string
	Accounts         database.AccountList
	Content          string
	Message          string
	DateFormatString string
	Amount           int
}

func GetZwitscherLatestViewPage(c *gin.Context) {
	if !generics.GetIfEmptyQuery(c, "uuid") {
		GetZwitscherSingleViewPage(c)
		return
	}

	acc, b := dataLogic.CheckUserPrivileged(c, database.HeadAdmin, database.Admin, database.MediaAdmin)
	res := &ZwitscherListViewStruct{}

	i := htmlHandler.ExtractAmount(c, 1, 50, 20)
	res.Amount = i
	err := res.Zwitscher.GetLatested(i, b)
	if err != nil {
		htmlBasics.MakeErrorPage(c, &acc, generics.CouldNotLoadTweets)
		return
	}

	htmlHandler.FillOwnAccounts(res, &acc)
	res.CanZwitscher = acc.Role != database.NotLoggedIn
	res.DateFormatString = generics.ZwitscherFormat
	htmlHandler.MakeSite(res, c, &acc)
}

func PostZwitscherLatestViewPage(c *gin.Context) {
	if !generics.GetIfEmptyQuery(c, "uuid") {
		PostZwitscherSingleViewPage(c)
		return
	}

	acc, b := dataLogic.CheckUserPrivileged(c, database.HeadAdmin, database.Admin, database.MediaAdmin, database.User)
	if !b {
		htmlBasics.MakeErrorPage(c, &acc, generics.NotAuthorizedToView)
		return
	}

	res := &ZwitscherListViewStruct{CanZwitscher: true}
	htmlHandler.FillOwnAccounts(res, &acc)

	res.validateZwitscherCreate(c, &acc)

	i := htmlHandler.ExtractAmount(c, 1, 50, 20)
	res.Amount = i
	err := res.Zwitscher.GetLatested(i, acc.Role != database.User)
	if err != nil {
		res.Message = generics.CouldNotLoadTweets + "\n" + res.Message
	}
	res.DateFormatString = generics.ZwitscherFormat

	htmlHandler.MakeSite(res, c, &acc)
}

func (zwitscher *ZwitscherListViewStruct) validateZwitscherCreate(c *gin.Context, acc *database.Account) {
	zwitscher.SelectedAccount = generics.GetText(c, "selectedAccount")
	zwitscher.Content = generics.GetText(c, "content")
	writer := &database.Account{}
	switch true {
	case generics.CheckWriter(zwitscher, writer, acc):
	case generics.CheckFieldNotEmpty(zwitscher, "Content", generics.ZwitscherIsNotAllowedToBeEmpty):
	case generics.CheckLengthField(zwitscher, generics.CharacterLimitZwitscher, "Content", generics.ZwitscherIsToLong):
	case zwitscher.tryCreation(writer):
	default:
		zwitscher.Content = ""
		zwitscher.Message = generics.ZwitscherCreationSuccessful + "\n" + zwitscher.Message
	}
}

func (zwitscher *ZwitscherListViewStruct) tryCreation(writer *database.Account) bool {
	z := database.Zwitscher{
		UUID:        uuid.New().String(),
		Author:      writer.DisplayName,
		Flair:       writer.Flair,
		HTMLContent: zwitscher.Content,
		ConnectedTo: sql.NullString{Valid: false, String: ""},
	}
	err := z.CreateMe()
	if err != nil {
		zwitscher.Message = generics.ZwitscherCreationError + "\n" + zwitscher.Message
		return true
	}
	return false
}
