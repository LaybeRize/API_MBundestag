package dataLogic

import (
	"API_MBundestag/database"
	"API_MBundestag/help/generics"
	"database/sql"
)

type Title struct {
	OldName   string
	Name      string
	Flair     string
	MainGroup string
	SubGroup  string
	Holder    []string
}

func (title *Title) GetMe(name string) (err error) {
	db := database.Title{}
	err = db.GetByName(name)
	if err != nil {
		return
	}
	*title = Title{
		OldName:   db.Name,
		Name:      db.Name,
		Flair:     db.Flair.String,
		MainGroup: db.MainGroup,
		SubGroup:  db.SubGroup,
		Holder:    make([]string, len(db.Holder)),
	}
	for i, acc := range db.Holder {
		title.Holder[i] = acc.DisplayName
	}
	return
}

var ErrorTitleCouldNotBeCreated generics.Message = "Titel konnte nicht erstellt werden"
var SuccessCreatedTitle generics.Message = "Titel erfolgreich erstellt"
var ErrorWhileRefreshingTitleHierarchy generics.Message = "Titelhierarchie konnte nicht aktualisiert werden"
var ErrorCouldNotFindTitle generics.Message = "Titel konnte nicht gefunden werden"
var ErrorWhileChangingTitle generics.Message = "Titel konnte nicht geändert werden"
var SuccessChangedTitle generics.Message = "Titel erfolgreich geändert"

func (title *Title) CreateMe(msg *generics.Message, positiv *bool) {
	titleLock.Lock()
	userLock.Lock()
	defer titleLock.Unlock()
	defer userLock.Unlock()
	creation := database.Title{
		Name:      title.Name,
		MainGroup: title.MainGroup,
		SubGroup:  title.SubGroup,
		Flair:     sql.NullString{Valid: title.Flair != "", String: title.Flair},
		Holder:    []database.Account{},
	}
	switch true {
	case addUsersTo(title.Holder, (*database.AccountList)(&creation.Holder), msg, nil):
	case tryCreateTitle(&creation, msg, positiv):
	default:
		err := RefreshTitleHierarchy()
		if err != nil {
			*msg = ErrorWhileRefreshingTitleHierarchy + "\n" + *msg
			*positiv = false
		}
		if !creation.Flair.Valid {
			return
		}
		title.addFlairs(&creation.Holder, msg, positiv)
	}
}

func tryCreateTitle(creation *database.Title, msg *generics.Message, positiv *bool) bool {
	err := creation.CreateMe()
	if err != nil {
		*msg = ErrorTitleCouldNotBeCreated + "\n" + *msg
		return true
	}
	*msg = SuccessCreatedTitle + "\n" + *msg
	*positiv = true
	return false
}

func (title *Title) addFlairs(i *[]database.Account, msg *generics.Message, positiv *bool) {
	for _, acc := range *i {
		switch true {
		case acc.GetByDisplayName(acc.DisplayName) != nil:
			fallthrough
		case addFlair(title.Flair, &acc) != nil:
			*msg = ErrorWhileAddingFlair + "\n" + *msg
			*positiv = false
			return
		}
	}
}

func (title *Title) ChangeMe(msg *generics.Message, positiv *bool) {
	titleLock.Lock()
	userLock.Lock()
	defer titleLock.Unlock()
	defer userLock.Unlock()
	change := database.Title{
		Name:      title.Name,
		MainGroup: title.MainGroup,
		SubGroup:  title.SubGroup,
		Flair:     sql.NullString{Valid: title.Flair != "", String: title.Flair},
		Holder:    []database.Account{},
	}
	old := database.Title{}
	switch true {
	case getTitle(&old, title.OldName, msg):
	case addUsersTo(title.Holder, (*database.AccountList)(&change.Holder), msg, nil):
	case title.tryChangeTitle(&change, msg, positiv):
	default:
		err := RefreshTitleHierarchy()
		if err != nil {
			*msg = ErrorWhileRefreshingTitleHierarchy + "\n" + *msg
			*positiv = false
		}
		if old.Flair.Valid {
			removeFlairsTitle(old.Flair.String, &old.Holder, msg, positiv)
		}
		if change.Flair.Valid {
			title.addFlairs(&change.Holder, msg, positiv)
		}
	}
}

func (title *Title) tryChangeTitle(new *database.Title, msg *generics.Message, positiv *bool) bool {
	old := database.Title{Name: title.OldName, Holder: []database.Account{}}
	switch true {
	case old.UpdateHolder() != nil:
	case new.ChangeTitleName(title.OldName) != nil:
	case new.UpdateHolder() != nil:
	default:
		*msg = SuccessChangedTitle + "\n" + *msg
		*positiv = true
		return false
	}
	*msg = ErrorWhileChangingTitle + "\n" + *msg
	return true
}

func removeFlairsTitle(flair string, i *[]database.Account, msg *generics.Message, positiv *bool) {
	for _, acc := range *i {
		err := removeFlairWithSave(flair, &acc)
		if err != nil {
			*msg = ErrorWhileAddingFlair + "\n" + *msg
			*positiv = false
			return
		}
	}
}

func getTitle(old *database.Title, name string, msg *generics.Message) bool {
	err := old.GetByName(name)
	if err != nil {
		*msg = ErrorCouldNotFindTitle + "\n" + *msg
		return true
	}
	return false
}

func (title *Title) DeleteMe(msg *generics.Message, positiv *bool) {
	titleLock.Lock()
	userLock.Lock()
	defer titleLock.Unlock()
	defer userLock.Unlock()
	old := database.Title{}
	switch true {
	case getTitle(&old, title.OldName, msg):
	case tryDeleteMe(&old, msg, positiv):
	default:
		err := RefreshTitleHierarchy()
		if err != nil {
			*msg = ErrorWhileRefreshingTitleHierarchy + "\n" + *msg
			*positiv = false
		}
		if old.Flair.Valid {
			removeFlairsTitle(old.Flair.String, &old.Holder, msg, positiv)
		}
	}
}

var ErrorWhileDeletingTitle generics.Message = "Titel konnte nicht gelöscht werden"
var SuccessDeletedTitle generics.Message = "Titel erfolgreich gelöscht"

func tryDeleteMe(d *database.Title, msg *generics.Message, positiv *bool) bool {
	old := database.Title{Name: d.Name, Holder: []database.Account{}}
	switch true {
	case old.UpdateHolder() != nil:
		fallthrough
	case d.DeleteMe() != nil:
		*msg = ErrorWhileDeletingTitle + "\n" + *msg
		return true
	}
	*msg = SuccessDeletedTitle + "\n" + *msg
	*positiv = true
	return false
}
