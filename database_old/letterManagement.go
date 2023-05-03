package database

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

var LettersSchema = `
CREATE TABLE IF NOT EXISTS letters (
    uuid text UNIQUE NOT NULL,
    written TIMESTAMP NOT NULL,
    author TEXT NOT NULL,
    flair TEXT NOT NULL,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    html_content TEXT NOT NULL,
    info jsonb NOT NULL
);
`

func TestLettersDB() {
	TestDatabase("DROP TABLE IF EXISTS letters;", "")
	InitLettersDatabase()
}

type (
	LetterList []Letter
	Letter     struct {
		UUID        string
		Written     time.Time
		Author      string
		Flair       string
		Title       string
		Content     string
		HTMLContent string `db:"html_content"`
		Info        LetterInfo
	}
	LetterInfo struct {
		ModMessage          bool     `json:"modMessage"`
		AllHaveToAgree      bool     `json:"allAgree"`
		NoSigning           bool     `json:"noSigning"`
		PeopleInvitedToSign []string `json:"viewer"`
		PeopleNotYetSigned  []string `json:"notSigned"`
		Signed              []string `json:"signed"`
		Rejected            []string `json:"rejected"`
	}
)

func InitLettersDatabase() {
	DB.MustExec(LettersSchema)
}

func (li *LetterInfo) Scan(val interface{}) error {
	switch v := val.(type) {
	case []byte:
		err := json.Unmarshal(v, &li)
		return err
	case string:
		err := json.Unmarshal([]byte(v), &li)
		return err
	default:
		return errors.New(fmt.Sprintf("Unsupported type: %T", v))
	}
}
func (li *LetterInfo) Value() driver.Value {
	l, _ := json.Marshal(&li)
	return l
}

func (letter *Letter) CreateMe() (err error) {
	_, err = DB.NamedExec("INSERT INTO letters (uuid, written, author, flair, title, content, html_content, info) VALUES (:uuid, :written, :author, :flair, :title,:content, :html_content, :info)", map[string]interface{}{
		"uuid":         letter.UUID,
		"written":      time.Now().UTC(),
		"author":       letter.Author,
		"flair":        letter.Flair,
		"title":        letter.Title,
		"html_content": letter.HTMLContent,
		"content":      letter.Content,
		"info":         letter.Info.Value(),
	})
	return
}

func (letter *Letter) GetByID(uuid string) (err error) {
	err = DB.Get(letter, "SELECT * FROM letters WHERE uuid=$1;", uuid)
	return
}

func (letter *Letter) SaveChanges() (err error) {
	_, err = DB.NamedExec("UPDATE letters SET info=:info WHERE uuid=:uuid", map[string]interface{}{
		"uuid": letter.UUID,
		"info": letter.Info.Value(),
	})
	return
}

func (letterList *LetterList) GetPublicationAfter(publicationUUID string, amount int, accountName string, modMail bool) (err error, exists bool) {
	exists = true
	pub := Letter{}
	err = pub.GetByID(publicationUUID)
	if err == sql.ErrNoRows {
		exists = false
		pub.Written = time.Now().UTC()
	} else if err != nil {
		return
	}
	err = DB.Select(letterList, "SELECT * FROM letters WHERE written < $2 AND uuid != $1 AND (info -> 'viewer' ? $5 OR ( $4 AND (info ->> 'modMessage')::bool = $4)) ORDER BY written DESC LIMIT $3;", publicationUUID, pub.Written, amount, modMail, accountName)
	return
}

func (letterList *LetterList) GetPublicationBefore(letterUUID string, amount int, accountName string, modMail bool) (err error, exists bool) {
	exists = true
	pub := Letter{}
	err = pub.GetByID(letterUUID)
	if err == sql.ErrNoRows {
		exists = false
		pub.Written = time.Now().UTC()
	} else if err != nil {
		return
	}
	err = DB.Select(letterList, "SELECT * FROM (SELECT * FROM letters WHERE written > $2 AND uuid != $1 AND (info -> 'viewer' ? $5 OR ( $4 AND (info ->> 'modMessage')::bool = $4)) ORDER BY written LIMIT $3) as X ORDER BY X.written DESC;", letterUUID, pub.Written, amount, modMail, accountName)
	return
}
