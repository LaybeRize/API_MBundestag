package database

import (
	"database/sql"
	"time"
)

var ZwitscherSchema = `
CREATE TABLE IF NOT EXISTS zwitscher (
    uuid text UNIQUE NOT NULL,
    written TIMESTAMP NOT NULL,
    blocked BOOLEAN NOT NULL,
    author TEXT NOT NULL,
    flair TEXT NOT NULL,
    html_content TEXT NOT NULL,
    connectedTo TEXT
);
`

func TestZwitscherDB() {
	TestDatabase("DROP TABLE IF EXISTS zwitscher;", "")
	InitZwitscherDatabase()
}

type (
	ZwitscherList []Zwitscher
	Zwitscher     struct {
		UUID        string
		Written     time.Time
		Blocked     bool
		Author      string
		Flair       string
		HTMLContent string `db:"html_content"`
		ConnectedTo sql.NullString
	}
)

func InitZwitscherDatabase() {
	DB.MustExec(ZwitscherSchema)
}

func (zwitscher *Zwitscher) CreateMe() (err error) {
	_, err = DB.NamedExec("INSERT INTO zwitscher (uuid, written, blocked, author, flair, html_content, connectedTo) VALUES (:uuid, :written, :blocked, :author, :flair, :html_content, :connectedTo)", map[string]interface{}{
		"uuid":         zwitscher.UUID,
		"written":      time.Now().UTC(),
		"blocked":      false,
		"author":       zwitscher.Author,
		"flair":        zwitscher.Flair,
		"html_content": zwitscher.HTMLContent,
		"connectedTo":  zwitscher.ConnectedTo,
	})
	return
}

func (zwitscher *Zwitscher) GetByID(uuid string) (err error) {
	err = DB.Get(zwitscher, "SELECT * FROM zwitscher WHERE uuid=$1;", uuid)
	return
}

func (zwitscher *Zwitscher) SaveChanges() (err error) {
	_, err = DB.NamedExec("UPDATE zwitscher SET blocked=:blocked WHERE uuid=:uuid", map[string]interface{}{
		"uuid":    zwitscher.UUID,
		"blocked": zwitscher.Blocked,
	})
	return
}

func (zwitscherList *ZwitscherList) GetCommentsFor(uuid string, allowBlocked bool) (err error) {
	err = DB.Select(zwitscherList, "SELECT * FROM zwitscher WHERE connectedTo = $1 AND ($2 OR blocked = false) ORDER BY written DESC;", uuid, allowBlocked)
	return
}

func (zwitscherList *ZwitscherList) GetLatested(amount int, allowBlocked bool) (err error) {
	err = DB.Select(zwitscherList, "SELECT * FROM zwitscher WHERE ($1 OR blocked = false) AND connectedTo IS NULL ORDER BY written DESC LIMIT $2;", allowBlocked, amount)
	return
}

func (zwitscherList *ZwitscherList) GetByUser(displayName string, amount int, allowBlocked bool) (err error) {
	err = DB.Select(zwitscherList, "SELECT * FROM zwitscher WHERE author = $1 AND ($2 OR blocked = false) ORDER BY written DESC LIMIT $3;", displayName, allowBlocked, amount)
	return
}
