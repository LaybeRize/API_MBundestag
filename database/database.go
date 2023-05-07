package database

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"strings"
	"sync"
	"time"
)

var db *gorm.DB
var first = false
var parallel = sync.Mutex{}

func TestSetup() {
	parallel.Lock()
	defer parallel.Unlock()

	var err error
	if db == nil {
		db, err = gorm.Open(postgres.Open("user=postgres password=root dbname=test sslmode=disable"), &gorm.Config{})
		if err != nil {
			fmt.Println(err)
		}
	}

	if first {
		return
	}
	first = true
	err = db.Migrator().DropTable(&Account{}, &Letter{}, "letter_account", &Publication{}, &Article{}, &Title{}, "title_account",
		&Organisation{}, "organisation_member", "organisation_admins", "organisation_account",
		&Document{}, "doc_viewer", "doc_poster", "doc_allowed", &Zwitscher{}, &Votes{}, &Chat{})
	if err != nil {
		fmt.Println(err)
	}
	err = db.AutoMigrate(&Account{}, &Letter{}, &Publication{}, &Article{}, &Title{}, &Organisation{}, &Document{}, &Zwitscher{}, &Votes{}, &Chat{})
	if err != nil {
		fmt.Println(err)
	}

	acc := &Account{
		Username:    "head_admin",
		DisplayName: "head_admin",
		Password:    "head_admin",
		Role:        HeadAdmin,
	}
	err = acc.CreateMe()
	if err != nil {
		fmt.Println(err)
	}
}

func Setup() {
	var err error
	db, err = gorm.Open(postgres.Open(fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DOCKER"),
		os.Getenv("DB_NAME"))), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Fatal()
	}

	err = db.AutoMigrate(&Account{}, &Letter{}, &Publication{}, &Article{}, &Title{}, &Organisation{}, &Document{}, &Zwitscher{}, &Votes{}, &Chat{})
	if err != nil {
		fmt.Println(err)
	}

	acc := Account{}
	err = acc.GetByID(1)
	if err == gorm.ErrRecordNotFound {
		createHeadAdmin()
	} else if err != nil {
		log.Fatal(err)
	}

	createNormalArticle()
}

func createHeadAdmin() {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(os.Getenv("INIT_PASSWORD")), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}

	acc := Account{
		DisplayName: strings.Trim(os.Getenv("INIT_NAME"), " '\""),
		Flair:       "",
		Username:    os.Getenv("INIT_USERNAME"),
		Password:    string(hashedPassword),
		Role:        HeadAdmin,
	}

	err = acc.CreateMe()
	if err != nil {
		log.Fatal(err)
	}
}

func createNormalArticle() {
	pub := Publication{}

	err := pub.GetByID(EternatityPublicationName)
	if err == nil {
		return
	} else if err != gorm.ErrRecordNotFound {
		log.Fatal(err)
	}

	pub = Publication{
		UUID:         EternatityPublicationName,
		PublishTime:  time.Now().UTC(),
		Publicated:   false,
		BreakingNews: false,
	}

	err = pub.CreateMe()
	if err != nil {
		log.Fatal(err)
	}
}
