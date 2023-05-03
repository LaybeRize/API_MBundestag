package database

import (
	"API_MBundestag/help"
	"database/sql"
	"errors"
	"github.com/lib/pq"
)

type (
	RoleString  string
	AccountList []Account
	Names       struct {
		DisplayName string `gorm:"unique;column:name"`
		Username    string `gorm:"unique"`
	}
	NameList []Names
)

type Account struct {
	ID            int64  `gorm:"primaryKey;autoIncrement:true"`
	DisplayName   string `gorm:"index:unique;column:name"`
	Flair         string
	Username      string `gorm:"index:unique"`
	Password      string
	Suspended     bool
	RefToken      sql.NullString `gorm:"column:refresh_token"`
	ExpDate       sql.NullTime   `gorm:"column:expiration_date"`
	LoginTries    int            `gorm:"column:login_tries"`
	NextLoginTime sql.NullTime   `gorm:"column:next_login_allowed"`
	Role          RoleString
	Linked        sql.NullInt64
	Parent        *Account  `gorm:"foreignKey:linked;joinReferences:id"`
	Children      []Account `gorm:"foreignKey:linked"`
}

const (
	User         RoleString = "user"
	MediaAdmin   RoleString = "media_admin"
	Admin        RoleString = "admin"
	HeadAdmin    RoleString = "head_admin"
	PressAccount RoleString = "press_account"
	NotLoggedIn  RoleString = "notLoggedIn"
)

var Roles = []string{string(PressAccount), string(User), string(MediaAdmin), string(Admin), string(HeadAdmin)}
var RoleTranslation = map[RoleString]string{
	PressAccount: "Presse-Account",
	User:         "Nutzer",
	MediaAdmin:   "Medien-Administrator",
	Admin:        "Administrator",
	HeadAdmin:    "Leitender Administrator",
}

func (user *Account) CreateMe() (err error) {
	err = db.Create(user).Error
	return
}

func (user *Account) GetByUserName(username string) (err error) {
	*user = Account{}
	err = db.First(user, "username=?", username).Error
	return
}

func (user *Account) GetByDisplayName(name string) (err error) {
	*user = Account{}
	err = db.First(user, "name=?", name).Error
	return
}

func (user *Account) GetByDisplayNameWithParent(name string) (err error) {
	*user = Account{}
	err = db.Preload("Parent").First(user, "name=? AND suspended = false", name).Error
	return
}

func (user *Account) GetByToken(token string) (err error) {
	*user = Account{}
	err = db.First(user, "refresh_token=?", token).Error
	return
}

func (user *Account) GetByID(id int64) (err error) {
	*user = Account{}
	err = db.First(user, "id=?", id).Error
	return
}

func (user *Account) GetByIDWithChildren(id int64) (err error) {
	*user = Account{}
	err = db.Preload("Children", "suspended = false").First(user, "id = ?", id).Error
	return
}

func (user *Account) SaveChanges() (err error) {
	err = db.Save(user).Error
	return
}

func (accountList *AccountList) GetAllPressAccountsFromAccountPlusSelf(acc *Account) (err error) {
	special := Account{}
	err = special.GetByIDWithChildren(acc.ID)
	*accountList = make(AccountList, len(special.Children)+1)
	(*accountList)[0] = *acc
	copy((*accountList)[1:], special.Children)
	return
}

func (accountList *AccountList) GetAllAccounts() (err error) {
	err = db.Order("name").Find(accountList).Error
	return
}

func (accountList *AccountList) GetAllAccountsNotSuspended() (err error) {
	err = db.Where("suspended = false").Order("name").Find(accountList).Error
	return
}

func (accountList *AccountList) DoAccountsExist(displayNames []string) (b bool, err error) {
	*accountList = AccountList{}

	err = db.Where("name = ANY($1) AND suspended = false", pq.StringArray(displayNames)).Order("name").Find(&accountList).Error
	if len(displayNames) == len(*accountList) {
		return true, err
	}

	for _, item := range *accountList {
		displayNames = help.RemoveFirstStringOccurrenceFromArray(displayNames, item.DisplayName)
	}

	*accountList = AccountList{}
	return false, errors.New(displayNames[0])
}

func (list *NameList) GetAllUserAndDisplayName() (err error) {
	err = db.Model(&Account{}).Order("name").Find(list).Error
	return
}
