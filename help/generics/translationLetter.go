package generics

//**values**

var LetterTitleLimit = 200
var LetterContentLimit = 15000

//**error**
//create

var LetterCouldNotBeCreated = "Brief konnte nicht korrekt erstellt werden"
var AuthorEmptyError = "Bei einer Moderationsnachricht muss die Person, in dessen Namen versendet wird, eingetragen werden"

//view

var AccountForLetterViewError = "Der angegebene Account für die Einsicht des Briefes existiert nicht oder du bist nicht für ihn berechtigt"
var LetterDoesntExistOrNotAccessable = "Der Brief existiert nicht, oder du hast keine Berechtigung mit diesem Account darauf zuzugreifen"

//admin view

var ErrorUUIDDoesNotExist = "Die gesuchte UUID existiert nicht"

//list view

var ErrorWhileLoadingLetters = "Es ist ein Fehler beim laden der Briefe aufgetreten"
