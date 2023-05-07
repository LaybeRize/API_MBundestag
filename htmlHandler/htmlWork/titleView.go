package htmlWork

import (
	"API_MBundestag/dataLogic"
	gen "API_MBundestag/htmlHandler"
	"github.com/gin-gonic/gin"
)

func GetTitleViewPage(c *gin.Context) {
	acc, _ := dataLogic.CheckUserPrivileged(c)
	val := dataLogic.GetTitleHierarchy()
	gen.MakeSite(&val, c, &acc)
}

//Sorry but this doesn't need a test lol
