package htmlWrapper

import (
	"errors"
	"regexp"
	"sort"
	"strings"
)

type (
	byteList    []byte
	elementList []Elements
	Elements    struct {
		Name       string
		Text       string
		Attributes []string
	}
)

func (elementSlice *elementList) fill(dir string, extension string) (err error) {
	var elements map[string][]byte
	elements, err = loadTemplateFiles(dir, "elements", extension)
	if err != nil {
		return
	}
	err = elementSlice.extractElements(elements)
	if err != nil {
		return
	}
	return elementSlice.cleanElements()
}

func removeElement(elementName string, content string, item string) string {

	var re = regexp.MustCompile(`(?s)<template-` + elementName + regexp.QuoteMeta(content) + `<\/template-` + elementName + `\s*?>\s*`)
	var substitution = ""

	res := re.ReplaceAllString(item, substitution)
	return res
}

func (elementSlice *elementList) extractElements(e map[string][]byte) (err error) {
	for _, i := range e {
		fileContent := strings.TrimSpace(string(i))
		err = elementSlice.extractElementFromFile(fileContent)
		if err != nil {
			return
		}
	}
	return nil
}

func (elementSlice *elementList) extractElementFromFile(fileContent string) (err error) {
	for fileContent != "" {
		nameAndParms := strings.TrimSpace(getStringInBetween(fileContent, "<template-", ">"))
		if !strings.Contains(nameAndParms, "=") {
			content := getStringInBetween(fileContent, "<template-"+nameAndParms, "</template-"+nameAndParms)
			newElement := Elements{
				Name:       nameAndParms,
				Attributes: []string{},
				Text:       strings.TrimSpace(strings.TrimLeft(content, "> ")),
			}
			fileContent = removeElement(newElement.Name, content, fileContent)
			*elementSlice = append(*elementSlice, newElement)
		} else {
			split := strings.Split(nameAndParms, "=")
			if len(split) != 2 {
				return errors.New("there are more attributes then allowed")
			}
			if !strings.Contains(split[0], "info") {
				return errors.New("the incorrect attribute has been provided")
			}
			//cleaning name string part
			name := strings.TrimSuffix(strings.ReplaceAll(split[0], " ", ""), "info")
			//removes " from the info string and then all spaces (to be able to cleanly split it)
			//before splitting it into a list from the comma seperation
			attr := strings.Split(strings.ReplaceAll(strings.ReplaceAll(split[1], "\"", ""), " ", ""), ",")
			content := getStringInBetween(fileContent, "<template-"+name, "</template-"+name)
			newElement := Elements{
				Name:       name,
				Attributes: attr,
				Text:       strings.TrimSpace(strings.TrimLeft(getStringInBetween(fileContent, "<template-"+nameAndParms, "</template-"+name), "> ")),
			}
			fileContent = removeElement(newElement.Name, content, fileContent)
			*elementSlice = append(*elementSlice, newElement)
		}
	}
	return nil
}

func (elementSlice *elementList) cleanElements() (err error) {
	for num, e := range *elementSlice {
		err = elementSlice.needReplacement(num, e)
		if err != nil {
			return
		}
	}
	if canFindElementInList(*elementSlice) {
		return errors.New("there seems to be a recursion in the templates")
	}
	return nil
}

func (elementSlice *elementList) needReplacement(num int, e Elements) (err error) {
	i := 0
	for {
		val, stat := elementSlice.canFindElementInElement(e)
		if !stat {
			break
		}
		err = (*elementSlice)[num].replaceElement(val)
		if err != nil {
			return errors.New("error: \"" + err.Error() + "\" occured in template: \"" + e.Name + "\" with template: \"" + val.Name + "\"")
		}
		err = e.replaceElement(val)
		if i > 500 {
			return errors.New("there seems to be a recursion in the templates")
		}
		i++
	}
	return nil
}

func (element *Elements) replaceElement(e Elements) error {
	val := getStringInBetween(element.Text, "<"+e.Name, "</"+e.Name)
	replace, err := e.getReplacerString(val)
	if err != nil {
		return err
	}
	var re = regexp.MustCompile(`<` + e.Name + regexp.QuoteMeta(val) + `<\/` + e.Name + `\s*?>`)
	replace = strings.ReplaceAll(replace, "$", "$$")
	element.Text = re.ReplaceAllString(element.Text, replace)
	return nil
}

func (element *byteList) replaceElement(e Elements) error {
	val := getStringInBetween(string(*element), "<"+e.Name, "</"+e.Name)
	replace, err := e.getReplacerString(val)
	if err != nil {
		return err
	}
	var re = regexp.MustCompile(`<` + e.Name + regexp.QuoteMeta(val) + `<\/` + e.Name + `\s*?>`)
	replace = strings.ReplaceAll(replace, "$", "$$")
	*element = byteList(re.ReplaceAllString(string(*element), replace))
	return nil
}

func (elementSlice *elementList) canFindElementInElement(e Elements) (Elements, bool) {
	for _, e2 := range *elementSlice {
		if strings.Contains(e.Text, "<"+e2.Name) {
			return e2, true
		}
	}
	return Elements{}, false
}

func (element *Elements) getReplacerString(str string) (string, error) {
	if len(element.Attributes) == 0 {
		return element.Text, nil
	}
	re := regexp.MustCompile(`(?ms)[>,](.*?)=\|\|(.*?)\|\|`)
	matchList := re.FindAllSubmatch([]byte(str), -1)
	if len(matchList) == 0 {
		re = regexp.MustCompile(`(?ms)[>,](.*?)=([^,<]*)`)
		matchList = re.FindAllSubmatch([]byte(str), -1)
	}
	//check if any attributes are not covered
	var list []int
	notAllCovered := false
	if len(element.Attributes) != len(matchList) {
		notAllCovered = true
	}
	result := element.Text
	for _, val := range matchList {
		if len(val) != 3 {
			return "", errors.New("provided attribute incorrect")
		}
		valSplit := []string{string(val[1]), string(val[2])}
		name := strings.TrimSpace(valSplit[0])
		if i := getPositionOfString(element.Attributes, name); i == -1 {
			return "", errors.New("provided attribute does not exist")
		} else {
			// add covered attributes
			list = append(list, i)
		}
		value := strings.TrimSpace(valSplit[1])
		result = replaceElement(result, name, value)
		//replace all not specified attributes with empty string
	}
	sort.Ints(list)
	for i := 0; i < len(element.Attributes) && notAllCovered; i++ {
		if len(list) != 0 && list[0] == i {
			list = list[1:]
			continue
		}
		result = replaceElement(result, element.Attributes[i], "")
	}
	return result, nil
}

func replaceElement(input string, value string, replacer string) string {
	var re = regexp.MustCompile(`\${\s*?` + value + `\s*?}`)
	var substitution = strings.ReplaceAll(replacer, "$", "$$")
	val := re.ReplaceAllString(input, substitution)
	return val
}

func canFindElementInList(list []Elements) bool {
	for _, e1 := range list {
		for _, e2 := range list {
			if strings.Contains(e1.Text, "<"+e2.Name) {
				return true
			}
		}
	}
	return false
}

func (elementSlice *elementList) replaceAllElements(htmlList *map[string][]byte) (err error) {
	for a, b := range *htmlList {
		b, err = elementSlice.replaceSingleElement(a, b)
		(*htmlList)[a] = b
	}
	return
}

func (elementSlice *elementList) replaceSingleElement(a string, b []byte) (res []byte, err error) {
	list := byteList(b)
	for _, e := range *elementSlice {
		for strings.Contains(string(list), "<"+e.Name) {
			err = list.replaceElement(e)
			if err != nil {
				return []byte{}, errors.New("error: \"" + err.Error() + "\" occured in page: \"" + a + "\" with template: \"" + e.Name + "\"")
			}
		}
	}
	return list, nil
}

func getStringInBetween(str string, start string, end string) (result string) {
	s := strings.Index(str, start)
	if s == -1 {
		return
	}
	s += len(start)
	e := strings.Index(str[s:], end)
	if e == -1 {
		return
	}
	return str[s : s+e]
}

func getPositionOfString(input []string, value string) int {
	for p, v := range input {
		if v == value {
			return p
		}
	}
	return -1
}
