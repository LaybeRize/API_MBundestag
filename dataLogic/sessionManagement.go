package dataLogic

import (
	"API_MBundestag/database"
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"time"
)

// timeUntilTokenRunsOut defines the time in seconds until a token becomes invalid
var timeUntilTokenRunsOut = 60 * 60 * 24 * 7

// CheckUserPrivileged reads the cookie from the request via requestingAccount and evaluates if the account has
// one of the roles specified returning the result and the account.
func CheckUserPrivileged(c *gin.Context, roles ...database.RoleString) (database.Account, bool) {
	acc := requestingAccount(c)

	return acc, CheckIfHasRole(&acc, roles...)
}

func CheckIfHasRole(acc *database.Account, roles ...database.RoleString) bool {
	for _, r := range roles {
		if r == acc.Role {
			return true
		}
	}
	return false
}

// requestingAccount extracts the cookie from the request and finds the corresponding account, if the cookie
// exists. Then updates the cookie and adds it to the request.
func requestingAccount(c *gin.Context) database.Account {
	inCookie, err := c.Request.Cookie("token")

	if err != nil || inCookie == nil {
		return database.Account{Role: database.NotLoggedIn}
	}

	token := inCookie.Value

	account := database.Account{}
	userLock.Lock()
	defer userLock.Unlock()
	err = account.GetByToken(token)

	if err != nil {
		return database.Account{Role: database.NotLoggedIn}
	}

	if !account.ExpDate.Valid || account.ExpDate.Time.Before(time.Now().UTC()) || account.Suspended {
		return database.Account{Role: database.NotLoggedIn}
	}

	account.ExpDate.Time = time.Now().UTC().Add(time.Second * time.Duration(timeUntilTokenRunsOut))
	account.ExpDate.Valid = true
	account.RefToken.String = uuid.New().String()
	account.RefToken.Valid = true

	err = account.SaveChanges()
	if err != nil {
		return database.Account{Role: database.NotLoggedIn}
	}

	cookie := http.Cookie{Name: "token", Value: account.RefToken.String, Path: "/", MaxAge: timeUntilTokenRunsOut}

	http.SetCookie(c.Writer, &cookie)

	return account
}

func RequestToSetCookieForAccount(acc database.Account, c *gin.Context) (err error) {
	acc.ExpDate.Time = time.Now().UTC().Add(time.Second * time.Duration(timeUntilTokenRunsOut))
	acc.ExpDate.Valid = true
	acc.RefToken.String = uuid.New().String()
	acc.RefToken.Valid = true

	err = acc.SaveChanges()
	if err != nil {
		return
	}

	cookie := http.Cookie{Name: "token", Value: acc.RefToken.String, Path: "/", MaxAge: timeUntilTokenRunsOut}
	http.SetCookie(c.Writer, &cookie)
	return nil
}

func ClearCookieFromUser(acc database.Account) (err error) {
	acc.ExpDate = sql.NullTime{}
	acc.RefToken = sql.NullString{}

	err = acc.SaveChanges()
	return
}
