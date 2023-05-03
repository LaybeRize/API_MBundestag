package database

import (
	"database/sql"
	"gorm.io/gorm"
)

//Sort Functions

func (orgList OrganisationList) Len() int {
	return len(orgList)
}

func (orgList OrganisationList) Less(i, j int) bool {
	return orgList[i].Name < orgList[j].Name
}

func (orgList OrganisationList) Swap(i, j int) {
	orgList[i], orgList[j] = orgList[j], orgList[i]
}

func (orgList SubGroupListOrg) Len() int {
	return len(orgList)
}

func (orgList SubGroupListOrg) Less(i, j int) bool {
	return orgList[i].SubGroup < orgList[j].SubGroup
}

func (orgList SubGroupListOrg) Swap(i, j int) {
	orgList[i], orgList[j] = orgList[j], orgList[i]
}

func (orgList MainGroupListOrg) Len() int {
	return len(orgList)
}

func (orgList MainGroupListOrg) Less(i, j int) bool {
	return orgList[i].MainGroup < orgList[j].MainGroup
}

func (orgList MainGroupListOrg) Swap(i, j int) {
	orgList[i], orgList[j] = orgList[j], orgList[i]
}

//Sort Functions End

type (
	StatusString     string
	MainGroupListOrg []Organisation
	SubGroupListOrg  []Organisation
	OrganisationList []Organisation
	Organisation     struct {
		Name      string `gorm:"primaryKey"`
		MainGroup string
		SubGroup  string
		Flair     sql.NullString
		Status    StatusString
		Members   []Account `gorm:"many2many:organisation_member;foreignKey:name;joinForeignKey:name;References:id;joinReferences:id"`
		Admins    []Account `gorm:"many2many:organisation_admins;foreignKey:name;joinForeignKey:name;References:id;joinReferences:id"`
		Accounts  []Account `gorm:"many2many:organisation_account;foreignKey:name;joinForeignKey:name;References:id;joinReferences:id"`
	}
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

func (org *Organisation) CreateMe() (err error) {
	err = db.Create(org).Error
	return
}

func (org *Organisation) GetByName(name string) (err error) {
	*org = Organisation{}
	err = db.Preload("Members").Preload("Admins").Preload("Accounts").Where("name = ?", name).First(org).Error
	return
}

func (org *Organisation) GetByNameAndOnlyWithAccount(name string, accountID int64) (err error) {
	*org = Organisation{}
	err = db.Joins("JOIN organisation_account ON organisations.name = organisation_account.name").
		Where("organisation_account.id = ?", accountID).Select("*").Table("organisations").Where("organisations.name = ?", name).
		Preload("Members").Preload("Admins").Preload("Accounts").First(org).Error
	return
}

func (org *Organisation) GetByNameAndOnlyWhenAccountIsMember(name string, accountID int64) (err error) {
	*org = Organisation{}
	err = db.Joins("JOIN organisation_member ON organisations.name = organisation_member.name").
		Where("organisation_member.id = ?", accountID).Select("*").Table("organisations").Where("organisations.name = ?", name).
		Preload("Members").Preload("Admins").Preload("Accounts").First(org).Error
	if err == gorm.ErrRecordNotFound {
		err = org.GetByNameAndOnlyWhenAccountAsAdmin(name, accountID)
	}
	return
}

func (org *Organisation) GetByNameAndOnlyWhenAccountAsAdmin(name string, accountID int64) (err error) {
	*org = Organisation{}
	err = db.Joins("JOIN organisation_admins ON organisations.name = organisation_admins.name").
		Where("organisation_admins.id = ?", accountID).Select("*").Table("organisations").Where("organisations.name = ?", name).
		Preload("Members").Preload("Admins").Preload("Accounts").First(org).Error
	return
}

func (org *Organisation) SaveChanges() (err error) {
	err = db.Save(org).Error
	return
}

func (org *Organisation) UpdateMembers() (err error) {
	err = db.Model(org).Association("Members").Replace(org.Members)
	return
}

func (org *Organisation) UpdateAdmins() (err error) {
	err = db.Model(org).Association("Admins").Replace(org.Admins)
	return
}

func (org *Organisation) UpdateAccounts() (err error) {
	err = db.Model(org).Association("Accounts").Replace(org.Accounts)
	return
}

func (orgList *OrganisationList) GetAllVisibleFor(accountID int64) (err error) {
	*orgList = OrganisationList{}
	err = db.Joins("LEFT JOIN organisation_account ON organisations.name = organisation_account.name").
		Where("organisation_account.id = ? OR status = 'public' OR status = 'private'", accountID).Select("organisations.name, main_group, sub_group, flair, status").Table("organisations").Order("organisations.name").
		Preload("Members").Preload("Admins").Preload("Accounts").Find(orgList).Error
	return
}

func (orgList *OrganisationList) GetAllPartOf(accountID int64) (err error) {
	*orgList = OrganisationList{}
	err = db.Joins("LEFT JOIN organisation_member ON organisations.name = organisation_member.name").
		Joins("LEFT JOIN organisation_admins ON organisations.name = organisation_admins.name").
		Where("organisation_member.id = ? OR organisation_admins.id = ?", accountID, accountID).Select("organisations.name, main_group, sub_group, flair, status").Table("organisations").Order("organisations.name").
		Preload("Members").Preload("Admins").Preload("Accounts").Find(orgList).Error
	return
}

func (orgList *OrganisationList) GetAllVisable() (err error) {
	*orgList = OrganisationList{}
	err = db.Where("status != 'hidden'").Order("name").
		Preload("Members").Preload("Admins").Preload("Accounts").Find(orgList).Error
	return
}

func (orgList *OrganisationList) GetAllInvisable() (err error) {
	*orgList = OrganisationList{}
	err = db.Where("status = 'hidden'").Order("name").
		Preload("Members").Preload("Admins").Preload("Accounts").Find(orgList).Error
	return
}

func (orgList *OrganisationList) GetAllSubGroups() (err error) {
	*orgList = OrganisationList{}
	err = db.Distinct("sub_group").Order("sub_group").Find(orgList).Error
	return
}

func (orgList *OrganisationList) GetAllMainGroups() (err error) {
	*orgList = OrganisationList{}
	err = db.Distinct("main_group").Order("main_group").Find(orgList).Error
	return
}

func DeleteMeFromOrganisations(accountID int64) (err error) {
	err = db.Exec("DELETE FROM organisation_admins WHERE id = ?", accountID).Error
	if err != nil {
		return
	}
	err = db.Exec("DELETE FROM organisation_member WHERE id = ?", accountID).Error
	return
}
