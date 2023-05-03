package htmlWork

import (
	"API_MBundestag/dataLogic"
	"API_MBundestag/database_old"
	"API_MBundestag/htmlHandler"
	"API_MBundestag/htmlHandler/htmlBasics"
	wr "API_MBundestag/htmlWrapper"
	"database/sql"
	"github.com/stretchr/testify/assert"
	"html/template"
	"testing"
)

func TestGetTitleViewPage(t *testing.T) {
	helper.TestBlocker.Lock()
	defer helper.TestBlocker.Unlock()
	Setup()
	htmlBasics.Setup()

	dataLogic.TitleHierarchy = dataLogic.TitleMainGroupArray{{
		Name: "Test1",
		Groups: []dataLogic.TitleSubGroup{
			{
				Name: "asd",
				Titles: []database.Title{
					{
						Name:      "asd",
						MainGroup: "bxc",
						SubGroup:  "eway",
						Flair:     sql.NullString{Valid: true, String: "bazinga"},
						Info:      database.TitleInfo{Names: []string{"asda", "svdcasd"}},
					},
				},
			},
		},
	}}
	temp, err := template.New("layout").Funcs(wr.DefaultFunctions).Parse("{{.Template}} {{range $i, $main := .Page}}{{$main.Name}} {{range $j, $sub := $main.Groups}}{{$sub.Name}} {{range $k, $title := $sub.Titles}}{{$title.Name}} {{range $z, $name := $title.Info.Names}} {{$name}}{{end}}{{end}}{{end}}{{end}}")
	assert.Nil(t, err)
	htmlHandler.Template = &wr.Templates{
		Extension: "",
		Dir:       "",
		Templates: map[string]*template.Template{
			"displayTitle": temp,
		},
	}
	w, ctx := htmlHandler.GetEmptyContext(t)
	GetTitleViewPage(ctx)
	assert.Equal(t, "displayTitle Test1 asd asd  asda svdcasd", w.Body.String())
}
