package dataLogic

import (
	"API_MBundestag/database"
	"API_MBundestag/help/generics"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Account struct {
	ID                     int64
	DisplayName            string
	Flair                  string
	ChangeFlair            bool
	Username               string
	Suspended              bool
	Role                   database.RoleString
	Linked                 int64
	RemoveFromTitle        bool
	RemoveFromOrganisation bool
}

var CouldNotFindAccount generics.Message = "Der Account konnte nicht gefunden werden"
var CouldFindAccount generics.Message = "Der Account wurde gefunden"

func (acc *Account) GetUser(displayName string, username string, msg *generics.Message, positiv *bool) {
	userLock.Lock()
	defer userLock.Unlock()
	var get = database.Account{}
	var err error
	switch true {
	case displayName == "":
		err = get.GetByUserName(username)
	case username == "":
		err = get.GetByDisplayName(displayName)
	default:
		err = gorm.ErrRecordNotFound
	}
	if err != nil {
		*msg = CouldNotFindAccount + "\n" + *msg
		return
	}
	*msg = CouldFindAccount + "\n" + *msg
	*positiv = true
	acc.ID = get.ID
	acc.Flair = get.Flair
	acc.ChangeFlair = false
	acc.Username = get.Username
	acc.DisplayName = get.DisplayName
	acc.Suspended = get.Suspended
	acc.Role = get.Role
	acc.Linked = get.Linked.Int64
}

var CouldNotChangeAccount generics.Message = "Der Account konnte nicht geändert werden"
var CouldChangeAccount generics.Message = "Der Account wurde geändert"
var AccountRetainsOrgs generics.Message = "Account konnte nicht von allen Organisationen entfernt werden"
var AccountRetainsTitles generics.Message = "Account konnte nicht von allen Titeln enfernt werden"

func (acc *Account) ChangeUser(msg *generics.Message, positiv *bool) {
	userLock.Lock()
	defer userLock.Unlock()
	var change = database.Account{}
	err := change.GetByUserName(acc.Username)
	if err != nil {
		*msg = CouldNotFindAccount + "\n" + *msg
		return
	}
	if acc.ChangeFlair {
		change.Flair = acc.Flair
	}
	change.Suspended = acc.Suspended
	change.Role = acc.Role
	change.Linked.Int64 = acc.Linked
	var linked = database.Account{}
	err = linked.GetByID(acc.Linked)
	change.Linked.Valid = err == nil && change.Role == database.PressAccount
	err = change.SaveChanges()
	if err != nil {
		*msg = CouldNotChangeAccount + "\n" + *msg
		return
	}
	*msg = CouldChangeAccount + "\n" + *msg
	if acc.RemoveFromTitle {
		err = RemoveSingleAccountFromTitles(&change)
		if err != nil {
			*msg = AccountRetainsTitles + "\n" + *msg
			return
		}
	}
	if acc.RemoveFromOrganisation {
		err = RemoveSingleAccountFromOrganisations(&change)
		if err != nil {
			*msg = AccountRetainsOrgs + "\n" + *msg
			return
		}
	}
	*positiv = true
	//Set to the actual values
	acc.ID = change.ID
	acc.Flair = change.Flair
	acc.Username = change.Username
	acc.DisplayName = change.DisplayName
	acc.Suspended = change.Suspended
	acc.Role = change.Role
	acc.Linked = change.Linked.Int64
	acc.ChangeFlair = false
	acc.RemoveFromTitle = false
	acc.RemoveFromOrganisation = false
}

var AccountCloudNotBeFound generics.Message = "Dein Account konnte nicht gefunden werden"
var OldPasswordNotcorrect generics.Message = "Das alte Password ist nicht korrekt"
var CouldNotHashPassword generics.Message = "Das neue Passwort konnte nicht korrekt gehashed werden"
var AccountPasswordCouldNotBeSaved generics.Message = "Es ist ein Fehler beim verändern des Passwords aufgetreten"
var AccountPasswordSuccessfulChanged generics.Message = "Password erfolgreich angepasst"

func ChangePassword(displayName string, oldPassword string, newPassword string, msg *generics.Message, positiv *bool) {
	userLock.Lock()
	defer userLock.Unlock()
	acc := database.Account{}
	err := acc.GetByDisplayName(displayName)
	if err != nil {
		*msg = AccountCloudNotBeFound + "\n" + *msg
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(acc.Password), []byte(oldPassword))
	if err != nil {
		*msg = OldPasswordNotcorrect + "\n" + *msg
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		*msg = CouldNotHashPassword + "\n" + *msg
		return
	}

	acc.Password = string(hash)
	err = acc.SaveChanges()
	if err != nil {
		*msg = AccountPasswordCouldNotBeSaved + "\n" + *msg
		return
	}
	*msg = AccountPasswordSuccessfulChanged + "\n" + *msg
	*positiv = true
}
