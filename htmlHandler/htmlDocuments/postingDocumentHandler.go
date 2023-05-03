package htmlDocuments

import (
	"API_MBundestag/dataLogic"
	"API_MBundestag/database_old"
	"API_MBundestag/help/generics"
	"API_MBundestag/htmlHandler"
	"API_MBundestag/htmlHandler/htmlBasics"
	"github.com/gin-gonic/gin"
)

type DocumentNavigationStruct struct {
	Accounts           database.AccountList
	Organisations      database.OrganisationList
	AccountSelect      bool
	OrgAlreadySelected bool
	OrganisationSelect bool
	DocumentSelect     bool
	CanPost            bool
	CanDiscussOrVote   bool
	OrgName            string
	AccName            string
	Message            string
}

func GetDocumentNavigationPage(c *gin.Context) {
	acc, b := dataLogic.CheckUserPrivileged(c, database.HeadAdmin, database.Admin, database.MediaAdmin, database.User)
	if !b {
		htmlBasics.MakeErrorPage(c, &acc, generics.NotAuthorizedToView)
		return
	}

	res := DocumentNavigationStruct{}

	res.getCorrectFunction(c, &acc)

	htmlHandler.MakeSite(&res, c, &acc)
}

func (doc *DocumentNavigationStruct) getCorrectFunction(c *gin.Context, acc *database.Account) {
	switch true {
	case generics.GetIfEmptyQuery(c, "usr") && generics.GetIfEmptyQuery(c, "org"):
		doc.selectFromAllYourAccounts(acc)
	case generics.GetIfEmptyQuery(c, "usr") && !generics.GetIfEmptyQuery(c, "org"):
		doc.selectFromAccountsForOrg(c, acc)
	case !generics.GetIfEmptyQuery(c, "usr") && generics.GetIfEmptyQuery(c, "org"):
		doc.selectFromAccountOrganisation(c, acc)
	case !generics.GetIfEmptyQuery(c, "usr") && !generics.GetIfEmptyQuery(c, "org"):
		doc.selectEverythingYouCanDoWithTheOrg(c, acc)
	}
}

func (doc *DocumentNavigationStruct) selectFromAllYourAccounts(acc *database.Account) {
	doc.AccountSelect = true
	doc.OrgAlreadySelected = false
	doc.OrganisationSelect = false
	doc.DocumentSelect = false
	htmlHandler.FillOwnAccounts(doc, acc)
}

func (doc *DocumentNavigationStruct) selectFromAccountsForOrg(c *gin.Context, acc *database.Account) {
	doc.AccountSelect = true
	doc.OrgAlreadySelected = true
	doc.OrganisationSelect = false
	doc.DocumentSelect = false
	doc.OrgName = c.Query("org")
	htmlHandler.FillOwnAccounts(doc, acc)
	org := database.Organisation{}
	err := org.GetByName(doc.OrgName)
	if err != nil {
		doc.Message = generics.OrganisationInURLDoesNotExist + "\n" + doc.Message
		doc.Accounts = database.AccountList{}
		return
	}
	newList := database.AccountList{}
	for _, rangeAcc := range doc.Accounts {
		if helper.GetPositionOfString(org.Info.User, rangeAcc.DisplayName) != -1 || helper.GetPositionOfString(org.Info.Admins, rangeAcc.DisplayName) != -1 {
			newList = append(newList, rangeAcc)
		}
	}
	doc.Accounts = newList
}

func (doc *DocumentNavigationStruct) selectFromAccountOrganisation(c *gin.Context, acc *database.Account) {
	doc.AccountSelect = false
	doc.OrganisationSelect = true
	doc.DocumentSelect = false
	doc.AccName = c.Query("usr")
	author := database.Account{}
	err := author.GetByDisplayName(doc.AccName)
	if err != nil {
		doc.Message = generics.AccountDoesNotExistOrIsNotYours + "\n" + doc.Message
		doc.Organisations = database.OrganisationList{}
		return
	}
	if (author.Linked.Int64 != acc.ID || author.Suspended) && !(author.DisplayName == acc.DisplayName) {
		doc.Message = generics.AccountDoesNotExistOrIsNotYours + "\n" + doc.Message
		doc.Organisations = database.OrganisationList{}
		return
	}

	htmlHandler.FillOwnOrganisations(doc, &author)
}

func (doc *DocumentNavigationStruct) selectEverythingYouCanDoWithTheOrg(c *gin.Context, acc *database.Account) {
	doc.AccountSelect = false
	doc.OrganisationSelect = false
	doc.DocumentSelect = true
	doc.CanPost = false
	doc.CanDiscussOrVote = false
	doc.AccName = c.Query("usr")
	doc.OrgName = c.Query("org")
	author := database.Account{}
	err := author.GetByDisplayName(doc.AccName)
	if err != nil {
		doc.Message = generics.AccountDoesNotExistOrIsNotYours + "\n" + doc.Message
		doc.Organisations = database.OrganisationList{}
		return
	}
	if (author.Linked.Int64 != acc.ID || author.Suspended) && !(author.DisplayName == acc.DisplayName) {
		doc.Message = generics.AccountDoesNotExistOrIsNotYours + "\n" + doc.Message
		doc.Organisations = database.OrganisationList{}
		return
	}

	org := database.Organisation{}
	err = org.GetByName(doc.OrgName)
	if err != nil {
		doc.Message = generics.OrganisationInURLDoesNotExist + "\n" + doc.Message
		return
	}

	if helper.GetPositionOfString(org.Info.User, author.DisplayName) != -1 {
		doc.CanDiscussOrVote = true
	}

	if helper.GetPositionOfString(org.Info.Admins, author.DisplayName) != -1 || author.Role == database.HeadAdmin {
		doc.CanDiscussOrVote = true
		doc.CanPost = true
	}
}
