package htmlWrapper

import (
	"API_MBundestag/htmlWrapper/xfl"
	"regexp"
	"strings"
)

type HTMLMap map[string]xfl.HTMLItem

func (h *HTMLMap) cleanMap() {
	for k := range *h {
		h.cleanElement(k)
	}
}

func (h *HTMLMap) cleanElement(k string) {
	e := (*h)[k]
	keys := h.getKeys()
	for i := 0; i < len(keys); i++ {
		if keys[i] != k {
			e.Content = h.generateReplacment(keys[i], e.Content, keys[i])
		}
	}
	for i := len(keys) - 1; i >= 0; i-- {
		if keys[i] != k {
			e.Content = h.generateReplacment(keys[i], e.Content, keys[i])
		}
	}
	(*h)[k] = e
}

func (h *HTMLMap) getKeys() []string {
	keys := make([]string, len(*h))

	i := 0
	for k := range *h {
		keys[i] = k
		i++
	}
	return keys
}

func (h *HTMLMap) generateReplacment(errorIn string, html string, item string) string {
	var re = regexp.MustCompile(`(?mUs)<` + regexp.QuoteMeta(item) + `([^<]*/>|.*</` + regexp.QuoteMeta(item) + `\s*>)`)
	res := re.FindAllString(html, -1)
	if len(res) == 0 {
		return html
	}
	strMap := h.parseSingleTag(errorIn, item, res)
	for key, content := range strMap {
		html = strings.Replace(html, key, content, -1)
	}
	return html
}

func (h *HTMLMap) parseSingleTag(errorIn string, item string, array []string) map[string]string {
	element := (*h)[item]
	res := map[string]string{}
	for _, str := range array {
		res[str] = scanAndGetReplacement(errorIn, element, str)
	}
	return res
}

func scanAndGetReplacement(errorIn string, element xfl.HTMLItem, str string) string {
	attrMap := scanString(str, errorIn)
	return getSingleReplacerString(element, attrMap)
}

func scanString(str string, errorIn string) map[string]string {
	return xfl.ParseHTML(errorIn, str)
}

func getSingleReplacerString(element xfl.HTMLItem, attributes map[string]string) string {
	res := element.Content
	for _, key := range element.Attributes {
		res = strings.ReplaceAll(res, "#"+key+"#", attributes[key])
	}
	res = strings.ReplaceAll(res, "#content#", attributes["content"])
	return res
}

func (h *HTMLMap) replaceAllElements(htmlList *map[string][]byte, folder string) {
	for key, byteArray := range *htmlList {
		errorIn := key + " in folder " + folder
		b := string(byteArray)
		for mapKey := range *h {
			b = h.generateReplacment(errorIn, b, mapKey)
		}
		(*htmlList)[key] = []byte(b)
	}
}
