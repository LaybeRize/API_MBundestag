package database

import (
	"database/sql"
	"time"
)

type (
	ZwitscherList []Zwitscher
	Zwitscher     struct {
		UUID           string `gorm:"primaryKey"`
		Written        time.Time
		Blocked        bool
		Author         string
		Flair          string
		HTMLContent    string
		ConnectedTo    sql.NullString
		AmountChildren int64
		Parent         *Zwitscher  `gorm:"foreignKey:connected_to;joinReferences:uuid"`
		Children       []Zwitscher `gorm:"foreignKey:connected_to"`
	}
)

func (zwitscher *Zwitscher) CreateMe() (err error) {
	err = db.Create(zwitscher).Error
	return
}

func (zwitscher *Zwitscher) GetByUUID(uuid string) (err error) {
	*zwitscher = Zwitscher{}
	err = db.Preload("Children").Preload("Parent").First(zwitscher, "uuid = ?", uuid).Error
	return
}

func (zwitscher *Zwitscher) SaveChanges() (err error) {
	err = db.Save(zwitscher).Error
	return
}

func (zwitscherList *ZwitscherList) GetLatested(amount int, allowBlocked bool) (err error) {
	err = db.Where("connected_to IS NULL AND (blocked = false OR ?)", allowBlocked).Order("written desc").Limit(amount).Find(zwitscherList).Error
	return
}
