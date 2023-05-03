package database

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/lib/pq"
	"gorm.io/gorm"
	"time"
)

type (
	DocumentType     string
	VoteType         string
	DocumentTypeList []DocumentType
	DocumentList     []Document
	Document         struct {
		UUID         string `gorm:"primaryKey"`
		Written      time.Time
		Organisation string
		Type         DocumentType
		Author       string
		Flair        string
		Title        string
		Subtitle     sql.NullString
		HTMLContent  string
		Private      bool
		Blocked      bool
		Info         DocumentInfo `gorm:"type:jsonb"`
		Viewer       []Account    `gorm:"many2many:doc_viewer;foreignKey:uuid;joinForeignKey:uuid;References:id;joinReferences:id"`
		Poster       []Account    `gorm:"many2many:doc_poster;foreignKey:uuid;joinForeignKey:uuid;References:id;joinReferences:id"`
		Allowed      []Account    `gorm:"many2many:doc_allowed;foreignKey:uuid;joinForeignKey:uuid;References:id;joinReferences:id"`
	}
	DocumentInfo struct {
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

const (
	LegislativeText    DocumentType = "legislative_text"
	RunningDiscussion  DocumentType = "running_discussion"
	FinishedDiscussion DocumentType = "finished_discussion"
	RunningVote        DocumentType = "running_vote"
	FinishedVote       DocumentType = "finished_vote"
)

func (documentation *Document) CreateMe() (err error) {
	err = db.Create(documentation).Error
	return
}

func (documentation *Document) GetByID(uuid string) (err error) {
	*documentation = Document{}
	err = db.Preload("Viewer").Preload("Poster").Preload("Allowed").Where("uuid = ?", uuid).First(documentation).Error
	return
}

func (documentation *Document) GetByIDOnlyWithAccount(uuid string, accountID int64) (err error) {
	*documentation = Document{}
	err = db.Joins("LEFT JOIN organisation_account ON documents.organisation = organisation_account.name").
		Joins("LEFT JOIN doc_allowed ON documents.uuid = doc_allowed.uuid").
		Select("documents.uuid, written, organisation, type, author, flair, title, subtitle, html_content, private, blocked, info").Table("documents").
		Preload("Viewer").Preload("Poster").Preload("Allowed").
		Where("documents.uuid = ? AND blocked = false", uuid).Where("private = false OR organisation_account.id = ? OR doc_allowed.id = ?", accountID, accountID).First(documentation).Error
	return
}

func (documentation *Document) SaveChanges() (err error) {
	err = db.Save(documentation).Error
	return
}

func (docList *DocumentList) GetDocumentsAfter(publicationUUID string, amount int, accountID int64, infos map[string]string, types ...DocumentType) (err error, exists bool) {
	return docList.getDocuments(publicationUUID, func(doc *Document) *gorm.DB {
		return getBasicDocumentQuery(doc.UUID, amount, accountID, infos, types).Where("written < ?", doc.Written).Order("written desc")
	})
}

func (docList *DocumentList) GetDocumentsBefore(publicationUUID string, amount int, accountID int64, infos map[string]string, types ...DocumentType) (err error, exists bool) {
	return docList.getDocuments(publicationUUID, func(doc *Document) *gorm.DB {
		return db.Select("*").Table("(?) as X", getBasicDocumentQuery(doc.UUID, amount, accountID, infos, types).
			Where("written > ?", doc.Written).Order("written")).Order("X.written desc")
	})
}

func (docList *DocumentList) GetAdminDocumentsAfter(publicationUUID string, amount int, infos map[string]string, types ...DocumentType) (err error, exists bool) {
	return docList.getDocuments(publicationUUID, func(doc *Document) *gorm.DB {
		return getAdminDocumentQuery(doc.UUID, amount, infos, types).Where("written < ?", doc.Written).Order("written desc")
	})
}

func (docList *DocumentList) GetAdminDocumentsBefore(publicationUUID string, amount int, infos map[string]string, types ...DocumentType) (err error, exists bool) {
	return docList.getDocuments(publicationUUID, func(doc *Document) *gorm.DB {
		return db.Select("*").Table("(?) as X", getAdminDocumentQuery(doc.UUID, amount, infos, types).
			Where("written > ?", doc.Written).Order("written")).Order("X.written desc")
	})
}

func (docList *DocumentList) getDocuments(publicationUUID string, query func(pub *Document) *gorm.DB) (err error, exists bool) {
	*docList = DocumentList{}
	exists = true
	doc := Document{}
	err = doc.GetByID(publicationUUID)
	if err == gorm.ErrRecordNotFound {
		exists = false
		doc.Written = time.Now().UTC()
	} else if err != nil {
		return
	}
	err = query(&doc).Find(docList).Error
	return
}

func getBasicDocumentQuery(uuid string, amount int, accountID int64, infos map[string]string, types []DocumentType) *gorm.DB {
	return generalMapQuery(getBlockDocumentQuery(uuid).
		Joins("LEFT JOIN organisation_account ON documents.organisation = organisation_account.name").
		Joins("LEFT JOIN doc_allowed ON documents.uuid = doc_allowed.uuid").
		Where("private = false OR organisation_account.id = ? OR doc_allowed.id = ?", accountID, accountID),
		infos, types).
		Where("blocked = false").Limit(amount)
}

func getAdminDocumentQuery(uuid string, amount int, infos map[string]string, types []DocumentType) *gorm.DB {
	return generalMapQuery(getBlockDocumentQuery(uuid),
		infos, types).
		Where("blocked = ? OR ?", infos["blocked"] == "true", infos["blocked"] == "all").Limit(amount)
}

func getBlockDocumentQuery(uuid string) *gorm.DB {
	return db.Select("documents.uuid, written, organisation, type, author, flair, title, subtitle, private, blocked").Table("documents").
		Where("documents.uuid != ?", uuid)
}

func generalMapQuery(qry *gorm.DB, infos map[string]string, types []DocumentType) *gorm.DB {
	return qry.Where("type = ANY(?) OR ?", pq.StringArray(DocumentTypeList(types).Value()), len(types) == 0).
		Where("organisation = ? OR ?", infos["organisation"], infos["organisation"] == "").
		Where("author LIKE ('%' || ? || '%') OR ?", infos["author"], infos["author"] == "").
		Where("title LIKE ('%' || ? || '%') OR ?", infos["title"], infos["title"] == "")
}
