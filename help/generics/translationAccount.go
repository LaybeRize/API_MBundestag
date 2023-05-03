package generics

import "API_MBundestag/help"

//**success**
//edit

var SuccessFullQueryedAccount help.Message = "Account konnte gefunden werden"
var SuccessFullChangedAccount help.Message = "Account wurde erfolgreich angepasst"
var SuccessFullChangedFlair help.Message = "Account-Flair wurde erfolgreich angepasst"

//create

var SuccesFullCreatedAccount help.Message = "Account wurde erfolgreich erstellt"

//**error**

//edit and create

var LinkedValueNotANumberError help.Message = "Die eingetragene Zahl für die Account-Verlinkung kann nicht gelesen werden"
var RoleCanNotBeSelectedError help.Message = "Die ausgewählte Rolle ist entweder nicht korrekt, oder dir fehlen die Berechtigung zur Zuweisung"

//edit

var AccountNotFoundError help.Message = "Account konnte nicht gefunden werden"
var CanNotChangeNoExistentAccount help.Message = "Account kann nicht verändert werden, da er nicht existiert"
var DisallowedChangeToHeadAdmin help.Message = "Du bist nicht dazu berechtigt Veränderungen an Head-Admins vorzunehmen"
var CanNotChangeRootAccount help.Message = "Veränderungen am Root-Account sind nicht erlaubt"
var InvalidType help.Message = "Der angegebene Suchtypus wird nicht unterstützt"

//create

var ErrorWhileGeneratingPasswordHash help.Message = "Es ist ein Fehler beim generieren des Passwordhashs aufgetreten"
var UserOrDisplaynameAlreadyExist help.Message = "Der Nutzer- oder Anzeigename des Account wurde bereits benutzt"
var NamesOrPasswordIsEmptyError help.Message = "Nutzername, Anzeigename und Password drüfen nicht leer sein"

//account list

var AccountQueryError help.Message = "Account-Details konnten der Datenbank nicht korrekt entnommen werden"

//password change

var NewPasswordIsNotTheSame help.Message = "Das neue Passwort stimmt nicht mit der Wiederholung überein"
var MinPasswordLength = 10
var NewPasswordIsNotMinimumOf10Characters help.Message = "Das neue Password muss mindestens 10 Zeichen haben"
var OldPasswordNotcorrect help.Message = "Das alte Password ist nicht korrekt"
var ErrorWhileChangingPassword help.Message = "Es ist ein Fehler beim ändern des Passworts aufgetreten"
var SuccessChangePassword help.Message = "Dein Passwort wurde erfolgreich geändert"

var CouldNotLoadAccountDetails = "Details über dich konnten nicht geladen werden"
