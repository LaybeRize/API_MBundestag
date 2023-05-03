package htmlBasics

import (
	"API_MBundestag/help"
	"testing"
)

func TestMarkdownHandler(t *testing.T) {
	help.UpdateAttributes()

	t.Run("testCorrectHTMLResponse", testCorrectHTMLResponse)
}

func testCorrectHTMLResponse(t *testing.T) {

}
