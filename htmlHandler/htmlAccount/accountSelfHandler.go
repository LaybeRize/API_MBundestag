package htmlAccount

import (
	"API_MBundestag/dataLogic"
	"API_MBundestag/database"
	"API_MBundestag/help/generics"
	"API_MBundestag/htmlHandler"
	"API_MBundestag/htmlHandler/htmlBasics"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type ViewUserInfoStruct []ViewUserInfoElement

type ViewUserInfoElement struct {
	DisplayName  string
	Flair        string
	Title        string
	Organisation string
}

type PasswordChangeStruct struct {
	htmlHandler.MessageStruct
}

func GetPasswordChangePage(c *gin.Context) {
	acc, b := dataLogic.CheckUserPrivileged(c, database.HeadAdmin, database.Admin, database.MediaAdmin, database.User)
	if !b {
		htmlBasics.MakeErrorPage(c, &acc, generics.NotAuthorizedToView)
		return
	}

	htmlHandler.MakeSite(&PasswordChangeStruct{}, c, &acc)
}

func PostPasswordChangePage(c *gin.Context) {
	acc, b := dataLogic.CheckUserPrivileged(c, database.HeadAdmin, database.Admin, database.MediaAdmin, database.User)
	if !b {
		htmlBasics.MakeErrorPage(c, &acc, generics.NotAuthorizedToView)
		return
	}

	htmlHandler.MakeSite(validateChangePassword(c, &acc), c, &acc)
}

func validateChangePassword(c *gin.Context, acc *database.Account) (changeStruct *PasswordChangeStruct) {
	changeStruct = &PasswordChangeStruct{}
	changeStruct.Message = generics.NewPasswordIsNotTheSame
	newPassword := c.PostForm("newPassword")
	if newPassword != c.PostForm("newPassword2") {
		return
	}

	if len([]rune(newPassword)) < generics.MinPasswordLength {
		changeStruct.Message = generics.NewPasswordIsNotMinimumOf10Characters
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(acc.Password), []byte(c.PostForm("password")))
	if err != nil {
		changeStruct.Message = generics.OldPasswordNotcorrect
		return
	}

	changeStruct.Message = generics.ErrorWhileChangingPassword
	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return
	}

	acc.Password = string(hash)
	err = acc.SaveChanges()
	if err != nil {
		return
	}

	changeStruct.Message = generics.SuccessChangePassword
	changeStruct.Positiv = true
	return
}

func GetViewOfProfilePage(c *gin.Context) {
	acc, b := dataLogic.CheckUserPrivileged(c, database.HeadAdmin, database.Admin, database.MediaAdmin, database.User)
	if !b {
		htmlBasics.MakeErrorPage(c, &acc, generics.NotAuthorizedToView)
		return
	}

	accounts := database.AccountList{}
	err := accounts.GetAllPressAccountsFromAccountPlusSelf(&acc)
	if err != nil {
		htmlBasics.MakeErrorPage(c, &acc, generics.CouldNotLoadAccountDetails)
		return
	}

	err, viewStruct := getViewStruct(accounts)
	if err != nil {
		htmlBasics.MakeErrorPage(c, &acc, generics.CouldNotLoadAccountDetails)
		return
	}

	htmlHandler.MakeSite(viewStruct, c, &acc)
}

func getViewStruct(accounts database.AccountList) (err error, view *ViewUserInfoStruct) {
	view = &ViewUserInfoStruct{}
	for _, acc := range accounts {
		var element ViewUserInfoElement
		err, element = getElementForUser(acc)
		if err != nil {
			return
		}
		*view = append(*view, element)
	}
	return
}

func getElementForUser(acc database.Account) (err error, result ViewUserInfoElement) {
	result = ViewUserInfoElement{
		DisplayName: acc.DisplayName,
		Flair:       acc.Flair,
	}
	err, result.Title = dataLogic.GetTitelList(acc.ID)
	if err != nil {
		return
	}
	err, result.Organisation = dataLogic.GetOrganisationList(acc.ID)
	return
}
