package htmlBasics

import (
	"API_MBundestag/database"
	gen "API_MBundestag/htmlHandler"
	"github.com/gin-gonic/gin"
)

type ErrorCode struct {
	Error string
}

func MakeErrorPage(c *gin.Context, acc *database.Account, errorText string) {
	gen.MakeSite(&ErrorCode{Error: errorText}, c, acc)
}
