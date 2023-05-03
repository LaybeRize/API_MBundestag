package database

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
)

type (
	StatusString string
)

const (
	Public  StatusString = "public"
	Private StatusString = "private"
	Secret  StatusString = "secret"
	Hidden  StatusString = "hidden"
)

var StatusTranslation = map[StatusString]string{
	Public:  "Öffentlich",
	Private: "Privat",
	Secret:  "Geheim",
	Hidden:  "Versteckt",
}
var Stati = []string{string(Public), string(Private), string(Secret), string(Hidden)}

var OrganisationsSchema = `
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'status') THEN
		CREATE TYPE STATUS AS ENUM ('public', 'private', 'secret', 'hidden');
    END IF;
END$$;
CREATE TABLE IF NOT EXISTS organisations (
    name text UNIQUE NOT NULL,
    main_group TEXT NOT NULL,
    sub_group TEXT NOT NULL,
    flair TEXT UNIQUE,
    status STATUS NOT NULL,
    info jsonb NOT NULL
);
`

func TestOrganisationsDB() {
	TestDatabase("DROP TABLE IF EXISTS organisations;", "DROP TYPE IF EXISTS STATUS;")
	InitOrganisationsDatabase()
}

type (
	MainGroupListOrg []Organisation
	SubGroupListOrg  []Organisation
	OrganisationList []Organisation
	Organisation     struct {
		Name      string
		MainGroup string `db:"main_group"`
		SubGroup  string `db:"sub_group"`
		Flair     sql.NullString
		Status    StatusString
		Info      OrganisationInfo
	}
	OrganisationInfo struct {
		Admins []string `json:"admins"`
		User   []string `json:"user"`
		Viewer []string `json:"viewer"`
	}
)

//Sort Functions

func (orgList *OrganisationList) Len() int {
	return len(*orgList)
}

func (orgList *OrganisationList) Less(i, j int) bool {
	strI := (*orgList)[i].Name
	strJ := (*orgList)[j].Name
	return strI < strJ
}

func (orgList *OrganisationList) Swap(i, j int) {
	(*orgList)[i], (*orgList)[j] = (*orgList)[j], (*orgList)[i]
}

func (orgList SubGroupListOrg) Len() int {
	return len(orgList)
}

func (orgList SubGroupListOrg) Less(i, j int) bool {
	strI := orgList[i].SubGroup
	strJ := orgList[j].SubGroup
	return strI < strJ
}

func (orgList SubGroupListOrg) Swap(i, j int) {
	orgList[i], orgList[j] = orgList[j], orgList[i]
}

func (orgList MainGroupListOrg) Len() int {
	return len(orgList)
}

func (orgList MainGroupListOrg) Less(i, j int) bool {
	strI := orgList[i].MainGroup
	strJ := orgList[j].MainGroup
	return strI < strJ
}

func (orgList MainGroupListOrg) Swap(i, j int) {
	orgList[i], orgList[j] = orgList[j], orgList[i]
}

//Sort Functions End

func InitOrganisationsDatabase() {
	DB.MustExec(OrganisationsSchema)
}

func (oi *OrganisationInfo) Scan(val interface{}) error {
	switch v := val.(type) {
	case []byte:
		err := json.Unmarshal(v, &oi)
		return err
	case string:
		err := json.Unmarshal([]byte(v), &oi)
		return err
	default:
		return errors.New(fmt.Sprintf("Unsupported type: %T", v))
	}
}
func (oi *OrganisationInfo) Value() driver.Value {
	l, _ := json.Marshal(&oi)
	return l
}

func (org *Organisation) CreateMe() (err error) {
	_, err = DB.NamedExec("INSERT INTO organisations (name, main_group, sub_group, flair, status, info) VALUES (:name, :main_group, :sub_group, :flair, :status, :info)", map[string]interface{}{
		"name":       org.Name,
		"main_group": org.MainGroup,
		"sub_group":  org.SubGroup,
		"flair":      org.Flair,
		"status":     org.Status,
		"info":       org.Info.Value(),
	})
	return
}

func (org *Organisation) GetByName(name string) (err error) {
	err = DB.Get(org, "SELECT * FROM organisations WHERE name=$1;", name)
	return
}

func (org *Organisation) GetByNameAndOnlyWithAccount(name string, displayName string) (err error) {
	err = DB.Get(org, "SELECT * FROM organisations WHERE name=$1 AND info -> 'viewer' ? $2;", name, displayName)
	return
}

func (org *Organisation) SaveChanges() (err error) {
	_, err = DB.NamedExec("UPDATE organisations SET main_group=:main_group, sub_group=:sub_group, flair=:flair, status=:status, info=:info WHERE name=:name", map[string]interface{}{
		"name":       org.Name,
		"main_group": org.MainGroup,
		"sub_group":  org.SubGroup,
		"flair":      org.Flair,
		"status":     org.Status,
		"info":       org.Info.Value(),
	})
	return
}

func (orgList *OrganisationList) GetAllVisibleFor(displayName string) (err error) {
	err = DB.Select(orgList, "SELECT * FROM organisations WHERE status = 'public' OR status = 'private' OR info -> 'viewer' ? $1 ORDER BY name;", displayName)
	return
}

func (orgList *OrganisationList) GetAllPartOf(displayName string) (err error) {
	err = DB.Select(orgList, "SELECT * FROM organisations WHERE info -> 'user' ? $1 OR info -> 'admins' ? $1 ORDER BY name;", displayName)
	return
}

func (orgList *OrganisationList) GetAllVisable() (err error) {
	err = DB.Select(orgList, "SELECT * FROM organisations WHERE status != 'hidden' ORDER BY name;")
	return
}

func (orgList *OrganisationList) GetAllInvisable() (err error) {
	err = DB.Select(orgList, "SELECT * FROM organisations WHERE status = 'hidden' ORDER BY name;")
	return
}

func (orgList *OrganisationList) GetAllSubGroups() (err error) {
	err = DB.Select(orgList, "SELECT DISTINCT sub_group FROM organisations ORDER BY sub_group;")
	return
}

func (orgList *OrganisationList) GetAllMainGroups() (err error) {
	err = DB.Select(orgList, "SELECT DISTINCT main_group FROM organisations ORDER BY main_group;")
	return
}
