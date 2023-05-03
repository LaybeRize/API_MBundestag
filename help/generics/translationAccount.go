package generics

import "API_MBundestag/htmlHandler"

//**success**

//create

var SuccesFullCreatedAccount htmlHandler.Message = "Account wurde erfolgreich erstellt"

//**error**

//edit and create

var LinkedValueNotANumberError htmlHandler.Message = "Die eingetragene Zahl für die Account-Verlinkung kann nicht gelesen werden"
var RoleCanNotBeSelectedError htmlHandler.Message = "Die ausgewählte Rolle ist entweder nicht korrekt, oder dir fehlen die Berechtigung zur Zuweisung"

//edit

var CanNotChangeNoExistentAccount htmlHandler.Message = "Account kann nicht verändert werden, da er nicht existiert"
var DisallowedChangeToHeadAdmin htmlHandler.Message = "Du bist nicht dazu berechtigt Veränderungen an Head-Admins vorzunehmen"
var CanNotChangeRootAccount htmlHandler.Message = "Veränderungen am Root-Account sind nicht erlaubt"
var InvalidType htmlHandler.Message = "Der angegebene Suchtypus wird nicht unterstützt"

//create

var ErrorWhileGeneratingPasswordHash htmlHandler.Message = "Es ist ein Fehler beim generieren des Passwordhashs aufgetreten"
var UserOrDisplaynameAlreadyExist htmlHandler.Message = "Der Nutzer- oder Anzeigename des Account wurde bereits benutzt"
var NamesOrPasswordIsEmptyError htmlHandler.Message = "Nutzername, Anzeigename und Password drüfen nicht leer sein"

//account list

var AccountQueryError htmlHandler.Message = "Account-Details konnten der Datenbank nicht korrekt entnommen werden"

//password change

var NewPasswordIsNotTheSame htmlHandler.Message = "Das neue Passwort stimmt nicht mit der Wiederholung überein"
var MinPasswordLength = 10
var NewPasswordIsNotMinimumOf10Characters htmlHandler.Message = "Das neue Password muss mindestens 10 Zeichen haben"
var OldPasswordNotcorrect htmlHandler.Message = "Das alte Password ist nicht korrekt"
var ErrorWhileChangingPassword htmlHandler.Message = "Es ist ein Fehler beim ändern des Passworts aufgetreten"
var SuccessChangePassword htmlHandler.Message = "Dein Passwort wurde erfolgreich geändert"

var CouldNotLoadAccountDetails = "Details über dich konnten nicht geladen werden"
