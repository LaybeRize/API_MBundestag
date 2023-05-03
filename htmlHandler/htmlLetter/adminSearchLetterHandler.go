package htmlLetter

import (
	"API_MBundestag/dataLogic"
	"API_MBundestag/database_old"
	"API_MBundestag/help/generics"
	"API_MBundestag/htmlHandler"
	"API_MBundestag/htmlHandler/htmlBasics"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/url"
)

type AdminLetterViewStruct struct {
	UUID    string
	Message string
}

func GetAdminLetterViewPage(c *gin.Context) {
	acc, b := dataLogic.CheckUserPrivileged(c, database.HeadAdmin, database.Admin)
	if !b {
		htmlBasics.MakeErrorPage(c, &acc, generics.NotAuthorizedToView)
		return
	}

	if generics.GetIfEmptyQuery(c, "uuid") {
		htmlHandler.MakeSite(&AdminLetterViewStruct{}, c, &acc)
		return
	}

	letterStruct, err := getLetterWithoutAccount(c.Query("uuid"))
	if err == nil {
		htmlHandler.MakeSite(letterStruct, c, &acc)
		return
	}

	htmlHandler.MakeSite(&AdminLetterViewStruct{
		UUID:    c.Query("uuid"),
		Message: generics.ErrorUUIDDoesNotExist,
	}, c, &acc)
}

func PostAdminLetterViewPage(c *gin.Context) {
	acc, b := dataLogic.CheckUserPrivileged(c, database.HeadAdmin, database.Admin)
	if !b {
		htmlBasics.MakeErrorPage(c, &acc, generics.NotAuthorizedToView)
		return
	}

	_, err := getLetterWithoutAccount(c.PostForm("uuid"))
	if err == nil {
		c.Redirect(http.StatusFound, "/admin-letter-view?uuid="+url.QueryEscape(c.PostForm("uuid")))
		return
	}

	htmlHandler.MakeSite(&AdminLetterViewStruct{
		UUID:    c.PostForm("uuid"),
		Message: generics.ErrorUUIDDoesNotExist,
	}, c, &acc)
}
