package dataLogic

import (
	"API_MBundestag/database"
	"API_MBundestag/help/generics"
	"database/sql"
	"fmt"
)

type Organisation struct {
	Name      string
	MainGroup string
	SubGroup  string
	Flair     string
	Status    database.StatusString
	Member    []string
	Admins    []string
}

func (org *Organisation) GetMeWhenAdmin(name string, id int64) (err error) {
	get := database.Organisation{}
	err = get.GetByNameAndOnlyWhenAccountAsAdmin(name, id)
	if err != nil {
		return
	}
	org.translateTo(&get)
	return
}

func (org *Organisation) GetMe(name string) (err error) {
	get := database.Organisation{}
	err = get.GetByName(name)
	if err != nil {
		return
	}
	org.translateTo(&get)
	return
}

func (org *Organisation) translateTo(get *database.Organisation) {
	*org = Organisation{
		Name:      get.Name,
		MainGroup: get.MainGroup,
		SubGroup:  get.SubGroup,
		Flair:     get.Flair.String,
		Status:    get.Status,
		Member:    make([]string, len(get.Members)),
		Admins:    make([]string, len(get.Admins)),
	}
	for i, acc := range get.Members {
		org.Member[i] = acc.DisplayName
	}
	for i, acc := range get.Admins {
		org.Admins[i] = acc.DisplayName
	}
}

var AccountDoesNotExistError = "Account \"%s\" existiert nicht"
var ErrorWhileQueryingAccounts generics.Message = "Fehler beim Abfragen der Accounts"
var ErrorWhileCreatingOrganisation generics.Message = "Fehler beim Erstellen der Organisation"
var ErrorWhileAddingFlair generics.Message = "Fehler beim Hinzufägen des Flairs"
var ErrorOrganisationNotFound = "Die Organisation \"%s\" existiert nicht"
var ErrorWhileChangingOrganisation generics.Message = "Fehler beim Ändern der Organisation"
var ErrorWhileRemovingFlair generics.Message = "Fehler beim Entfernen des Flairs"
var OrganisationSuccessfulCreated generics.Message = "Die Organisation wurde erfolgreich erstellt"
var SucessfulChangedOrganisation generics.Message = "Die Organisation wurde erfolgreich geändert"

func (org *Organisation) CreateMe(msg *generics.Message, positiv *bool) {
	creation := database.Organisation{
		Name:      org.Name,
		MainGroup: org.MainGroup,
		SubGroup:  org.SubGroup,
		Flair:     sql.NullString{Valid: org.Flair != "", String: org.Flair},
		Status:    org.Status,
		Members:   []database.Account{},
		Admins:    []database.Account{},
		Accounts:  []database.Account{},
	}
	userLock.Lock()
	defer userLock.Unlock()
	accMap := map[string]*database.Account{}
	switch true {
	case addUsersTo(org.Member, (*database.AccountList)(&creation.Members), msg, &accMap):
	case addUsersTo(org.Admins, (*database.AccountList)(&creation.Admins), msg, &accMap):
	case addAccountsTo(&creation.Accounts, &accMap):
	case tryCreation(&creation, msg, positiv):
	default:
		if !creation.Flair.Valid {
			return
		}
		org.addFlairs(&creation.Members, msg, positiv)
		org.addFlairs(&creation.Admins, msg, positiv)
	}
	return
}

func tryCreation(d *database.Organisation, msg *generics.Message, positiv *bool) bool {
	err := d.CreateMe()
	if err != nil {
		*msg = ErrorWhileCreatingOrganisation + "\n" + *msg
		return true
	}
	*msg = OrganisationSuccessfulCreated + "\n" + *msg
	*positiv = true
	return false
}

func addAccountsTo(i *[]database.Account, m *map[string]*database.Account) bool {
	*i = make([]database.Account, len(*m))
	counter := 0
	for _, acc := range *m {
		(*i)[counter] = *acc
		counter++
	}
	return false
}

func addUsersTo(user []string, accounts *database.AccountList, msg *generics.Message, accMap *map[string]*database.Account) bool {
	exists, err := accounts.DoAccountsExist(user)
	if err != nil {
		if !exists {
			*msg = generics.Message(fmt.Sprintf(AccountDoesNotExistError, err.Error())) + "\n" + *msg
		} else {
			*msg = ErrorWhileQueryingAccounts + "\n" + *msg
		}
		return true
	}
	if accMap == nil {
		return false
	}
	for _, acc := range *accounts {
		err = acc.GetByDisplayNameWithParent(acc.DisplayName)
		if err != nil {
			*msg = ErrorWhileQueryingAccounts + "\n" + *msg
			return true
		}
		if acc.Parent == nil {
			(*accMap)[acc.DisplayName] = &acc
		} else {
			(*accMap)[acc.Parent.DisplayName] = acc.Parent
		}
	}
	return false
}

func (org *Organisation) ChangeMe(msg *generics.Message, positiv *bool) {
	orgnisationLock.Lock()
	userLock.Lock()
	defer orgnisationLock.Unlock()
	defer userLock.Unlock()
	change := database.Organisation{}
	old := database.Organisation{}
	accMap := map[string]*database.Account{}
	switch true {
	case org.getOrganistion(&change, &old, msg):
	case org.Update(&change):
	case addUsersTo(org.Member, (*database.AccountList)(&change.Members), msg, &accMap):
	case addUsersTo(org.Admins, (*database.AccountList)(&change.Admins), msg, &accMap):
	case addAccountsTo(&change.Accounts, &accMap):
	case tryChanging(&change, msg, positiv):
	default:
		if old.Flair.Valid {
			removeFlairsOrganisation(&old.Members, &old, msg, positiv)
			removeFlairsOrganisation(&old.Admins, &old, msg, positiv)
		}
		if change.Flair.Valid {
			org.addFlairs(&change.Members, msg, positiv)
			org.addFlairs(&change.Admins, msg, positiv)
		}
	}
	return
}

func tryChanging(d *database.Organisation, msg *generics.Message, positiv *bool) bool {
	switch true {
	case d.SaveChanges() != nil:
	case d.UpdateAccounts() != nil:
	case d.UpdateMembers() != nil:
	case d.UpdateAdmins() != nil:
	default:
		*msg = SucessfulChangedOrganisation + "\n" + *msg
		*positiv = true
		return false
	}
	*msg = ErrorWhileChangingOrganisation + "\n" + *msg
	return true
}

func (org *Organisation) getOrganistion(d *database.Organisation, old *database.Organisation, msg *generics.Message) bool {
	err := d.GetByName(org.Name)
	err2 := old.GetByName(org.Name)
	if err != nil || err2 != nil {
		*msg = generics.Message(fmt.Sprintf(ErrorOrganisationNotFound, org.Name)) + "\n" + *msg
		return true
	}
	return false
}

func (org *Organisation) Update(d *database.Organisation) bool {
	d.MainGroup = org.MainGroup
	d.SubGroup = org.SubGroup
	d.Flair = sql.NullString{Valid: org.Flair != "", String: org.Flair}
	d.Status = org.Status
	return false
}

func (org *Organisation) addFlairs(i *[]database.Account, msg *generics.Message, positiv *bool) {
	for _, acc := range *i {
		switch true {
		case acc.GetByDisplayName(acc.DisplayName) != nil:
			fallthrough
		case addFlair(org.Flair, &acc) != nil:
			*msg = ErrorWhileAddingFlair + "\n" + *msg
			*positiv = false
			return
		}
	}
}

func removeFlairsOrganisation(i *[]database.Account, old *database.Organisation, msg *generics.Message, positiv *bool) {
	for _, acc := range *i {
		err := removeFlairWithSave(old.Flair.String, &acc)
		if err != nil {
			*msg = ErrorWhileRemovingFlair + "\n" + *msg
			*positiv = false
			return
		}
	}
}

func (org *Organisation) ChangeOnlyMembers(msg *generics.Message, positiv *bool) {
	orgnisationLock.Lock()
	userLock.Lock()
	defer orgnisationLock.Unlock()
	defer userLock.Unlock()
	change := database.Organisation{}
	old := database.Organisation{}
	accMap := map[string]*database.Account{}
	switch true {
	case org.getOrganistion(&change, &old, msg):
	case addUsersTo(org.Member, (*database.AccountList)(&change.Members), msg, &accMap):
	case addToMap(&change.Admins, msg, &accMap):
	case addAccountsTo(&change.Accounts, &accMap):
	case tryChanging(&change, msg, positiv):
	default:
		if old.Flair.Valid {
			removeFlairsOrganisation(&old.Members, &old, msg, positiv)
		}
		if change.Flair.Valid {
			org.addFlairs(&change.Members, msg, positiv)
		}
	}
	return
}

func addToMap(accounts *[]database.Account, msg *generics.Message, accMap *map[string]*database.Account) bool {
	var err error
	for _, acc := range *accounts {
		err = acc.GetByDisplayNameWithParent(acc.DisplayName)
		if err != nil {
			*msg = ErrorWhileQueryingAccounts + "\n" + *msg
			return true
		}
		if acc.Parent == nil {
			(*accMap)[acc.DisplayName] = &acc
		} else {
			(*accMap)[acc.Parent.DisplayName] = acc.Parent
		}
	}
	return false
}
