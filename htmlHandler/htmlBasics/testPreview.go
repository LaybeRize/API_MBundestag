package htmlBasics

import (
	"API_MBundestag/dataLogic"
	"API_MBundestag/help"
	gen "API_MBundestag/htmlHandler"
	"github.com/gin-gonic/gin"
	"html/template"
)

type TestPreviewStruct struct {
	Text    string
	Preview template.HTML
}

func GetPreviewPage(c *gin.Context) {
	acc, _ := dataLogic.CheckUserPrivileged(c)
	gen.MakeSite(&TestPreviewStruct{}, c, &acc)
}

func PostPreviewPage(c *gin.Context) {
	acc, _ := dataLogic.CheckUserPrivileged(c)
	text := c.PostForm("text")
	html := help.CreateHTML(text)
	gen.MakeSite(&TestPreviewStruct{
		Text:    text,
		Preview: template.HTML(html),
	}, c, &acc)
}
