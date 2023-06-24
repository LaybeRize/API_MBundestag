package htmlWrapper

import (
	"API_MBundestag/database"
	"API_MBundestag/help/generics"
	"encoding/json"
	"fmt"
	"html/template"
	"net/url"
	"strings"
	"time"
	"unicode"
)

// Functions the original author found to be required in most every web-site template engine (extended by those, I need too)
// Many borrowed from https://github.com/Masterminds/sprig

// DefaultFunctions for templates
var DefaultFunctions = template.FuncMap{
	"add":               add,
	"optionValue":       optionValue,
	"showNumbers":       showNumbers,
	"showNames":         showNames,
	"notZero":           notZero,
	"arrayLengthEq":     arrayLengthEq,
	"lenStrNotZero":     lenStrNotZero,
	"option":            option,
	"ueq":               ueq,
	"showPublishButton": showPublishButton,
	"eqRole":            eqRole,
	"fullFillsClass":    fulfillsClass,
	"orgStatus":         orgStatus,
	"oneOfValues":       oneOfValues,
	"oneOfValuesArray":  oneOfValuesInArray,
	"yesno":             yesno,
	"plural":            plural,
	"dateFormat":        dateFormat,
	"withFlair":         withFlair,
	"unixtimestamp":     unixtimestamp,
	"json":              jsonFunc,
	/*	"sha256":             sha256Encoding,
		"sha1":               sha1Encoding,
		"md5":                md5Encoding,
		"base64encode":       base64encode,
		"base64decode":       base64decode,
		"base32encode":       base32encode,
		"base32decode":       base32decode, */
	"roleTranslations":   roleTranslations,
	"statusTranslations": statusTranslations,
	"voteTranslations":   voteTranslations,
	"arrayOrEmpty":       arrayOrEmpty,
	"userArrayOrEmpty":   userArrayOrEmpty,
	"title":              title,
	"messageToString":    messageToString,
	"noescape":           noescape,
	"getIcon":            getIcon,
	"noescapeurl":        noescapeurl,
	"queryEscape":        queryEscape,
	"getQueryString":     getQueryString,
	"valueLoop":          valueLoop,
	"headerLoop":         headerLoop,
	"userLoop":           userLoop,
	"roleLoop":           roleLoop,
	"statusLoop":         statusLoop,
	"voteLoop":           voteLoop,

	"upper":     strings.ToUpper,
	"lower":     strings.ToLower,
	"trim":      strings.TrimSpace,
	"urlencode": url.QueryEscape,
}

/******
FOR INT
******/

func add(i, j int) int {
	return i + j
}

func optionValue(m map[string]int, selected string) int {
	return m[selected]
}

/******
FOR BOOLEAN
******/

func showNumbers(v database.Votes) bool {
	return v.Finished || v.ShowNumbersWhileVoting
}

func showNames(v database.Votes) bool {
	return (v.Finished && v.ShowNamesAfterVoting) || (!v.Finished && v.ShowNamesWhileVoting)
}

func notZero(num int) bool {
	return num != 0
}

func arrayLengthEq(arr []database.Posts, num int) bool {
	return len(arr) == num
}

func lenStrNotZero(arr []string) bool {
	return len(arr) != 0
}

func option(m map[string]int, selected string, expected int) bool {
	return m[selected] == expected
}

func ueq(s1 string, s2 string) bool {
	return s1 != s2
}

func showPublishButton(pub bool, arts database.ArticleList) bool {
	if !pub && len(arts) != 0 {
		return true
	}
	return false
}

func eqRole(role database.RoleString, str interface{}) bool {
	switch str.(type) {
	case string:
		return string(role) == str
	case database.RoleString:
		return role == str
	default:
		return false
	}
}

func fulfillsClass(num int, role database.RoleString) bool {
	roleNum := 6
	switch role {
	case database.NotLoggedIn:
		roleNum = 6
	case database.User:
		roleNum = 5
	case database.MediaAdmin:
		roleNum = 4
	case database.Admin:
		roleNum = 3
	case database.HeadAdmin:
		roleNum = 2
	}
	return roleNum <= num
}

func orgStatus(str string, status database.StatusString) bool {
	return str == string(status)
}

func oneOfValues(str ...string) bool {
	for i, val := range str {
		if i == 0 {
			continue
		}
		if val == str[0] {
			return true
		}
	}
	return false
}

func oneOfValuesInArray(str string, strArray []string) bool {
	res := []string{str}
	res = append(res, strArray...)
	return oneOfValues(res...)
}

/******
FOR STRINGS
******/

// Often used for tables of rows
func yesno(yes string, no string, value bool) string {
	if value {
		return yes
	}
	return no
}

func plural(one, many string, count int) string {
	if count == 1 {
		return one
	}
	return many
}

// Current Date (Local server time)

func dateFormat(format string, timeValue time.Time) string {
	timeValue = timeValue.In(time.Local)
	return timeValue.Format(format)
}

func withFlair(author string, flair string) string {
	if flair == "" {
		return author
	}
	return author + ", " + flair
}

// Current Unix timestamp
func unixtimestamp() string {
	return fmt.Sprintf("%d", time.Now().Unix())
}

// json encodes an item into a JSON string
func jsonFunc(v interface{}) string {
	output, _ := json.Marshal(v)
	return string(output)
}

/*
// Modern Hash
func sha256Encoding(input string) string {
	hash := sha256.Sum256([]byte(input))
	return hex.EncodeToString(hash[:])
}

// Legacy
func sha1Encoding(input string) string {
	hash := sha1.Sum([]byte(input))
	return hex.EncodeToString(hash[:])
}

// Gravatar
func md5Encoding(input string) string {
	hash := md5.Sum([]byte(input))
	return hex.EncodeToString(hash[:])
}

// Popular encodings
func base64encode(v string) string {
	return base64.StdEncoding.EncodeToString([]byte(v))
}

func base64decode(v string) string {
	data, err := base64.StdEncoding.DecodeString(v)
	if err != nil {
		return err.Error()
	}
	return string(data)
}

func base32encode(v string) string {
	return base32.StdEncoding.EncodeToString([]byte(v))
}

func base32decode(v string) string {
	data, err := base32.StdEncoding.DecodeString(v)
	if err != nil {
		return err.Error()
	}
	return string(data)
}
*/

func roleTranslations(roleName interface{}) string {
	str := fmt.Sprintf("%v", roleName)
	return database.RoleTranslation[database.RoleString(str)]
}

func statusTranslations(statusName interface{}) string {
	str := fmt.Sprintf("%v", statusName)
	return database.StatusTranslation[database.StatusString(str)]
}
func voteTranslations(statusName interface{}) string {
	str := fmt.Sprintf("%v", statusName)
	return database.VoteTranslation[database.VoteType(str)]
}

func arrayOrEmpty(msg string, arr []string) string {
	if len(arr) == 0 {
		return msg
	}
	return strings.Join(arr, ", ")
}

func userArrayOrEmpty(msg string, arr []database.Account) string {
	if len(arr) == 0 {
		return msg
	}
	str := ""
	for _, acc := range arr {
		str += acc.DisplayName + ", "
	}
	return str[:len(str)-2]
}

func title(s string) string {
	// Use a closure here to remember state.
	// Hackish but effective. Depends on Map scanning in order and calling
	// the closure once per rune.
	prev := ' '
	return strings.Map(
		func(r rune) rune {
			if isSeparator(prev) {
				prev = r
				return unicode.ToTitle(r)
			}
			prev = r
			return r
		},
		s)
}

func messageToString(m generics.Message) string {
	return string(m)
}

// isSeperator is specifically for the european writing systems
func isSeparator(r rune) bool {
	// ASCII alphanumerics and underscore are not separators
	if r <= 0x7F {
		switch {
		case '0' <= r && r <= '9':
			return false
		case 'a' <= r && r <= 'z':
			return false
		case 'A' <= r && r <= 'Z':
			return false
		case r == '_':
			return false
		}
		return true
	}
	// Letters and digits are not separators
	if unicode.IsLetter(r) || unicode.IsDigit(r) {
		return false
	}
	// Otherwise, all we can do for now is treat spaces as separators.
	return unicode.IsSpace(r)
}

/******
FOR TEMPLATE HTML OR URL
******/

// Allow unsafe injection into HTML
func noescape(a ...interface{}) template.HTML {
	return template.HTML(fmt.Sprint(a...))
}

func getIcon(t database.DocumentType) template.HTML {
	switch t {
	case database.RunningDiscussion:
		fallthrough
	case database.FinishedDiscussion:
		return "<i class=\"bi bi-chat-right-text text-2xl pl-6\"></i>"
	case database.RunningVote:
		fallthrough
	case database.FinishedVote:
		return "<i class=\"bi bi-archive text-2xl pl-6\"></i>"
	case database.LegislativeText:
		return "<i class=\"bi bi-file-text text-2xl pl-6\"></i>"
	}
	return ""
}

// Allow unsafe URL injections into HTML
func noescapeurl(u string) template.URL {
	return template.URL(u)
}

func queryEscape(str string) template.URL {
	return template.URL(url.QueryEscape(str))
}

func getQueryString(uuid string, acc string, search bool) template.URL {
	if !search {
		return template.URL("uuid=" + url.QueryEscape(uuid))
	}
	return template.URL("uuid=" + url.QueryEscape(uuid) + "&usr=" + url.QueryEscape(acc))
}

/******
LOOPS FOR RANGES
******/

type ElementLoop struct {
	IsEnd  bool
	Header string
	Value  int
}

func valueLoop(array []string, m map[string]int) <-chan ElementLoop {
	loop := make(chan ElementLoop)
	go func() {
		end := len(array) - 1
		for i, str := range array {
			loop <- ElementLoop{
				IsEnd: i == end,
				Value: m[str],
			}
		}
		close(loop)
	}()
	return loop
}

func headerLoop(array []string) <-chan ElementLoop {
	loop := make(chan ElementLoop)
	go func() {
		end := len(array) - 1
		for i, str := range array {
			loop <- ElementLoop{
				IsEnd:  i == end,
				Header: str,
			}
		}
		close(loop)
	}()
	return loop
}

type UserLoop struct {
	Div    template.HTMLAttr
	Input  template.HTMLAttr
	Button template.HTMLAttr
}

func userLoop(uniqueID string, array []string) <-chan UserLoop {
	loop := make(chan UserLoop)
	go func() {
		loop <- UserLoop{
			Div:    template.HTMLAttr("style=\"display: none\" id=\"divClasses" + uniqueID + "\" hidden"),
			Input:  template.HTMLAttr("id=\"inputClasses" + uniqueID + "\""),
			Button: template.HTMLAttr("id=\"buttonClasses" + uniqueID + "\""),
		}
		for _, str := range array {
			str = template.HTMLEscapeString(str)
			loop <- UserLoop{
				Div:    "",
				Input:  template.HTMLAttr("type=\"text\" value=\"" + str + "\""),
				Button: "onclick=\"deleteDivFromSelf(this)\"",
			}
		}
		close(loop)
	}()
	return loop
}

type LoopStruct struct {
	Loop      string
	Attribute template.HTMLAttr
}

// generates a function that can be ranged over.
func roleLoop(consValue database.RoleString) <-chan LoopStruct {
	loop := make(chan LoopStruct)
	go func() {
		for _, str := range database.Roles {
			if str == string(consValue) {
				loop <- LoopStruct{Loop: str, Attribute: "selected"}
				continue
			}
			loop <- LoopStruct{Loop: str, Attribute: ""}
		}
		close(loop)
	}()
	return loop
}

func statusLoop(consValue database.StatusString) <-chan LoopStruct {
	loop := make(chan LoopStruct)
	go func() {
		for _, str := range database.Stati {
			if str == string(consValue) {
				loop <- LoopStruct{Loop: str, Attribute: "selected"}
				continue
			}
			loop <- LoopStruct{Loop: str, Attribute: ""}
		}
		close(loop)
	}()
	return loop
}

// generates a function that can be ranged over.
func voteLoop(consValue database.VoteType) <-chan LoopStruct {
	loop := make(chan LoopStruct)
	go func() {
		for _, str := range database.VoteTypes {
			if str == consValue {
				loop <- LoopStruct{Loop: string(str), Attribute: "selected"}
				continue
			}
			loop <- LoopStruct{Loop: string(str), Attribute: ""}
		}
		close(loop)
	}()
	return loop
}
