package generics

//**error**
//edit and create
//org

var ViewerError = "Etwas ist schiefgelaufen dabei die Hauptaccounts zuordnen zu können"
var OrganisationCreationError = "Es ist ein Fehler beim erstellen der Organisation aufgetreten"
var OrganisationEditError = "Es ist ein fehler beim Abspeichern der Veränderungen aufgetreten"
var OrgFindingError = "Konnte die Organisation mit diesem Namen nicht finden"
var OrgEditNonExistantElement = "Eine Organisation die nicht exisitiert kann nicht angepasst werden"

//title

var TitelUpdateError = "Es ist ein Fehler beim bearbeiten des Titels ausgetreten.\n Das kann unter anderem daran liegen, dass der Flair bereits vergeben ist"
var TitleDoesNotExists = "Angegebener Titel konnte nicht gefunden werden"
var ErrorWhileDeletingTitle = "Beim löschen des Titels ist ein Fehler aufgetreten"
var TitleCreationError = "Titel konnte nicht erstellt werden.\n Das kann unter anderem daran liegen, dass der Flair bereits vergeben ist"
var RefresingTitleHierachyDidNotWork = "Während des updaten der Titelübersicht ist ein Fehler aufgetreten"

//both

var FlairUpdateError = "Beim updaten der Flairs der Nutzer ist ein Fehler aufgetreten"

//view
//org

var ErrorWhileLodingOrganisationView = "Es ist ein Fehler beim laden der Organisationen aufgetreten"

//**success**

//org

var SuccessFullCreationOrg = "Organisation wurde erfolgreich erstellt"
var SuccessFullFindOrg = "Konnte Organisation finden"
var SuccessFullChangeOrg = "Organisation erfolgreich verändert"

//title

var SuccessFullEditTitle = "Titel erfolgreich verändert"
var SuccessFullFoundTitle = "Titel erfolgreich gefunden"
var SuccesfulDeletedTitle = "Titel erfolgreich gelöscht"
var SuccessFullCreationTitle = "Titel erfolgreich erstellt"
