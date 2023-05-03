package generics

//**error**

//article create

var TextOrHeadlineAreEmpty = "Inhalt oder Titel des Artikels nicht angegeben"
var ArticleLimitError = "Die Überschrift des Artikels überschreitet das Zeichenlimit von %d"
var ErrorWhileCreatingArticle = "Bei der Erstellung des Artikels ist ein Fehler aufgetreten"
var ArticleContentLimit = 15000
var ArticleTitleLimit = 100
var ArticleSubtitleLimit = 200

//newspaper

var ErrorWhileLoadingNewsPaper = "Es ist ein Fehler beim laden der Zeitungen aufgetreten"

//publication

var PublicationDoesNotExistsOrNotAllowedToView = "Diese Zeitung existiert nicht, oder du besitzt keine Berechtigung sie anzuschauen"
var ErrorWhileLoadingArticles = "Es ist ein Fehler beim Laden der Artikel aufgetreten"
var ErrorWhilePublishingNews = "Es ist ein Fehler beim veröffentlichen der Zeitung aufgetreten"
var NewsIsAlreadyPublished = "Eine Zeitung kann nicht erneut veröffentlicht werden"

//publication and rejection

var FormatTimeForArticle = "Verfasst am 02.01.2006 um 15:04 Uhr"

//rejection

var CanNotFindArticle = "Artikel existiert nicht"
var CanNotFindPublicationForArticle = "Konnte Zeitung des Artikels nicht finden"
var ArticleAlreadyPublished = "Artikel ist bereits veröffentlicht, damit kann er nicht mehr zurückgewiesen werden"
var LetterRejectText = "Untertitel: %s\n\n```\n%s\n```\n\n# Begründung der Moderation\n\n%s"
var LetterRejectTitle = "Der Artikel \"%s\" wurde abgelehnt"
var AuthorQualityCheck = "der Qualitätskontrolle"
var RejectionCouldNotBeCreated = "Der Ablehnungsbrief konnte nicht korrekt erstellt werden"
var CouldNotDeleteArticle = "Artikel konnte nicht gelöscht werden"
var CouldNotDeletePublication = "Die Zeitung zum Artikel konnte nicht korrekt gelöscht werden"

//**success**

//article create

var SuccessfulCreateArticle = "Artikel wurde erfolgreich erstellt"

//**format**
//newspaper and publication

var FormatHiddenNormalNews = "Ausstehende Zeitung"
var FormatHiddenBreakingNews = "Eilmeldung geschrieben am 02.01.2006 um 15:04"
var FormatNormalNews = "Zeitung vom 02.01.2006"
var FormatBreakingNews = "Eilmeldung vom 02.01.2006 um 15:04 Uhr"
