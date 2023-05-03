package database

import "database/sql"

//Sort Functions

func (titleList TitleList) Len() int {
	return len(titleList)
}

func (titleList TitleList) Less(i, j int) bool {
	return titleList[i].Name < titleList[j].Name
}

func (titleList TitleList) Swap(i, j int) {
	titleList[i], titleList[j] = titleList[j], titleList[i]
}

func (titleList SubGroupListTitle) Len() int {
	return len(titleList)
}

func (titleList SubGroupListTitle) Less(i, j int) bool {
	return titleList[i].SubGroup < titleList[j].SubGroup
}

func (titleList SubGroupListTitle) Swap(i, j int) {
	titleList[i], titleList[j] = titleList[j], titleList[i]
}

func (titleList MainGroupListTitle) Len() int {
	return len(titleList)
}

func (titleList MainGroupListTitle) Less(i, j int) bool {
	return titleList[i].MainGroup < titleList[j].MainGroup
}

func (titleList MainGroupListTitle) Swap(i, j int) {
	titleList[i], titleList[j] = titleList[j], titleList[i]
}

// Structs and functions
type (
	MainGroupListTitle []Title
	SubGroupListTitle  []Title
	TitleList          []Title
	Title              struct {
		Name      string `gorm:"primaryKey"`
		MainGroup string
		SubGroup  string
		Flair     sql.NullString
		Holder    []Account `gorm:"many2many:title_account;foreignKey:name;joinForeignKey:name;References:id;joinReferences:id"`
	}
)

func (title *Title) CreateMe() (err error) {
	err = db.Create(title).Error
	return
}

func (title *Title) GetByName(name string) (err error) {
	*title = Title{}
	err = db.Preload("Holder").Where("name = ?", name).First(title).Error
	return
}

func (title *Title) SaveChanges() (err error) {
	err = db.Save(title).Error
	return
}

func (title *Title) ChangeTitleName(oldName string) (err error) {
	err = db.Model(Title{}).Where("name = ?", oldName).Updates(title).Error
	return
}

func (title *Title) UpdateHolder() (err error) {
	err = db.Model(title).Association("Holder").Replace(title.Holder)
	return
}

func (title *Title) DeleteMe() (err error) {
	err = db.Delete(title).Error
	return
}

func (titleList *TitleList) GetAll() (err error) {
	err = db.Preload("Holder").Order("name").Find(titleList).Error
	return
}

func (titleList *TitleList) GetAllForUserID(userID int64) (err error) {
	err = db.Joins("JOIN title_account ON titles.name = title_account.name").
		Where("title_account.id = ?", userID).Select("*").Table("titles").Order("titles.name").Find(titleList).Error
	return
}

func DeleteMeFromTitles(accountID int64) (err error) {
	err = db.Exec("DELETE FROM title_account WHERE id = ?", accountID).Error
	return
}
