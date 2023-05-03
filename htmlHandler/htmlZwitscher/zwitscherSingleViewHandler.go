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
	"time"
)

type ZwitscherSingleViewStruct struct {
	Parent              database.Zwitscher
	Self                database.Zwitscher
	Zwitscher           database.ZwitscherList
	CanZwitscher        bool
	CanSuspendZwitscher bool
	SelectedAccount     string
	Accounts            database.AccountList
	Content             string
	Message             string
	DateFormatString    string
}

func getBasicFillForZwitscherSingleView(acc *database.Account) *ZwitscherSingleViewStruct {
	res := ZwitscherSingleViewStruct{}
	htmlHandler.FillOwnAccounts(&res, acc)
	switch acc.Role {
	case database.HeadAdmin, database.Admin, database.MediaAdmin:
		res.CanSuspendZwitscher = true
		res.CanZwitscher = true
	case database.User:
		res.CanSuspendZwitscher = false
		res.CanZwitscher = true
	default:
		res.CanSuspendZwitscher = false
		res.CanZwitscher = false
	}
	res.DateFormatString = generics.ZwitscherFormat
	return &res
}

func GetZwitscherSingleViewPage(c *gin.Context) {
	acc, _ := dataLogic.CheckUserPrivileged(c)
	res := getBasicFillForZwitscherSingleView(&acc)

	switch true {
	case res.getSelfTweet(c, &acc):
	case res.trySelfConnect(c, &acc):
	case res.getComments(c, &acc):
	default:
		res.displaySite(c, &acc)
	}
}

func PostZwitscherSingleViewPage(c *gin.Context) {
	acc, b := dataLogic.CheckUserPrivileged(c, database.HeadAdmin, database.Admin, database.MediaAdmin, database.User)
	if !b {
		htmlBasics.MakeErrorPage(c, &acc, generics.NotAuthorizedToView)
		return
	}
	res := getBasicFillForZwitscherSingleView(&acc)

	switch true {
	case res.getSelfTweet(c, &acc):
	case res.tryBlocking(c, &acc):
	case res.trySelfConnect(c, &acc):
	case res.tryNotBlocked(c, &acc):
	default:
		res.displaySite(c, &acc)
	}
}

func (zwitscher *ZwitscherSingleViewStruct) getSelfTweet(c *gin.Context, acc *database.Account) bool {
	err := zwitscher.Self.GetByID(c.Query("uuid"))
	if err != nil {
		htmlBasics.MakeErrorPage(c, acc, generics.CouldNotFindTweet)
		return true
	}
	return false
}

func (zwitscher *ZwitscherSingleViewStruct) tryBlocking(c *gin.Context, acc *database.Account) bool {
	if htmlHandler.GetBool(c, "block") {
		if acc.Role == database.User {
			htmlBasics.MakeErrorPage(c, acc, generics.NotAuthorizedToView)
			return true
		}
		zwitscher.Self.Blocked = !zwitscher.Self.Blocked
		err := zwitscher.Self.SaveChanges()
		if err != nil {
			zwitscher.Message = generics.TweetCouldNotBeBlockedOrDeblocked + zwitscher.Message
			zwitscher.Self.Blocked = !zwitscher.Self.Blocked
		}
	}
	return false
}

func (zwitscher *ZwitscherSingleViewStruct) trySelfConnect(c *gin.Context, acc *database.Account) bool {
	if zwitscher.Self.ConnectedTo.Valid {
		err := zwitscher.Parent.GetByID(zwitscher.Self.ConnectedTo.String)
		if err != nil {
			htmlBasics.MakeErrorPage(c, acc, generics.CouldNotFindTweet)
			return true
		}
	}
	return false
}

func (zwitscher *ZwitscherSingleViewStruct) tryNotBlocked(c *gin.Context, acc *database.Account) bool {
	if !htmlHandler.GetBool(c, "block") {
		zwitscher.validateMakeComment(c, acc)
	}

	return zwitscher.getComments(c, acc)
}

func (zwitscher *ZwitscherSingleViewStruct) getComments(c *gin.Context, acc *database.Account) bool {
	err := zwitscher.Zwitscher.GetCommentsFor(zwitscher.Self.UUID, zwitscher.CanSuspendZwitscher)
	if err != nil {
		htmlBasics.MakeErrorPage(c, acc, generics.CouldNotFindTweet)
		return true
	}
	return false
}

func (zwitscher *ZwitscherSingleViewStruct) displaySite(c *gin.Context, acc *database.Account) {
	if zwitscher.Self.Blocked && !zwitscher.CanSuspendZwitscher {
		zwitscher.Self.HTMLContent = generics.ZwitscherBlockText
	}
	if !zwitscher.CanSuspendZwitscher && zwitscher.Self.ConnectedTo.Valid && zwitscher.Parent.Blocked {
		zwitscher.Parent.HTMLContent = generics.ZwitscherBlockText
	}

	htmlHandler.MakeSite(zwitscher, c, acc)
}

func (zwitscher *ZwitscherSingleViewStruct) validateMakeComment(c *gin.Context, acc *database.Account) {
	zwitscher.SelectedAccount = c.PostForm("selectedAccount")
	zwitscher.Content = c.PostForm("content")
	writer := &database.Account{}
	switch true {
	case htmlHandler.CheckWriter(zwitscher, writer, acc):
	case htmlHandler.CheckFieldNotEmpty(zwitscher, "Content", generics.ZwitscherIsNotAllowedToBeEmpty):
	case htmlHandler.CheckLengthField(zwitscher, generics.CharacterLimitZwitscher, "Content", generics.ZwitscherIsToLong):
	case zwitscher.tryCreation(writer):
	default:
		zwitscher.Content = ""
		zwitscher.Message = generics.ZwitscherCreationSuccessful + "\n" + zwitscher.Message
	}
}

func (zwitscher *ZwitscherSingleViewStruct) tryCreation(writer *database.Account) bool {
	z := database.Zwitscher{
		UUID:        uuid.New().String(),
		Written:     time.Now().UTC(),
		Author:      writer.DisplayName,
		Flair:       writer.Flair,
		HTMLContent: zwitscher.Content,
		ConnectedTo: sql.NullString{Valid: true, String: zwitscher.Self.UUID},
	}
	err := z.CreateMe()
	if err != nil {
		zwitscher.Message = generics.ZwitscherCreationError + "\n" + zwitscher.Message
		return true
	}
	return false
}
