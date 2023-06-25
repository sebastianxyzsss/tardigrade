package utils

import "strings"

func BoolToCommand(b bool) string {
	if b {
		return "command"
	}
	return "group"
}

func ReplaceContentWithChoices(content string, choiceStr string) string {
	choices := strings.Split(choiceStr, ",")
	result := content
	for _, choice := range choices {
		result = strings.Replace(result, "<>", strings.TrimSpace(choice), 1)
	}
	return result
}
