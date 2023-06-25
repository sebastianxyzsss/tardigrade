package parser

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/sebastianxyzsss/tardigrade/globals"
	"github.com/sebastianxyzsss/tardigrade/logger"
	"github.com/sebastianxyzsss/tardigrade/utils"
)

var ll = logger.SetupLog()

var HistoryTemp *Element = nil

// represents one element in the yaml file that could be a group or a command
type Element struct {
	Content        string
	IsCommand      bool
	Description    string
	Tags           []string
	ChildKeys      *[]string
	Parent         *Element
	ChildrenSorted []*Element
	Children       map[string]*Element
}

func (e Element) String() string {
	return fmt.Sprintf("{ Content:%s, IsComm:%s, Desc:%s, Tags:%s} ", e.Content, boolToString(e.IsCommand), e.Description, e.Tags)
}

func boolToString(b bool) string {
	if b {
		return "true"
	}
	return "false"
}

func SimpleElementToMap(element *Element) map[interface{}]interface{} {
	m := make(map[interface{}]interface{})
	simpleList := make([]interface{}, len(element.ChildrenSorted))
	for i, s := range element.ChildrenSorted {
		simpleList[i] = s.Content
	}
	m[element.Content] = simpleList
	return m
}

// main entry point for the parser
func MainRecurseMap(m map[interface{}]interface{}, parent *Element) error {

	RecurseMap(m, parent)

	HistoryTemp = parent.Children["history"]

	if globals.FlatParse {
		ll.Debug().Msg("Flat Parse !!")

		flatParent := NewFlatParent()

		PostProcess(parent, flatParent)

		parent.Children = flatParent.Children
		parent.ChildrenSorted = flatParent.ChildrenSorted

		if len(parent.ChildrenSorted) < 1 {
			return fmt.Errorf("No children found")
		}
	}

	return nil
}

func PostProcess(element *Element, flatParent *Element) {
	if element.Content == "history" {
		return
	}
	if element.IsCommand {
		if element.Parent != nil {
			if !DoesElementHaveCommandTag(element.Parent) {
				appendToFlatParent(flatParent, element)
			} else {
				content := utils.ReplaceContentWithChoices(element.Parent.Content, element.Content)
				element.Content = content
				appendToFlatParent(flatParent, element)
			}
		}
	}

	for _, child := range element.ChildrenSorted {
		PostProcess(child, flatParent)
	}
}

func appendToFlatParent(flatParent *Element, element *Element) {
	if elementPassesPreCriteria(element) {
		flatParent.ChildrenSorted = append(flatParent.ChildrenSorted, element)
		flatParent.Children[element.Content] = element
	}
}

func RecurseMap(m map[interface{}]interface{}, parent *Element) {
	for key, val := range m {
		valType := reflect.TypeOf(val)
		keyStr := fmt.Sprintf("%v", key)
		element := NewElement(keyStr, false, parent)
		if valType.Kind() == reflect.Map {
			processElement(element, parent)
			RecurseMap(val.(map[interface{}]interface{}), element)
		} else if valType.Kind() == reflect.Slice {
			processElement(element, parent)
			RecurseSlice(val.([]interface{}), element)
		}
	}
}

func RecurseSlice(lst []interface{}, parent *Element) {
	for _, item := range lst {
		itemStr := fmt.Sprintf("%v", item)
		element := NewElement(itemStr, true, parent)
		processElement(element, parent)
	}
}

func NewFlatParent() *Element {
	return NewElement("all", false, nil)
}

func NewElement(key string, isCommand bool, parent *Element) *Element {
	children := make(map[string]*Element)
	childrenSorted := make([]*Element, 0)
	childKeys := make([]string, 0)
	tags := make([]string, 0)
	element := Element{
		Content:        key,
		IsCommand:      isCommand,
		Tags:           tags,
		ChildKeys:      &childKeys,
		Children:       children,
		ChildrenSorted: childrenSorted,
		Parent:         parent,
	}
	return &element
}

func processElement(element *Element, parent *Element) {
	tokens := strings.Split(element.Content, "^")
	if len(tokens) < 2 { // without any comments
		currentContent := strings.TrimSpace(tokens[0])
		element.Content = currentContent
		processElementToParent(element, parent, currentContent)
	} else { // comments were found
		currentContent := strings.TrimSpace(tokens[0])
		currentDescription := strings.TrimSpace(tokens[1])
		element.Content = currentContent
		element.Description = currentDescription
		tags := getTagsFromDescription(currentDescription)
		element.Tags = tags
		processElementToParent(element, parent, currentContent)
	}
}

func processElementToParent(element *Element, parent *Element, elementContent string) {
	parentChildKeys := *parent.ChildKeys
	parentChildKeys = append(parentChildKeys, TruncateString(element.Content, globals.ChildKeyMaxSize))
	parent.ChildKeys = &parentChildKeys
	// copy children
	parent.Children[elementContent] = element
	parent.ChildrenSorted = append(parent.ChildrenSorted, element)
	// copy tags
	for _, parentTag := range parent.Tags {
		if parentTag != "comm" {
			element.Tags = append(element.Tags, parentTag)
		}
	}
}

func elementPassesPreCriteria(element *Element) bool {

	if globals.FilterAction == globals.FilterTags {
		if element.IsCommand == true && !AreFilterStringsInTags(element.Tags) {
			return false
		}
	}

	if globals.FilterAction == globals.FilterAnything {
		others := []string{element.Content, element.Description}
		if element.IsCommand == true && !AreFilterStringsInTags(element.Tags) && !AreFilterStringsInTags(others) {
			return false
		}
	}

	if strings.HasPrefix(element.Content, "(dummy)") { // dummy test element
		return false
	}
	return true
}

func AreFilterStringsInTags(strs []string) bool {

	if len(globals.FilterStrings) > 0 {
		for _, filterString := range globals.FilterStrings {
			for _, str := range strs {
				if strings.Contains(str, filterString) {
					ll.Debug().Str("str", str).Str("filterString", filterString).Msg("match")
					return true
				}
			}
		}
	}
	return false
}

func getTagsFromDescription(description string) []string {
	re := regexp.MustCompile("@\\S+")
	tags := re.FindAllString(description, -1)
	// cleaning at sign
	for i := 0; i < len(tags); i++ {
		tags[i] = strings.Replace(tags[i], "@", "", 1)
	}
	return tags
}

func DoesElementHaveCommandTag(element *Element) bool {
	doesIt := areAnyTagsCommand(element.Tags)
	return doesIt
}

func areAnyTagsCommand(tags []string) bool {
	for _, tag := range tags {
		if tag == "comm" {
			return true
		}
	}
	return false
}

func GetChildKeysAsString(element *Element) string {
	asString := strings.Join(getChildKeys(element), " | ")
	return "| " + asString + " |"
}

func getChildKeys(element *Element) []string {
	keys := make([]string, 0)
	if element.ChildKeys != nil {
		keys = *element.ChildKeys
	}
	return keys
}

func TruncateString(str string, length int) string {
	if length <= 0 {
		return ""
	}
	if length >= len(str) {
		return str
	}
	truncated := ""
	count := 0
	for _, char := range str {
		truncated += string(char)
		count++
		if count >= length {
			break
		}
	}
	return truncated + ".."
}
