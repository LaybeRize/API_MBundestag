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
		DisplayName string `db:"name"`
		Username    string
	}
	NameList []Names
)

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

var AccountSchema = `
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'level') THEN
		CREATE TYPE LEVEL AS ENUM ('admin', 'user', 'head_admin', 'media_admin', 'press_account');
    END IF;
END$$;
CREATE TABLE IF NOT EXISTS account (
    id BIGSERIAL,
    name TEXT UNIQUE NOT NULL,
    flair TEXT NOT NULL,
    username TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL,
    suspended BOOLEAN NOT NULL,
    refresh_token TEXT,
    expiration_date TIMESTAMP,
    login_tries INT NOT NULL,
    next_login_allowed TIMESTAMP,
    role LEVEL,
    linked BIGINT
);
`

func TestAccountDB() {
	TestDatabase("DROP TABLE IF EXISTS account;", "DROP TYPE IF EXISTS LEVEL;")
	InitAccountDatabase()
}

type Account struct {
	ID            int64
	DisplayName   string `db:"name"`
	Flair         string
	Username      string
	Password      string
	Suspended     bool
	RefToken      sql.NullString `db:"refresh_token"`
	ExpDate       sql.NullTime   `db:"expiration_date"`
	LoginTries    int            `db:"login_tries"`
	NextLoginTime sql.NullTime   `db:"next_login_allowed"`
	Role          RoleString
	Linked        sql.NullInt64
}

func InitAccountDatabase() {
	DB.MustExec(AccountSchema)
}

func (user *Account) CreateMe() (err error) {
	_, err = DB.NamedExec("INSERT INTO account (username, password, name, role, suspended, login_tries, flair, linked) VALUES (:username, :password, :name, :role, false, 0,:flair, :linked)", map[string]interface{}{
		"username": user.Username,
		"password": user.Password,
		"name":     user.DisplayName,
		"role":     user.Role,
		"flair":    user.Flair,
		"linked":   user.Linked,
	})
	return
}

func (user *Account) GetByUserName(username string) (err error) {
	err = DB.Get(user, "SELECT * FROM account WHERE username=$1;", username)
	return
}

func (user *Account) GetByDisplayName(name string) (err error) {
	err = DB.Get(user, "SELECT * FROM account WHERE name=$1;", name)
	return
}

func (user *Account) GetByToken(token string) (err error) {
	err = DB.Get(user, "SELECT * FROM account WHERE refresh_token=$1;", token)
	return
}

func (user *Account) GetByID(id int64) (err error) {
	err = DB.Get(user, "SELECT * FROM account WHERE id=$1;", id)
	return
}

func (user *Account) SaveChanges() (err error) {
	_, err = DB.NamedExec("UPDATE account SET flair=:flair, password=:password, suspended=:suspended, refresh_token=:refresh_token, expiration_date=:expiration_date, login_tries=:login_tries, next_login_allowed=:next_login_allowed, role=:role, linked=:linked WHERE name=:name", map[string]interface{}{
		"name":               user.DisplayName,
		"flair":              user.Flair,
		"password":           user.Password,
		"suspended":          user.Suspended,
		"refresh_token":      user.RefToken,
		"expiration_date":    user.ExpDate,
		"login_tries":        user.LoginTries,
		"next_login_allowed": user.NextLoginTime,
		"role":               user.Role,
		"linked":             user.Linked,
	})
	return
}

func (accountList *AccountList) GetAllPressAccountsFromAccountPlusSelf(acc Account) (err error) {
	err = DB.Select(accountList, "SELECT * FROM account WHERE linked = $1 AND suspended = false ORDER BY name;", acc.ID)
	*accountList = append([]Account{acc}, []Account(*accountList)...)
	return
}

func (accountList *AccountList) GetAllAccounts() (err error) {
	err = DB.Select(accountList, "SELECT * FROM account ORDER BY name;")
	return
}

func (accountList *AccountList) GetAllAccountsNotSuspended() (err error) {
	err = DB.Select(accountList, "SELECT * FROM account WHERE suspended = false ORDER BY name;")
	return
}

func DoAccountsExist(displayNames []string, suspendedAllowed bool) (b bool, err error) {
	var list AccountList

	err = DB.Select(&list, "SELECT * FROM account WHERE name = ANY($1) AND ( $2 OR suspended = false ) ORDER BY name;", pq.StringArray(displayNames), suspendedAllowed)
	if len(displayNames) == len(list) {
		return true, err
	}

	for _, item := range list {
		displayNames = helper.RemoveFirstStringOccurrenceFromArray(displayNames, item.DisplayName)
	}

	return false, errors.New(displayNames[0])
}

func (list *NameList) GetAllUserAndDisplayName() (err error) {
	err = DB.Select(list, "SELECT name, username FROM account ORDER BY name;")
	return
}
