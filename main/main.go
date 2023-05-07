package main

import (
	"API_MBundestag/dataLogic"
	"API_MBundestag/database"
	"API_MBundestag/help"
	"API_MBundestag/htmlHandler"
	"API_MBundestag/htmlHandler/htmlAccount"
	"API_MBundestag/htmlHandler/htmlBasics"
	"API_MBundestag/htmlHandler/htmlDocuments"
	"API_MBundestag/htmlHandler/htmlLetter"
	"API_MBundestag/htmlHandler/htmlPress"
	"API_MBundestag/htmlHandler/htmlWork"
	"API_MBundestag/htmlHandler/htmlZwitscher"
	"API_MBundestag/htmlHandler/websocket"
	wr "API_MBundestag/htmlWrapper"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"html/template"
	"log"
	"net/http"
	"os"
)

func main() {
	help.UpdateAttributes()
	setup()
}

func setup() {
	//gin.SetMode(gin.ReleaseMode)
	database.Setup()

	htmlAccount.Setup()
	htmlBasics.Setup()
	htmlDocuments.Setup()
	htmlLetter.Setup()
	htmlPress.Setup()
	htmlWork.Setup()
	htmlZwitscher.Setup()

	err := dataLogic.RefreshTitleHierarchy()
	if err != nil {
		log.Fatal(err)
	}

	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())

	err = router.SetTrustedProxies(nil)
	if err != nil {
		log.Fatal(err)
	}

	router.Static("/public", "./public")
	router.SetFuncMap(template.FuncMap{})
	templates, err := wr.New("templates", ".html", wr.DefaultFunctions)
	if err != nil {
		log.Fatal(err)
	}
	htmlHandler.Template = templates

	router.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/start")
	})
	router.GET("/login", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/start")
	})

	initRouter(router)

	err = router.Run(os.Getenv("ADRESS") + ":8080")
	if err != nil {
		log.Fatal(err)
	}
}

func initRouter(router *gin.Engine) {
	//router.GET("/reload/*path", htmlHandler.MiddleHardwareForTests)
	//router.GET("/preview", htmlBasics.GetPreviewPage)
	//router.POST("/preview", htmlBasics.PostPreviewPage)
	//TODO: make to comment all code above this in production
	//websockets
	router.GET("/chat/:token/:user", websocket.GetWebsocket)
	//json
	router.POST("/markdown", htmlBasics.PostJsonMarkdown)
	//html
	router.GET("/start", htmlBasics.GetStartPage)
	router.POST("/start", htmlBasics.PostStartPage)
	router.GET("/create-user", htmlAccount.GetCreateUserPage)
	router.GET("/edit-user", htmlAccount.GetEditUserPage)
	router.POST("/create-user", htmlAccount.PostCreateUserPage)
	router.POST("/edit-user", htmlAccount.PostEditUserPage)
	router.GET("/view-user", htmlAccount.GetAdminListUserPage)
	router.GET("/create-organisation", htmlWork.GetCreateOrganisationPage)
	router.POST("/create-organisation", htmlWork.PostCreateOrganisationPage)
	router.GET("/edit-organisation", htmlWork.GetEditOrganisationPage)
	router.POST("/edit-organisation", htmlWork.PostEditOrganisationPage)
	router.GET("/organisation", htmlWork.GetOrganisationViewPage)
	router.GET("/hidden-organisation", htmlWork.GetHiddenOrganisationViewPage)
	router.GET("/create-title", htmlWork.GetCreateTitlePage)
	router.POST("/create-title", htmlWork.PostCreateTitlePage)
	router.GET("/edit-title", htmlWork.GetEditTitlePage)
	router.POST("/edit-title", htmlWork.PostEditTitlePage)
	router.GET("/title", htmlWork.GetTitleViewPage)
	router.GET("/create-letter", htmlLetter.GetCreateLetterPage)
	router.POST("/create-letter", htmlLetter.PostCreateLetterPage)
	router.GET("/letter", htmlLetter.GetViewSingleLetter)
	router.GET("/create-mod-mail", htmlLetter.GetCreateModMailPage)
	router.POST("/create-mod-mail", htmlLetter.PostCreateModMailPage)
	router.GET("/create-article", htmlPress.GetCreateArticlePage)
	router.POST("/create-article", htmlPress.PostCreateArticlePage)
	router.GET("/admin-letter-view", htmlLetter.GetAdminLetterViewPage)
	router.POST("/admin-letter-view", htmlLetter.PostAdminLetterViewPage)
	router.GET("/newspaper-approval", htmlPress.GetNewsPaperHiddenListPage)
	router.GET("/newspaper", htmlPress.GetNewsPaperListPage)
	router.GET("/publication", htmlPress.GetPublicationViewPage)
	router.POST("/publication", htmlPress.PostPublicationViewPage)
	router.GET("/reject-article", htmlPress.GetRejectArticlePage)
	router.POST("/reject-article", htmlPress.PostRejectArticlePage)
	router.GET("/mod-mails", htmlLetter.GetViewModMailListPage)
	router.GET("/letter-list", htmlLetter.GetViewLetterListPage)
	router.POST("/letter-list", htmlLetter.PostViewLetterListPage)
	router.GET("/self-info", htmlAccount.GetViewOfProfilePage)
	router.GET("/password", htmlAccount.GetPasswordChangePage)
	router.POST("/password", htmlAccount.PostPasswordChangePage)
	router.GET("/zwitscher", htmlZwitscher.GetZwitscherLatestViewPage)
	router.POST("/zwitscher", htmlZwitscher.PostZwitscherLatestViewPage)
	router.GET("/create-post", htmlDocuments.GetPostsCreateHandler)
	router.POST("/create-post", htmlDocuments.PostPostsCreateHandler)
	router.GET("/create-discussion", htmlDocuments.GetDiscussionCreatePage)
	router.POST("/create-discussion", htmlDocuments.PostDiscussionCreatePage)
	router.GET("/create-vote", htmlDocuments.GetVoteCreatePage)
	router.POST("/create-vote", htmlDocuments.PostVoteCreatePage)
	router.GET("/create-document", htmlDocuments.GetDocumentNavigationPage)
	router.GET("/document", htmlDocuments.GetDocumentViewPage)
	router.POST("/document", htmlDocuments.PostDocumentViewPage)
	router.GET("/documents", htmlDocuments.GetDocumentListView)
	router.GET("/vote", htmlDocuments.GetVoteHandler)
	router.POST("/vote", htmlDocuments.PostVoteHandler)
	router.GET("/edit-user-organisation", htmlWork.GetOrganisationUserHandler)
	router.POST("/edit-user-organisation", htmlWork.PostOrganisationUserHandler)
}
