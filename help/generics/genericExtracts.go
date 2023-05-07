package generics

import (
	"API_MBundestag/help"
	"database/sql"
	"github.com/gin-gonic/gin"
	"strconv"
	"strings"
)

func GetText(c *gin.Context, field string) string {
	return strings.TrimSpace(c.PostForm(field))
}

func GetNumber(c *gin.Context, field string, standard int, min int, max int) int {
	i, err := strconv.Atoi(strings.TrimSpace(c.PostForm(field)))
	if err != nil {
		return standard
	}
	if i < min {
		return min
	}
	if i > max {
		return max
	}
	return i
}

func GetBool(c *gin.Context, field string) bool {
	return c.PostForm(field) == "true"
}

func GetNullString(c *gin.Context, field string) sql.NullString {
	s := sql.NullString{
		String: GetText(c, field),
		Valid:  false,
	}
	s.Valid = s.String != ""
	return s
}

func GetStringArray(c *gin.Context, field string) []string {
	return help.DeleteMultiplesAndEmpty(c.PostFormArray(field))
}

// GetIfType querys c.Query for the string "type" and compares it to value
func GetIfType(c *gin.Context, value string) bool {
	return c.Query("type") == value
}

// GetIfEmptyQuery querys value and checks if the string is empty
func GetIfEmptyQuery(c *gin.Context, value string) bool {
	return c.Query(value) == ""
}
