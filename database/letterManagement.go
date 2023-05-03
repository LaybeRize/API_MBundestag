package database

import (
	"gorm.io/gorm"
	"time"
)

type (
	LetterList []Letter
	Letter     struct {
		UUID        string `gorm:"primaryKey"`
		Written     time.Time
		Author      string
		Flair       string
		Title       string
		Content     string
		HTMLContent string     `gorm:"column:html_content"`
		Info        LetterInfo `gorm:"type:jsonb;serializer:json"`
		Viewer      []Account  `gorm:"many2many:letter_account;foreignKey:uuid;joinForeignKey:uuid;References:id;joinReferences:id"`
		Removed     bool
		ModMessage  bool `gorm:"column:mod_message"`
	}
	LetterInfo struct {
		AllHaveToAgree     bool     `json:"allAgree"`
		NoSigning          bool     `json:"noSigning"`
		PeopleNotYetSigned []string `json:"notSigned"`
		Signed             []string `json:"signed"`
		Rejected           []string `json:"rejected"`
	}
)

func (letter *Letter) CreateMe() (err error) {
	letter.Written = time.Now().UTC()
	err = db.Create(letter).Error
	return
}

func (letter *Letter) GetByID(uuid string) (err error) {
	*letter = Letter{}
	err = db.Where("uuid = ?", uuid).First(letter).Error
	return
}

func (letter *Letter) GetByIDWithViewer(uuid string) (err error) {
	*letter = Letter{}
	err = db.Preload("Viewer").Where("uuid = ?", uuid).First(letter).Error
	return
}

func (letter *Letter) SaveChanges() (err error) {
	err = db.Save(letter).Error
	return
}

func (letterList *LetterList) GetLettersAfter(publicationUUID string, amount int, accountId int64) (err error, exists bool) {
	return letterList.getLetters(publicationUUID, func(pub *Letter) *gorm.DB {
		return getBasicLetterQuery(pub.UUID, amount, accountId).Where("written < ?", pub.Written).Order("written desc")
	})
}

func (letterList *LetterList) GetLettersBefore(publicationUUID string, amount int, accountId int64) (err error, exists bool) {
	return letterList.getLetters(publicationUUID, func(pub *Letter) *gorm.DB {
		return db.Select("*").Table("(?) as X", getBasicLetterQuery(pub.UUID, amount, accountId).Where("written > ?", pub.Written).Order("written")).Order("X.written desc")
	})
}

func (letterList *LetterList) GetModMailsAfter(publicationUUID string, amount int) (err error, exists bool) {
	return letterList.getLetters(publicationUUID, func(pub *Letter) *gorm.DB {
		return getBasicModmailQuery(pub.UUID, amount).Where("written < ?", pub.Written).Order("written desc")
	})
}

func (letterList *LetterList) GetModMailsBefore(publicationUUID string, amount int) (err error, exists bool) {
	return letterList.getLetters(publicationUUID, func(pub *Letter) *gorm.DB {
		return db.Select("*").Table("(?) as X", getBasicModmailQuery(pub.UUID, amount).Where("written > ?", pub.Written).Order("written")).Order("X.written desc")
	})
}

func (letterList *LetterList) getLetters(publicationUUID string, query func(pub *Letter) *gorm.DB) (err error, exists bool) {
	*letterList = LetterList{}
	exists = true
	pub := Letter{}
	err = pub.GetByID(publicationUUID)
	if err == gorm.ErrRecordNotFound {
		exists = false
		pub.Written = time.Now().UTC()
	} else if err != nil {
		return
	}
	err = query(&pub).Find(letterList).Error
	return
}

func getBasicLetterQuery(uuid string, amount int, accountID int64) *gorm.DB {
	return db.Joins("JOIN letter_account ON letters.uuid = letter_account.uuid").
		Where("letter_account.id = ?", accountID).Select("*").Where("letters.uuid != ?", uuid).Limit(amount).Table("letters")
}

func getBasicModmailQuery(uuid string, amount int) *gorm.DB {
	return db.Where("letters.uuid != ? AND mod_message = true", uuid).Limit(amount).Table("letters")
}
