package htmlBasics

import (
	"API_MBundestag/help"
	"github.com/gin-gonic/gin"
)

type MarkdownRequest struct {
	Markdown string `json:"markdown"`
}

type HtmlResponse struct {
	Html string `json:"html"`
}

func PostJsonMarkdown(c *gin.Context) {
	var markdownRequest MarkdownRequest

	if err := c.BindJSON(&markdownRequest); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	html := help.CreateHTML(markdownRequest.Markdown)
	response := HtmlResponse{Html: html}
	c.JSON(200, response)
}
