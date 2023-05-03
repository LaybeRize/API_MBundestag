package htmlAccount

import (
	"API_MBundestag/dataLogic"
	"API_MBundestag/database"
	"API_MBundestag/help"
	"API_MBundestag/help/generics"
	"API_MBundestag/htmlHandler"
	"API_MBundestag/htmlHandler/htmlBasics"
	"github.com/gin-gonic/gin"
	"strconv"
)

type EditUserStruct struct {
	Account dataLogic.Account
	Names   database.NameList
	help.MessageStruct
}

func getEmptyEditUserStruct() *EditUserStruct {
	result := &EditUserStruct{}
	htmlHandler.FillUserAndDisplayNames(result)
	return result
}

func GetEditUserPage(c *gin.Context) {
	acc, b := dataLogic.CheckUserPrivileged(c, database.HeadAdmin)
	if !b {
		htmlBasics.MakeErrorPage(c, &acc, generics.NotAuthorizedToView)
		return
	}

	htmlHandler.MakeSite(getEmptyEditUserStruct(), c, &acc)
}

func PostEditUserPage(c *gin.Context) {
	acc, b := dataLogic.CheckUserPrivileged(c, database.HeadAdmin)
	if !b {
		htmlBasics.MakeErrorPage(c, &acc, generics.NotAuthorizedToView)
		return
	}

	if c.Query("change") == "true" {
		htmlHandler.MakeSite(validateChangeAccount(c, &acc), c, &acc)
		return
	}

	htmlHandler.MakeSite(validateGetAccount(c), c, &acc)
}

func validateGetAccount(c *gin.Context) *EditUserStruct {
	result := getEmptyEditUserStruct()
	if c.Query("type") == "user" {
		result.Account.GetUser("", c.PostForm("name"), &result.Message, &result.Positiv)
	} else if c.Query("type") == "display" {
		result.Account.GetUser(c.PostForm("name"), "", &result.Message, &result.Positiv)
	} else {
		result.Message = generics.InvalidType + "\n" + result.Message
	}

	return result
}

func validateChangeAccount(c *gin.Context, self *database.Account) (editStruct *EditUserStruct) {
	editStruct = getEmptyEditUserStruct()

	switch true {
	case editStruct.setAccount(c):
	case editStruct.checkIfOnlyFlair(c, self):
	case editStruct.checkRolePrivileges(self):
	case editStruct.setFlairAndCheckLinked(c):
	case editStruct.changeRoleCheck(c, self):
	default:
		editStruct.Account.ChangeUser(&editStruct.Message, &editStruct.Positiv)
	}
	return
}

func (s *EditUserStruct) setAccount(c *gin.Context) bool {
	var temp help.Message
	s.Account.GetUser("", generics.GetText(c, "username"), &temp, &s.Positiv)

	if !s.Positiv {
		s.Message = generics.CanNotChangeNoExistentAccount + "\n" + s.Message
		return true
	}
	s.Positiv = false
	return false
}

func (s *EditUserStruct) checkIfOnlyFlair(c *gin.Context, self *database.Account) bool {
	switch true {
	case s.Account.ID != 1 && self.ID == 1:
		return false
	case !(s.Account.ID == self.ID):
		return false
	}

	s.Account.Flair = generics.GetText(c, "flair")
	s.Account.ChangeFlair = generics.GetBool(c, "changeFlair")
	s.Account.ChangeUser(&s.Message, &s.Positiv)
	return true
}

func (s *EditUserStruct) checkRolePrivileges(self *database.Account) bool {
	if s.Account.ID == 1 {
		s.Message = generics.CanNotChangeRootAccount + "\n" + s.Message
		return true
	}

	if s.Account.Role == database.HeadAdmin && self.ID != 1 {
		s.Message = generics.DisallowedChangeToHeadAdmin + "\n" + s.Message
		return true
	}
	return false
}

func (s *EditUserStruct) setFlairAndCheckLinked(c *gin.Context) bool {
	s.Account.Flair = generics.GetText(c, "flair")
	s.Account.ChangeFlair = generics.GetBool(c, "changeFlair")
	s.Account.Suspended = generics.GetBool(c, "suspended")
	s.Account.RemoveFromTitle = generics.GetBool(c, "removeTitles")
	s.Account.RemoveFromOrganisation = generics.GetBool(c, "removeOrgs")

	i, err := strconv.Atoi(generics.GetText(c, "linked"))
	if err != nil {
		s.Message = generics.LinkedValueNotANumberError + "\n" + s.Message
		return true
	}

	s.Account.Linked = int64(i)
	return false
}

func (s *EditUserStruct) changeRoleCheck(c *gin.Context, self *database.Account) bool {
	if s.Account.Role == database.PressAccount {
		return false
	}

	array := database.Roles[1 : len(database.Roles)-1]
	if self.ID == 1 {
		array = append(array, string(database.HeadAdmin))
	}
	r := generics.GetText(c, "role")
	if help.GetPositionOfString(array, r) == -1 {
		s.Message = generics.RoleCanNotBeSelectedError + "\n" + s.Message
		return true
	}
	s.Account.Role = database.RoleString(r)

	return false
}
