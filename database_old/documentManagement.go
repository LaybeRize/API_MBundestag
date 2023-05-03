package database

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/lib/pq"
	"time"
)

var DocumentSchema = `
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'doc_type') THEN
		CREATE TYPE DOC_TYPE AS ENUM ('legislative_text', 'discussion','finished_discussion', 'finished_vote', 'unfinished_vote');
    END IF;
END$$;
CREATE TABLE IF NOT EXISTS documents (
    uuid text UNIQUE NOT NULL,
    written TIMESTAMP NOT NULL,
    organisation TEXT NOT NULL,
	type DOC_TYPE NOT NULL,
    author TEXT NOT NULL,
    flair TEXT NOT NULL,
    title TEXT NOT NULL,
    subtitle TEXT,
    html_content TEXT NOT NULL,
    private BOOLEAN NOT NULL,
    blocked BOOLEAN NOT NULL,
    info jsonb NOT NULL
);
`

func TestDocumentsDB() {
	TestDatabase("DROP TABLE IF EXISTS documents;", "DROP TYPE IF EXISTS DOC_TYPE;")
	InitDocumentsDatabase()
}

type (
	DocumentType     string
	VoteType         string
	DocumentTypeList []DocumentType
	DocumentList     []Document
	Document         struct {
		UUID         string
		Written      time.Time
		Organisation string
		Type         DocumentType
		Author       string
		Flair        string
		Title        string
		Subtitle     sql.NullString
		HTMLContent  string `db:"html_content"`
		Private      bool
		Blocked      bool
		Info         DocumentInfo
	}
	DocumentInfo struct {
		Viewer                    []string      `json:"viewer"`
		Poster                    []string      `json:"poster"`
		Allowed                   []string      `json:"allowed"` //for the main accounts that are allowed to view the post if it is private
		AnyPosterAllowed          bool          `json:"anyPosterAllowed"`
		OrganisationPosterAllowed bool          `json:"organisationPosterAllowed"`
		Finishing                 time.Time     `json:"time"`
		Post                      []Posts       `json:"post"`
		Discussion                []Discussions `json:"discussion"`
		Votes                     []string      `json:"vote"`
	}
	Posts struct {
		UUID      string    `json:"uuid"`
		Hidden    bool      `json:"hidden"`
		Submitted time.Time `json:"submitted"`
		Info      string    `json:"info"`
		Color     string    `json:"color"`
	}
	Discussions struct {
		UUID        string    `json:"uuid"`
		Hidden      bool      `json:"hidden"`
		Written     time.Time `json:"written"`
		Author      string    `json:"author"`
		Flair       string    `json:"flair"`
		HTMLContent string    `json:"htmlContent"`
	}
)

const (
	LegislativeText    DocumentType = "legislative_text"
	Discussion         DocumentType = "discussion"
	FinishedDiscussion DocumentType = "finished_discussion"
	FinishedVote       DocumentType = "finished_vote"
	UnfinishedVote     DocumentType = "unfinished_vote"
)

func InitDocumentsDatabase() {
	DB.MustExec(DocumentSchema)
}

func (docI *DocumentInfo) Scan(val interface{}) error {
	switch v := val.(type) {
	case []byte:
		err := json.Unmarshal(v, &docI)
		return err
	case string:
		err := json.Unmarshal([]byte(v), &docI)
		return err
	default:
		return errors.New(fmt.Sprintf("Unsupported type: %T", v))
	}
}

func (docI *DocumentInfo) Value() driver.Value {
	l, _ := json.Marshal(&docI)
	return l
}

func (docIValue DocumentTypeList) Value() (arr []string) {
	for _, val := range docIValue {
		arr = append(arr, string(val))
	}
	return
}

func (documentation *Document) CreateMe() (err error) {
	_, err = DB.NamedExec("INSERT INTO documents (uuid, written, organisation, type, author, flair, title, subtitle, html_content, private,blocked, info) VALUES (:uuid, :written, :organisation, :type, :author, :flair, :title, :subtitle, :html_content, :private,:blocked, :info)", map[string]interface{}{
		"uuid":         documentation.UUID,
		"written":      time.Now().UTC(),
		"organisation": documentation.Organisation,
		"type":         documentation.Type,
		"author":       documentation.Author,
		"flair":        documentation.Flair,
		"title":        documentation.Title,
		"subtitle":     documentation.Subtitle,
		"html_content": documentation.HTMLContent,
		"private":      documentation.Private,
		"blocked":      false,
		"info":         documentation.Info.Value(),
	})
	return
}

func (documentation *Document) GetByID(uuid string) (err error) {
	err = DB.Get(documentation, "SELECT * FROM documents WHERE uuid=$1;", uuid)
	return
}

func (documentation *Document) SaveChanges() (err error) {
	_, err = DB.NamedExec("UPDATE documents SET private=:private, info=:info, blocked=:blocked WHERE uuid=:uuid", map[string]interface{}{
		"uuid":    documentation.UUID,
		"private": documentation.Private,
		"blocked": documentation.Blocked,
		"info":    documentation.Info.Value(),
	})
	return
}

var GetDocumentsAfterSchema = `
WITH orgView AS (SELECT name FROM organisations WHERE info -> 'viewer' ? $10 ) 
SELECT doc.uuid, doc.written, doc.organisation, doc.type, doc.author, doc.flair, doc.title, doc.subtitle, doc.html_content, doc.private, doc.blocked, doc.info 
FROM documents AS doc LEFT JOIN orgView ON orgView.name = doc.organisation WHERE 
doc.written < $2 AND 
doc.uuid != $1 AND 
(doc.blocked = false OR $14) AND 
(doc.private = false OR doc.info -> 'allowed' ? $10 OR (doc.private = true AND orgView.name IS NOT NULL) OR $11) AND 
($4 OR doc.type = ANY($5)) AND 
($6 OR doc.organisation = $7) AND 
($8 OR doc.author  LIKE ('%' || $9 || '%')) AND 
($12 OR doc.title LIKE ('%' || $13 || '%')) 
ORDER BY doc.written DESC LIMIT $3;
`

func (docList *DocumentList) GetDocumentsAfter(docUUID string, amount int, infos map[string]string, types ...DocumentType) (err error, exists bool) {
	exists = true
	doc := Document{}
	err = doc.GetByID(docUUID)
	if err == sql.ErrNoRows {
		exists = false
		doc.Written = time.Now().UTC()
	} else if err != nil {
		return
	}
	err = DB.Select(docList, GetDocumentsAfterSchema, docUUID, doc.Written, amount,
		len(types) == 0, pq.StringArray(DocumentTypeList(types).Value()),
		infos["organisation"] == "", infos["organisation"],
		infos["author"] == "", infos["author"],
		infos["displayname"], infos["admin"] == "true",
		infos["title"] == "", infos["title"], infos["blocked"] == "true")
	return
}

var GetDocumentsBeforeSchema = `
WITH orgView AS (SELECT name FROM organisations WHERE info -> 'viewer' ? $10 ) 
SELECT X.uuid, X.written, X.organisation, X.type, X.author, X.flair, X.title, X.subtitle, X.html_content, X.private, X.blocked, X.info 
FROM (SELECT * FROM documents AS doc LEFT JOIN orgView ON orgView.name = doc.organisation WHERE 
doc.written > $2 AND 
doc.uuid != $1 AND 
(doc.blocked = false OR $14) AND 
(doc.private = false OR doc.info -> 'allowed' ? $10 OR (doc.private = true AND orgView.name IS NOT NULL) OR $11) AND 
($4 OR doc.type = ANY($5)) AND 
($6 OR doc.organisation = $7) AND 
($8 OR doc.author  LIKE ('%' || $9 || '%')) AND 
($12 OR doc.title LIKE ('%' || $13 || '%')) 
ORDER BY doc.written LIMIT $3) as X ORDER BY X.written DESC;
`

func (docList *DocumentList) GetGetDocumentsBefore(docUUID string, amount int, infos map[string]string, types ...DocumentType) (err error, exists bool) {
	exists = true
	doc := Document{}
	err = doc.GetByID(docUUID)
	if err == sql.ErrNoRows {
		exists = false
		doc.Written = time.Now().UTC()
	} else if err != nil {
		return
	}
	err = DB.Select(docList, GetDocumentsBeforeSchema, docUUID, doc.Written, amount,
		len(types) == 0, pq.StringArray(DocumentTypeList(types).Value()),
		infos["organisation"] == "", infos["organisation"],
		infos["author"] == "", infos["author"],
		infos["displayname"], infos["admin"] == "true",
		infos["title"] == "", infos["title"], infos["blocked"] == "true")
	return
}

func (docList *DocumentList) GetAllOpenDiscussionsAndVotes() (err error) {
	err = DB.Select(docList, "SELECT * FROM documents WHERE type = $1 OR type = $2", Discussion, UnfinishedVote)
	return
}

func (documentation *Document) GetByIDForAccount(uuid string, displayName string) (err error) {
	err = DB.Get(documentation, `WITH orgView AS (SELECT name FROM organisations WHERE info -> 'viewer' ? $2 ) 
SELECT doc.uuid, doc.written, doc.organisation, doc.type, doc.author, doc.flair, doc.title, doc.subtitle, doc.html_content, doc.private, doc.blocked, doc.info 
FROM documents AS doc LEFT JOIN orgView ON orgView.name = doc.organisation WHERE doc.uuid=$1 AND 
(doc.private = false OR doc.info -> 'allowed' ? $2 OR (doc.private = true AND orgView.name IS NOT NULL));`, uuid, displayName)
	return
}
