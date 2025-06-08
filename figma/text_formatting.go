package figma

import (
	"fmt"
	"strings"
	"text/template"
	"unicode"
)

func normaliseWords(input string) []string {
	var result []rune
	for i, r := range input {
		if unicode.IsUpper(r) && i > 0 && (unicode.IsLower(rune(input[i-1])) || unicode.IsDigit(rune(input[i-1]))) {
			result = append(result, ' ')
		}
		result = append(result, r)
	}
	// Now normalize separators (space, underscore, dash) to space
	separators := []string{"_", "-", "/", " "} // TODO: may need more separators

	input = strings.ToLower(string(result))
	for _, sep := range separators {
		input = strings.ReplaceAll(input, sep, " ")
	}
	return strings.Fields(input)
}

func ToPascalCase(input string) string {
	words := normaliseWords(input)
	if len(words) == 0 {
		return ""
	}

	var pascal string
	for _, word := range words {
		if len(word) > 0 {
			pascal += strings.ToUpper(string(word[0])) + strings.ToLower(word[1:])
		}
	}

	return pascal
}

func ToCamelCase(input string) string {
	words := normaliseWords(input)
	if len(words) == 0 {
		return ""
	}

	// First word lowercase
	camel := strings.ToLower(words[0])

	// Capitalize the first letter of the rest
	for _, word := range words[1:] {
		if len(word) > 0 {
			camel += strings.ToUpper(string(word[0])) + strings.ToLower(word[1:])
		}
	}

	return camel
}

func ToKebabCase(input string) string {
	words := normaliseWords(input)
	if len(words) == 0 {
		return ""
	}

	return strings.Join(words, "-")
}

func TokenValues(name string) (string, string) {
	variable := fmt.Sprintf("--%v", ToKebabCase(name))
	theme := ":root"

	if strings.Contains(name, "/") {
		list := strings.Split(name, "/")

		if strings.Contains(list[0], "theme") {
			theme = ToKebabCase(list[0])
			variable = fmt.Sprintf("--%v", ToKebabCase(strings.Join(list[1:], "-")))
		}
	}

	return variable, theme
}

func CreateTmpl(name, t string) *template.Template {
	return template.Must(template.New(name).Parse(t))
}
