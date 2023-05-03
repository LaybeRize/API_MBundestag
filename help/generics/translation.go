package generics

/* genericChecks */
//error

var ContentAndTitelAreEmpty = "Inhalt und Titel dürfen nicht leer sein"
var ContentTooLong = "Der Inhalt überschreitet das Zeichenlimit von %d"
var TitleTooLong = "Der Titel überschreitet das Zeichenlimit von %d"
var SubtitleTooLong = "Der Untertitel überschreitet das Zeichenlimit von %d"
var AccountDoesNotExists = "Der ausgewählte Account existiert nicht"
var AccountIsNotYours = "Der ausgewählte Account steht dir nicht zur Verfügung"
var OrganisationDoesNotExist = "Die ausgewählte Organisation existiert nicht"
var AccountDoesNotExistError = "Account \"%s\" existiert nicht"
var StatusIsInvalid = "Der ausgewälte Status ist nicht zulässig"
var NoMainGroupSubGroupOrNameProvided = "Es wurde keine Hauptkategorie, Unterkategorie oder ein Name angegeben"

/* documentView */
//error

var DocumentDoesNotExists = "Dokument existiert nicht"
var TypeDoesNotExist = "Dieser Anfragetyp existiert nicht"

/* discussionCreateHandler */
//error

var TimeParseDiscussion = "2006-01-02T15:04"

/* postingDocumentHandler */
//error

var OrganisationInURLDoesNotExist = "Die Organisation in der URL existiert nicht"

/* postsCreateHandler */
//error

var PostTagLimit = 100
var PostTitleLimit = 100
var PostSubtitleLimit = 200
var PostContentLimit = 15000
var PostContentTooLong = "Der Inhalt überschreitet das Zeichenlimit von 15.000"
var PostTitleTooLong = "Der Titel überschreitet das Zeichenlimit von 100"
var PostSubtitleTooLong = "Der Untertitel überschreitet das Zeichenlimit von 200"
var TagTooLong = "Der Erstellungs-Tag überschreitet das Zeichenlimit von 100"
var PostIsNotAllowedWithEmptyTitleOrContent = "Ein Dokument darf keinen leeren Inhalt oder Titel haben"

var YouAreNotAllowedForOrganisation = "Du darfst in dieser Organisation nicht veröffentlichen"
var PostCreationFailed = "Es ist ein Fehler beim erstellen des Dokuments aufgetreten"
var SecretOrgsCanNotCreatePosts = "Geheime Organisationen können keine Posts veröffentlichen"

/* templater */
//error

var NamesQueryError = "Namen konnten der Datenbank nicht korrekt entnommen werden"
var OwnAccountsCouldNotBeFound = "Eigene Accounts konnten nicht gefunden werden"
var GroupQueryError = "Namen der Gruppenkategorien konnten der Datenbank nicht korrekt entnommen werden"
var OrgNamesQueryError = "Namen der Organisationen konnten der Datenbank nicht korrekt entnommen werden"

/* Multiple Ussage */
//for error page

var AccountDoesNotExistOrIsNotYours = "Dieser Account existiert nicht oder du bist nicht für ihn berechtigt"

//time formats

var LongTimeString = "02.01.2006 um 15:04 Uhr"

//for posting previews

var PreviewText = "Vorschau"

//for letter and articles

var WriteAccountDoesNotExistInDatabase = "Der ausgwählte Account für den Author, existiert nicht"
var ThatAccountIsNotAllowedToWrite = "Mit ausgewählte Account für den Autor darf nicht geschreiben werden"
