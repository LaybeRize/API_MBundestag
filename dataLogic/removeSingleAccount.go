package dataLogic

import (
	"API_MBundestag/database"
)

func RemoveSingleAccountFromOrganisations(account *database.Account) (err error) {
	orgList := database.OrganisationList{}
	orgnisationLock.Lock()
	defer orgnisationLock.Unlock()
	err = orgList.GetAllPartOf(account.ID)
	if err != nil {
		return
	}
	err = database.DeleteMeFromOrganisations(account.ID)
	if err != nil {
		return
	}
	for _, org := range orgList {
		err = updateOrganisations(&org)
		if err != nil {
			return
		}
		if org.Flair.Valid {
			removeFlair(org.Flair.String, account)
		}
	}
	err = account.SaveChanges()
	return
}

func updateOrganisations(org *database.Organisation) (err error) {
	err = org.GetByName(org.Name)
	if err != nil {
		return
	}
	accMap := map[string]*database.Account{}
	for _, acc := range org.Members {
		err = acc.GetByDisplayNameWithParent(acc.DisplayName)
		if err != nil {
			return
		}
		if acc.Parent == nil {
			accMap[acc.DisplayName] = &acc
		} else {
			accMap[acc.Parent.DisplayName] = acc.Parent
		}
	}
	for _, acc := range org.Admins {
		err = acc.GetByDisplayNameWithParent(acc.DisplayName)
		if err != nil {
			return
		}
		if acc.Parent == nil {
			accMap[acc.DisplayName] = &acc
		} else {
			accMap[acc.Parent.DisplayName] = acc.Parent
		}
	}
	org.Accounts = make([]database.Account, len(accMap))
	counter := 0
	for _, acc := range accMap {
		org.Accounts[counter] = *acc
		counter++
	}
	err = org.UpdateAccounts()
	return
}

func RemoveSingleAccountFromTitles(acc *database.Account) (err error) {
	titleLock.Lock()
	defer titleLock.Unlock()
	list := database.TitleList{}
	err = list.GetAllForUserID(acc.ID)
	if err != nil {
		return
	}
	err = database.DeleteMeFromTitles(acc.ID)
	if err != nil {
		return
	}
	for _, t := range list {
		if t.Flair.Valid {
			removeFlair(t.Flair.String, acc)
		}
	}
	err = acc.SaveChanges()
	return
}
