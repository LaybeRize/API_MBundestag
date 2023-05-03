package dataLogic

import (
	"API_MBundestag/database"
	"API_MBundestag/help"
	"regexp"
	"strings"
)

func addFlair(flair string, acc *database.Account) (err error) {
	if acc.Flair == "" {
		acc.Flair += flair
	} else {
		acc.Flair += ", " + flair
	}
	err = acc.SaveChanges()
	return
}

func removeFlair(flair string, acc *database.Account) {
	if acc.Flair == flair {
		acc.Flair = ""
	} else if strings.HasPrefix(acc.Flair, flair+",") {
		acc.Flair = help.TrimPrefix(acc.Flair, flair+", ")
	} else if strings.Contains(acc.Flair, ", "+flair+",") {
		var re = regexp.MustCompile(`(?m), ` + regexp.QuoteMeta(flair) + `,`)
		var substitution = ","
		acc.Flair = re.ReplaceAllString(acc.Flair, substitution)
	} else if strings.HasSuffix(acc.Flair, ", "+flair) {
		acc.Flair = help.TrimSuffix(acc.Flair, ", "+flair)
	}
}

func removeFlairWithSave(flair string, acc *database.Account) (err error) {
	removeFlair(flair, acc)
	err = acc.SaveChanges()
	return
}
