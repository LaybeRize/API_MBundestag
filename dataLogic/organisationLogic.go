package dataLogic

import (
	"API_MBundestag/database"
	"API_MBundestag/help"
	"database/sql"
	"fmt"
)

type Organsation struct {
	Name      string
	MainGroup string
	SubGroup  string
	Flair     string
	Status    database.StatusString
	Member    []string
	Admins    []string
}

func (org *Organsation) GetMe(name string) (err error) {
	get := database.Organisation{}
	err = get.GetByName(name)
	if err != nil {
		return
	}
	*org = Organsation{
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
	return
}

var AccountDoesNotExistError = "Account \"%s\" existiert nicht"
var ErrorWhileQueryingAccounts help.Message = "Fehler beim Abfragen der Accounts"
var ErrorWhileCreatingOrganisation help.Message = "Fehler beim Erstellen der Organisation"
var ErrorWhileAddingFlair help.Message = "Fehler beim Hinzufägen des Flairs"
var ErrorOrganisationNotFound = "Die Organisation \"%s\" existiert nicht"
var ErrorWhileChangingOrganisation help.Message = "Fehler beim Ändern der Organisation"
var ErrorWhileRemovingFlair help.Message = "Fehler beim Entfernen des Flairs"
var OrganisationSuccessfulCreated help.Message = "Die Organisation wurde erfolgreich erstellt"
var SucessfulChangedOrganisation help.Message = "Die Organisation wurde erfolgreich geändert"

func (org *Organsation) CreateMe(msg *help.Message, positiv *bool) {
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

func tryCreation(d *database.Organisation, msg *help.Message, positiv *bool) bool {
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

func addUsersTo(user []string, accounts *database.AccountList, msg *help.Message, accMap *map[string]*database.Account) bool {
	exists, err := accounts.DoAccountsExist(user)
	if err != nil {
		if !exists {
			*msg = help.Message(fmt.Sprintf(AccountDoesNotExistError, err.Error())) + "\n" + *msg
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

func (org *Organsation) ChangeMe(msg *help.Message, positiv *bool) {
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

func tryChanging(d *database.Organisation, msg *help.Message, positiv *bool) bool {
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

func (org *Organsation) getOrganistion(d *database.Organisation, old *database.Organisation, msg *help.Message) bool {
	err := d.GetByName(org.Name)
	err2 := old.GetByName(org.Name)
	if err != nil || err2 != nil {
		*msg = help.Message(fmt.Sprintf(ErrorOrganisationNotFound, org.Name)) + "\n" + *msg
		return true
	}
	return false
}

func (org *Organsation) Update(d *database.Organisation) bool {
	d.MainGroup = org.MainGroup
	d.SubGroup = org.SubGroup
	d.Flair = sql.NullString{Valid: org.Flair != "", String: org.Flair}
	d.Status = org.Status
	return false
}

func (org *Organsation) addFlairs(i *[]database.Account, msg *help.Message, positiv *bool) {
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

func removeFlairsOrganisation(i *[]database.Account, old *database.Organisation, msg *help.Message, positiv *bool) {
	for _, acc := range *i {
		err := removeFlairWithSave(old.Flair.String, &acc)
		if err != nil {
			*msg = ErrorWhileRemovingFlair + "\n" + *msg
			*positiv = false
			return
		}
	}
}

func (org *Organsation) ChangeOnlyMembers(msg *help.Message, positiv *bool) {
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

func addToMap(accounts *[]database.Account, msg *help.Message, accMap *map[string]*database.Account) bool {
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
