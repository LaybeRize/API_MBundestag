package htmlAccount

import (
	"API_MBundestag/dataLogic"
	"API_MBundestag/database"
	"API_MBundestag/help"
	"API_MBundestag/help/generics"
	"API_MBundestag/htmlHandler"
	"API_MBundestag/htmlHandler/htmlBasics"
	"github.com/gin-gonic/gin"
)

type ListUserStruct struct {
	Accounts database.AccountList
	help.MessageStruct
}

func GetAdminListUserPage(c *gin.Context) {
	acc, b := dataLogic.CheckUserPrivileged(c, database.HeadAdmin)
	if !b {
		htmlBasics.MakeErrorPage(c, &acc, generics.NotAuthorizedToView)
		return
	}

	htmlHandler.MakeSite(validateListOfAccounts(c), c, &acc)
}

func validateListOfAccounts(c *gin.Context) (listAccounts *ListUserStruct) {
	listAccounts = &ListUserStruct{}
	acc := database.Account{}
	err := acc.GetByDisplayName(c.Query("acc"))

	//if the account does not exist, the query was empty, or the account is a press account
	if c.Query("acc") == "" || err != nil || acc.Role == database.PressAccount {
		err = listAccounts.Accounts.GetAllAccounts()
	} else {
		err = listAccounts.Accounts.GetAllPressAccountsFromAccountPlusSelf(&acc)
	}

	//if the accounts could not be loaded, put a message on it
	if err != nil {
		listAccounts.Message = generics.AccountQueryError
	}
	return
}
