package htmlBasics

import (
	"API_MBundestag/dataLogic"
	"API_MBundestag/database"
	gen "API_MBundestag/help/generics"
	"API_MBundestag/htmlHandler"
	"database/sql"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

type StartPageStruct struct {
	LoggedIn bool
	Account  database.Account
	gen.MessageStruct
}

func GetStartPage(c *gin.Context) {
	acc, b := dataLogic.CheckUserPrivileged(c, database.User, database.MediaAdmin, database.Admin, database.HeadAdmin)

	htmlHandler.MakeSite(&StartPageStruct{LoggedIn: b,
		Account: acc,
		MessageStruct: gen.MessageStruct{
			Message: "",
			Positiv: false,
		},
	}, c, &acc)
}

func PostStartLogout(c *gin.Context) {
	var err error
	acc, b := dataLogic.CheckUserPrivileged(c, database.User, database.MediaAdmin, database.Admin, database.HeadAdmin)
	if b {
		err = dataLogic.ClearCookieFromUser(acc)
	}
	//if the cookies could not be cleared correctly give back an error page
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	htmlHandler.MakeSite(&StartPageStruct{LoggedIn: false,
		Account: acc,
		MessageStruct: gen.MessageStruct{
			Message: gen.SuccessfullLoggedOut,
			Positiv: true,
		},
	}, c, &database.Account{Role: database.NotLoggedIn})
}

func PostStartPage(c *gin.Context) {
	//logout routing
	if c.Query("type") == "logout" {
		PostStartLogout(c)
		return
	}
	var err error
	page := validateUserLogin(c)
	if page.LoggedIn == true {
		err = dataLogic.RequestToSetCookieForAccount(page.Account, c)
	}
	//if login was successfull but an error occured on RequestToSetCookieForAccount
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	htmlHandler.MakeSite(page, c, &page.Account)
}

func validateUserLogin(c *gin.Context) *StartPageStruct {
	var err error
	acc := database.Account{}
	//no need for checks if either password or username was left blank
	if c.PostForm("username") == "" || c.PostForm("password") == "" {
		return getLoggedOutStartPageStruct(gen.PasswordOrUsernameNotTypedIn)
	}
	//check if user account exists
	err = acc.GetByUserName(c.PostForm("username"))
	if err == sql.ErrNoRows {
		return getLoggedOutStartPageStruct(gen.PasswordOrUsernameWrong)
	}
	//if the database throws an error other than object not found, return an Internal Error
	if err != nil {
		return getLoggedOutStartPageStruct(gen.InternalValidationError)
	}
	//if the login block timer has not run out yet, return the time until it runs out
	if acc.NextLoginTime.Valid && !acc.NextLoginTime.Time.Before(time.Now().UTC()) {
		return getLoggedOutStartPageStruct(gen.Message(
			acc.NextLoginTime.Time.Format(gen.FormatStringForLoginTimeout)))
	}
	//check password
	err = bcrypt.CompareHashAndPassword([]byte(acc.Password), []byte(c.PostForm("password")))
	if err != nil {
		//if the password is wrong update the login tries and return the correct error message
		if success := updateLoginTries(&acc); success != nil {
			if success == AccountCanNotBeLoggindBecauseOfTimeout {
				return getLoggedOutStartPageStruct(gen.Message(
					acc.NextLoginTime.Time.Format(gen.FormatStringForLoginTimeout)))
			} else {
				return getLoggedOutStartPageStruct(gen.InternalValidationError)
			}
		}
		return getLoggedOutStartPageStruct(gen.PasswordOrUsernameWrong)
	}

	if acc.Suspended {
		return getLoggedOutStartPageStruct(gen.AccountIsSuspended)
	}
	//reset account login tries and make the login timer invalid before returning the correct struct
	acc.LoginTries = 0
	acc.NextLoginTime.Valid = false
	err = acc.SaveChanges()
	if err != nil {
		return getLoggedOutStartPageStruct(gen.InternalValidationError)
	}

	return &StartPageStruct{
		Account:  acc,
		LoggedIn: true,
		MessageStruct: gen.MessageStruct{
			Message: gen.SuccessfulLoggedIn,
			Positiv: true,
		},
	}
}

func getLoggedOutStartPageStruct(info gen.Message) *StartPageStruct {
	return &StartPageStruct{
		Account:  database.Account{Role: database.NotLoggedIn},
		LoggedIn: false,
		MessageStruct: gen.MessageStruct{
			Message: info,
			Positiv: false,
		},
	}
}

var AccountCanNotBeLoggindBecauseOfTimeout = htmlHandler.ValidationErrors{Info: "account can not be logged in at the moment, because the timer has not run out"}

func updateLoginTries(acc *database.Account) (err error) {
	acc.LoginTries += 1
	//set the timer appropriate for the tries
	switch acc.LoginTries {
	case 1, 2, 3:
	case 4, 5:
		acc.NextLoginTime.Time = time.Now().UTC().Add(time.Second * 5)
	case 6, 7:
		acc.NextLoginTime.Time = time.Now().UTC().Add(time.Minute)
	case 8, 9:
		acc.NextLoginTime.Time = time.Now().UTC().Add(time.Minute * 5)
	default:
		min := acc.LoginTries * acc.LoginTries * 10
		acc.NextLoginTime.Time = time.Now().UTC().Add(time.Second * time.Duration(min))
	}
	//make it valid if it had been set
	if acc.LoginTries > 3 {
		acc.NextLoginTime.Valid = true
	}
	err = acc.SaveChanges()
	//check if the timer was saved correctly
	if err == nil && acc.LoginTries > 3 {
		err = AccountCanNotBeLoggindBecauseOfTimeout
	}
	return
}
