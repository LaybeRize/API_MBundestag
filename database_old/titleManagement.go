package database

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
)

var TitlesSchema = `
CREATE TABLE IF NOT EXISTS titles (
    name text UNIQUE NOT NULL,
    main_group TEXT NOT NULL,
    sub_group TEXT NOT NULL,
    flair TEXT UNIQUE,
    info jsonb NOT NULL
);
`

func TestTitlesDB() {
	TestDatabase("DROP TABLE IF EXISTS titles;", "")
	InitTitlesDatabase()
}

type (
	MainGroupListTitle []Title
	SubGroupListTitle  []Title
	TitleList          []Title
	Title              struct {
		Name      string
		MainGroup string `db:"main_group"`
		SubGroup  string `db:"sub_group"`
		Flair     sql.NullString
		Info      TitleInfo
	}
	TitleInfo struct {
		Names []string `json:"names"`
	}
)

//Sort Functions

func (titleList *TitleList) Len() int {
	return len(*titleList)
}

func (titleList *TitleList) Less(i, j int) bool {
	strI := (*titleList)[i].Name
	strJ := (*titleList)[j].Name
	return strI < strJ
}

func (titleList *TitleList) Swap(i, j int) {
	(*titleList)[i], (*titleList)[j] = (*titleList)[j], (*titleList)[i]
}

func (titleList SubGroupListTitle) Len() int {
	return len(titleList)
}

func (titleList SubGroupListTitle) Less(i, j int) bool {
	strI := titleList[i].SubGroup
	strJ := titleList[j].SubGroup
	return strI < strJ
}

func (titleList SubGroupListTitle) Swap(i, j int) {
	titleList[i], titleList[j] = titleList[j], titleList[i]
}

func (titleList MainGroupListTitle) Len() int {
	return len(titleList)
}

func (titleList MainGroupListTitle) Less(i, j int) bool {
	strI := titleList[i].MainGroup
	strJ := titleList[j].MainGroup
	return strI < strJ
}

func (titleList MainGroupListTitle) Swap(i, j int) {
	titleList[i], titleList[j] = titleList[j], titleList[i]
}

//Sort Functions End

func InitTitlesDatabase() {
	DB.MustExec(TitlesSchema)
}

func (ti *TitleInfo) Scan(val interface{}) error {
	switch v := val.(type) {
	case []byte:
		err := json.Unmarshal(v, &ti)
		return err
	case string:
		err := json.Unmarshal([]byte(v), &ti)
		return err
	default:
		return errors.New(fmt.Sprintf("Unsupported type: %T", v))
	}
}
func (ti *TitleInfo) Value() driver.Value {
	l, _ := json.Marshal(&ti)
	return l
}

func (title *Title) CreateMe() (err error) {
	_, err = DB.NamedExec("INSERT INTO titles (name, main_group, sub_group, flair, info) VALUES (:name, :main_group, :sub_group, :flair, :info)", map[string]interface{}{
		"name":       title.Name,
		"main_group": title.MainGroup,
		"sub_group":  title.SubGroup,
		"flair":      title.Flair,
		"info":       title.Info.Value(),
	})
	return
}

func (title *Title) GetByName(name string) (err error) {
	err = DB.Get(title, "SELECT * FROM titles WHERE name=$1;", name)
	return
}

func (title *Title) SaveChanges() (err error) {
	return title.ChangeTitleName(title.Name)
}

func (title *Title) ChangeTitleName(oldName string) (err error) {
	_, err = DB.NamedExec("UPDATE titles SET name=:name, main_group=:main_group, sub_group=:sub_group, flair=:flair, info=:info WHERE name=:oldName", map[string]interface{}{
		"oldName":    oldName,
		"name":       title.Name,
		"main_group": title.MainGroup,
		"sub_group":  title.SubGroup,
		"flair":      title.Flair,
		"info":       title.Info.Value(),
	})
	return
}

func (title *Title) DeleteMe() (err error) {
	_, err = DB.NamedExec("DELETE FROM titles WHERE name=:name", map[string]interface{}{
		"name": title.Name,
	})
	return
}

func (titleList *TitleList) GetAll() (err error) {
	err = DB.Select(titleList, "SELECT * FROM titles ORDER BY name;")
	return
}

func (titleList *TitleList) GetAllForDisplayName(displayName string) (err error) {
	err = DB.Select(titleList, "SELECT * FROM titles WHERE info -> 'names' ? $1 ORDER BY name;", displayName)
	return
}
