package htmlAccount

import (
	"API_MBundestag/dataLogic"
	"API_MBundestag/database"
	"API_MBundestag/help"
	"API_MBundestag/help/generics"
	"API_MBundestag/htmlHandler"
	"API_MBundestag/htmlHandler/htmlBasics"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"strconv"
)

type CreateUserStruct struct {
	Account database.Account
	help.MessageStruct
}

func GetCreateUserPage(c *gin.Context) {
	acc, b := dataLogic.CheckUserPrivileged(c, database.HeadAdmin)
	if !b {
		htmlBasics.MakeErrorPage(c, &acc, generics.NotAuthorizedToView)
		return
	}

	createStruct := &CreateUserStruct{Account: database.Account{Role: database.User}}
	htmlHandler.MakeSite(createStruct, c, &acc)
}

func PostCreateUserPage(c *gin.Context) {
	acc, b := dataLogic.CheckUserPrivileged(c, database.HeadAdmin)
	if !b {
		htmlBasics.MakeErrorPage(c, &acc, generics.NotAuthorizedToView)
		return
	}
	htmlHandler.MakeSite(validateCreateAccount(&acc, c), c, &acc)
}

func validateCreateAccount(self *database.Account, c *gin.Context) (createStruct *CreateUserStruct) {
	createStruct = &CreateUserStruct{Account: database.Account{
		DisplayName: generics.GetText(c, "displayname"),
		Flair:       generics.GetText(c, "flair"),
		Username:    generics.GetText(c, "username"),
		Password:    generics.GetText(c, "password"),
		Role:        database.User,
	}}
	result := &createStruct.Account
	hash := &[]byte{}

	switch true {
	case createStruct.linkValue(c):
	case createStruct.getRole(c, self):
	case createStruct.generatePassword(c, hash):
	case createStruct.checkNamesAndPassword(result):
	default:
		createStruct.finishCreation(result, hash)
	}

	return
}

func (s *CreateUserStruct) linkValue(c *gin.Context) bool {
	i, err := strconv.Atoi(generics.GetText(c, "linked"))
	if err != nil {
		s.Message = generics.LinkedValueNotANumberError
		return true
	}
	s.Account.Linked.Int64 = int64(i)
	return false
}

func (s *CreateUserStruct) getRole(c *gin.Context, self *database.Account) bool {
	//Add Head-Admin if the creator is the root account
	array := database.Roles[:len(database.Roles)-1]
	if self.ID == 1 {
		array = database.Roles
	}
	//Check if the role exists
	r := generics.GetText(c, "role")
	if help.GetPositionOfString(array, r) == -1 {
		s.Message = generics.RoleCanNotBeSelectedError
		return true
	}
	//Set role and set if the linked value is valid
	s.Account.Role = database.RoleString(r)
	s.Account.Linked.Valid = database.PressAccount == s.Account.Role
	return false
}

func (s *CreateUserStruct) generatePassword(c *gin.Context, hash *[]byte) bool {
	var err error
	*hash, err = bcrypt.GenerateFromPassword([]byte(generics.GetText(c, "password")), bcrypt.DefaultCost)
	if err != nil {
		s.Message = generics.ErrorWhileGeneratingPasswordHash
		return true
	}
	return false
}

func (s *CreateUserStruct) checkNamesAndPassword(result *database.Account) bool {
	if result.Username == "" && result.Role == database.PressAccount {
		result.Username = result.DisplayName
	}

	if result.Username == "" || result.DisplayName == "" || (result.Password == "" && result.Role != database.PressAccount) {
		s.Message = generics.NamesOrPasswordIsEmptyError
		return true
	}
	return false
}

func (s *CreateUserStruct) finishCreation(result *database.Account, hash *[]byte) {
	if result.Role == database.PressAccount {
		result.Password = ""
	} else {
		result.Password = string(*hash)
	}

	err := result.CreateMe()
	result.Password = ""
	if err != nil {
		s.Message = generics.UserOrDisplaynameAlreadyExist
		return
	}

	s.Positiv = true
	s.Message = generics.SuccesFullCreatedAccount
}
