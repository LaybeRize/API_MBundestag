package generics

//**success**

//create

var SuccesFullCreatedAccount Message = "Account wurde erfolgreich erstellt"

//**error**

//edit and create

var LinkedValueNotANumberError Message = "Die eingetragene Zahl für die Account-Verlinkung kann nicht gelesen werden"
var RoleCanNotBeSelectedError Message = "Die ausgewählte Rolle ist entweder nicht korrekt, oder dir fehlen die Berechtigung zur Zuweisung"

//edit

var CanNotChangeNoExistentAccount Message = "Account kann nicht verändert werden, da er nicht existiert"
var DisallowedChangeToHeadAdmin Message = "Du bist nicht dazu berechtigt Veränderungen an Head-Admins vorzunehmen"
var CanNotChangeRootAccount Message = "Veränderungen am Root-Account sind nicht erlaubt"
var InvalidType Message = "Der angegebene Suchtypus wird nicht unterstützt"

//create

var ErrorWhileGeneratingPasswordHash Message = "Es ist ein Fehler beim generieren des Passwordhashs aufgetreten"
var UserOrDisplaynameAlreadyExist Message = "Der Nutzer- oder Anzeigename des Account wurde bereits benutzt"
var NamesOrPasswordIsEmptyError Message = "Nutzername, Anzeigename und Password drüfen nicht leer sein"

//account list

var AccountQueryError Message = "Account-Details konnten der Datenbank nicht korrekt entnommen werden"

//password change

var NewPasswordIsNotTheSame Message = "Das neue Passwort stimmt nicht mit der Wiederholung überein"
var MinPasswordLength = 10
var NewPasswordIsNotMinimumOf10Characters Message = "Das neue Password muss mindestens 10 Zeichen haben"
var OldPasswordNotcorrect Message = "Das alte Password ist nicht korrekt"
var ErrorWhileChangingPassword Message = "Es ist ein Fehler beim ändern des Passworts aufgetreten"
var SuccessChangePassword Message = "Dein Passwort wurde erfolgreich geändert"

var CouldNotLoadAccountDetails = "Details über dich konnten nicht geladen werden"
