package help

import (
	"fmt"
	"github.com/gomarkdown/markdown"
	"github.com/microcosm-cc/bluemonday"
	"log"
	"os"
	"regexp"
	"strings"
)

func ArrayContainsString(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

func GetPositionOfString(input []string, value string) int {
	for p, v := range input {
		if v == value {
			return p
		}
	}
	return -1
}

func RemoveFromArray(s []string, i int) []string {
	if i == -1 {
		return s
	}
	if i == 0 && len(s) == 1 {
		return []string{}
	}
	s[i] = s[0]
	return s[1:]
}

func RemoveFirstStringOccurrenceFromArray(s []string, str string) []string {
	i := GetPositionOfString(s, str)
	return RemoveFromArray(s, i)
}

func TrimSuffix(s, suffix string) string {
	s = s[:len(s)-len(suffix)]
	return s
}

func TrimPrefix(s, prefix string) string {
	s = s[len(prefix):]
	return s
}

func RemoveDuplicates(array []string) []string {
	var result []string
	result = []string{}
	for _, val := range array {
		if GetPositionOfString(result, val) == -1 {
			result = append(result, val)
		}
	}
	return result
}

func ClearStringArray(array *[]string) {
	clone := make([]string, len(*array))
	copy(clone, *array)
	*array = []string{}
	for _, str := range clone {
		if str != "" && GetPositionOfString(*array, str) == -1 {
			*array = append(*array, str)
		}
	}
}

func DeleteMultiplesAndEmpty(a []string) []string {
	ClearStringArray(&a)
	return a
}

func CreateHTML(md string) string {
	intermediate := markdown.NormalizeNewlines([]byte(md))
	maybeUnsafeHTML := markdown.ToHTML(intermediate, nil, nil)
	htmlResult := bluemonday.UGCPolicy().SanitizeBytes(maybeUnsafeHTML)
	return updateHtmlResult(htmlResult)
}

var ReplacerMap map[string]string

func updateHtmlResult(htmlResult []byte) string {
	result := string(htmlResult)
	result = strings.ReplaceAll(result, "<code>\n", "<code>")
	for key, val := range ReplacerMap {
		var withAttr = regexp.MustCompile(`(?m)(<` + regexp.QuoteMeta(key) + ` )`)
		var withoutAttr = regexp.MustCompile(`(?m)(<` + regexp.QuoteMeta(key) + `)>`)
		intermediate := fmt.Sprintf("$1 %s ", val)
		result = withAttr.ReplaceAllString(result, intermediate)
		intermediate = fmt.Sprintf("$1 %s>", val)
		result = withoutAttr.ReplaceAllString(result, intermediate)
	}
	return result
}

func UpdateAttributes() {
	ReplacerMap = make(map[string]string)
	var re = regexp.MustCompile(`(?m)<(\w*?) (.*?)>`)
	var getTemplate = regexp.MustCompile(`(?s)<!-- Test start -->(.*)<!-- Test end -->`)
	b, err := os.ReadFile("templates/includes/markdown.html")
	if err != nil {
		log.Fatalln(err)
	}
	b = getTemplate.FindAllSubmatch(b, -1)[0][1]
	for _, match := range re.FindAllSubmatch(b, -1) {
		ReplacerMap[string(match[1])] = string(match[2])
	}
}
