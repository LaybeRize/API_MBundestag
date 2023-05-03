package dataLogic

import "API_MBundestag/database"

func GetAllAccountNamesNotSuspended() (acc []string, err error) {
	accs := database.AccountList{}
	err = accs.GetAllAccountsNotSuspended()
	if err != nil {
		return
	}
	acc = make([]string, len(accs))
	for i, e := range accs {
		acc[i] = e.DisplayName
	}
	return
}
