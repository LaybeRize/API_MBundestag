package htmlBasics

import (
	"API_MBundestag/help"
	"API_MBundestag/htmlHandler"
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io"
	"testing"
)

func TestMarkdownHandler(t *testing.T) {
	htmlHandler.ChangePath()

	help.UpdateAttributes()

	t.Run("testCorrectHTMLResponse", testCorrectHTMLResponse)
}

func testCorrectHTMLResponse(t *testing.T) {
	w, ctx := htmlHandler.GetEmptyContext(t)
	ctx.Request.Method = "POST" // or PUT
	ctx.Request.Header.Set("Content-Type", "application/json")
	jsonbytes, err := json.Marshal(MarkdownRequest{Markdown: "*test* and **test**"})
	if err != nil {
		assert.Fail(t, "there should be no marshal error")
	}
	ctx.Request.Body = io.NopCloser(bytes.NewBuffer(jsonbytes))
	PostJsonMarkdown(ctx)

	var dat map[string]interface{}

	if err = json.Unmarshal([]byte(w.Body.String()), &dat); err != nil {
		panic(err)
	}
	assert.Equal(t, `<p class="text-justify break-words"><em>test</em> and <strong>test</strong></p>
`, dat["html"].(string))
}
