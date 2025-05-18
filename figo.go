package figo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"reflect"
	"sort"
	"strings"
	"time"

	"github.com/vpaulo/figo/figma"
	fg "github.com/vpaulo/figo/figma"
)

type Figma struct {
	FILE_KEY string
	API_KEY  string
}

func (figma *Figma) getUri() (string, error) {
	// TODO: download maybe become very large with geometry=path param
	component_url := `https://api.figma.com/v1/files/{{.FILE_KEY}}?geometry=paths`
	// t, parsingFailure := template.New("figma_uri").Parse(component_url)
	// if parsingFailure != nil {
	// 	return "", parsingFailure
	// }

	t := fg.CreateTmpl("figma_uri", component_url)

	var result bytes.Buffer

	err := t.Execute(&result, figma)
	if err != nil {
		return "", err
	}

	return result.String(), nil
}

func (figma *Figma) getVariablesUri() (string, error) {
	component_url := `https://api.figma.com/v1/files/{{.FILE_KEY}}/variables/local`
	t := fg.CreateTmpl("figma_uri", component_url)

	var result bytes.Buffer

	err := t.Execute(&result, figma)
	if err != nil {
		return "", err
	}

	return result.String(), nil
}

func (f *Figma) GetData() (figma.File, error) {
	var file figma.File

	// Create a new HTTP client with a timeout
	client := &http.Client{
		Timeout: 10 * time.Second, // may need longer timeout as figma files tend to get big
	}

	uri, uriError := f.getUri()
	if uriError != nil {
		return file, uriError
	}

	req, requestError := http.NewRequest("GET", uri, nil)
	if requestError != nil {
		return file, requestError
	}

	req.Header.Set("X-Figma-Token", f.API_KEY)

	httpResp, httpError := client.Do(req)
	if httpError != nil {
		return file, httpError
	}

	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		return file, fmt.Errorf("HTTP status code is %+v", httpResp.StatusCode)
	}

	body, readBodyError := io.ReadAll(httpResp.Body)
	if readBodyError != nil {
		return file, readBodyError
	}

	if unmarshallingError := json.Unmarshal(body, &file); unmarshallingError != nil {
		return file, unmarshallingError
	}

	return file, nil
}

func (f *Figma) GetVariablesData() (figma.Variables, error) {
	var variables figma.Variables

	// Create a new HTTP client with a timeout
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	uri, uriError := f.getVariablesUri()
	if uriError != nil {
		return variables, uriError
	}

	req, requestError := http.NewRequest("GET", uri, nil)
	if requestError != nil {
		return variables, requestError
	}

	req.Header.Set("X-Figma-Token", f.API_KEY)

	httpResp, httpError := client.Do(req)
	if httpError != nil {
		return variables, httpError
	}

	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		return variables, fmt.Errorf("HTTP status code is %+v", httpResp.StatusCode)
	}

	body, readBodyError := io.ReadAll(httpResp.Body)
	if readBodyError != nil {
		return variables, readBodyError
	}

	if unmarshallingError := json.Unmarshal(body, &variables); unmarshallingError != nil {
		return variables, unmarshallingError
	}

	return variables, nil
}

func (f *Figma) GetDataFromFile(path string) (figma.File, error) {
	var file figma.File

	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Println("Error Reading Figma file:", err)
		return file, err
	}

	if unmarshallingError := json.Unmarshal(data, &file); unmarshallingError != nil {
		fmt.Println("Error Unmarshal Figma file:", unmarshallingError)
		return file, unmarshallingError
	}

	return file, nil
}

func (f *Figma) GetVariablesFromFile(path string) (figma.Variables, error) {
	var variables figma.Variables

	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Println("Error Reading Figma Variables:", err)
		return variables, err
	}

	if unmarshallingError := json.Unmarshal(data, &variables); unmarshallingError != nil {
		fmt.Println("Error Unmarshal Figma Variables:", unmarshallingError)
		return variables, unmarshallingError
	}

	return variables, nil
}

func (f *Figma) Pages(file figma.File) []figma.Node {
	var pages []figma.Node

	for _, page := range file.Document.Children {
		pages = append(pages, page)
	}

	return pages
}

func (f *Figma) ParseTokens(file figma.File) map[string]figma.Token {
	tokens := make(map[string]figma.Token)
	pages := f.Pages(file)
	styles := file.Styles

	for _, page := range pages {
		children := page.Children

		for _, node := range children {
			f.mapTokens(node, &styles, &tokens)
		}
	}

	return tokens
}

func (f *Figma) mapTokens(node figma.Node, styles *map[string]figma.Style, tokens *map[string]figma.Token) {
	if node.IsFrame() {
		for key, id := range node.Styles {
			_, hasToken := (*tokens)[id]
			s, hasStyle := (*styles)[id]

			if !hasToken && hasStyle {
				var value string
				var className string
				variable, theme := figma.TokenValues(s.Name)
				switch key {
				case "fills":
					value = node.Background()
				case "strokes":
					value = node.BorderColor()
				case "effect":
					value = node.BoxShadow()
				case "grid":
					value = ""
					className = ""
				}

				if value != "" {
					token := figma.Token{
						Name:      s.Name,
						Variable:  variable,
						Value:     value,
						Theme:     theme,
						ClassName: className,
					}

					(*tokens)[id] = token
				}
			}
		}

		for _, child := range node.Children {
			if child.IsText() {
				for key, id := range child.Styles {
					_, hasToken := (*tokens)[id]
					s, hasStyle := (*styles)[id]

					if !hasToken && hasStyle {
						var value string
						var className string
						variable, theme := figma.TokenValues(s.Name)
						switch key {
						case "text":
							value = child.Font()
							className = fmt.Sprintf("text__style--%v", figma.ToKebabCase(s.Name))
							theme = ""
						case "fill":
							value = child.Background()
						case "stroke":
							value = child.BorderColor()
						case "effect":
							value = child.BoxShadow()
						}

						if value != "" {
							token := figma.Token{
								Name:      s.Name,
								Variable:  variable,
								Value:     value,
								Theme:     theme,
								ClassName: className,
							}

							(*tokens)[id] = token
							// fmt.Printf("[token] : %+v \n\n", token)
						}
					}
				}
			}
			f.mapTokens(child, styles, tokens)
		}
	}
}

func (f *Figma) GenerateTokensCSS(tokens map[string]figma.Token) (string, error) {
	tk := make(map[string][]string)

	// Group tokens by theme and generate CSS rules
	for _, token := range tokens {
		if _, exists := tk[token.Theme]; !exists && token.Theme != "" {
			var rules []string
			for _, t := range tokens {
				if t.Theme == token.Theme {
					rules = append(rules, fmt.Sprintf("%s: %s;", t.Variable, t.Value))
				}
			}
			// Sort rules alphabetically (case-insensitive)
			sort.Slice(rules, func(i, j int) bool {
				return strings.ToLower(rules[i]) < strings.ToLower(rules[j])
			})
			tk[token.Theme] = removeDuplicates(rules)
		} else if _, exists := tk[token.ClassName]; !exists && token.ClassName != "" {
			tk[token.ClassName] = strings.Split(token.Value, "|")
		}
	}

	var out bytes.Buffer
	tmp := figma.CreateTmpl("tokens", figma.CssVariablesTemplate)
	err := tmp.Execute(&out, tk)
	if err != nil {
		return "", err
	}

	return out.String(), nil
}

func (f *Figma) ParseVariables(variables figma.Variables) map[string]figma.Token {
	collections := variables.Meta.VariableCollections
	vars := variables.Meta.Variables
	tokens := make(map[string]figma.Token)
	themes := getVariablesThemes(collections)

	if len(collections) == 0 || len(vars) == 0 {
		return tokens
	}

	for _, v := range vars {
		if (v.ResolvedType == fg.ResolvedTypeColor || v.ResolvedType == fg.ResolvedTypeFloat) && !v.DeletedButReferenced {
			collectionName := fg.ToKebabCase(collections[v.VariableCollectionId].Name)
			varName := fg.ToKebabCase(v.Name)

			for key, value := range v.ValuesByMode {
				result := ""
				theme := themes[key]
				switch reflect.TypeOf(value).Kind() {
				case reflect.Float64:
					result = fmt.Sprintf("%vpx", value)
				case reflect.Map:
					jsonData, _ := json.Marshal(value) // encode back to JSON

					var color fg.Color
					var alias fg.VariableAlias

					json.Unmarshal(jsonData, &alias)
					json.Unmarshal(jsonData, &color)

					result = color.Rgba()

					if alias.ID != "" && vars[alias.ID].VariableCollectionId != "" {
						prefix := fg.ToKebabCase(collections[vars[alias.ID].VariableCollectionId].Name)
						result = fmt.Sprintf("var(--%v-%v)", prefix, fg.ToKebabCase(vars[alias.ID].Name))
					}
				}

				if result != "" {
					token := figma.Token{
						Name:     v.Name,
						Variable: fmt.Sprintf("--%v-%v", collectionName, varName),
						Value:    result,
						Theme:    theme,
					}

					tokens[v.ID] = token
				}
			}
		}
	}

	return tokens
}

func getVariablesThemes(collections map[string]fg.VariableCollection) map[string]string {
	themes := make(map[string]string)

	for _, c := range collections {
		for _, mode := range c.Modes {
			name := ":root"
			if strings.HasSuffix(strings.ToLower(mode.Name), "theme") {
				name = fg.ToKebabCase(mode.Name)
			}

			themes[mode.ModeId] = name
		}
	}

	return themes
}

// TODO: move this functions to a common place
func removeDuplicates(input []string) []string {
	seen := make(map[string]bool)
	var result []string

	for _, val := range input {
		if !seen[val] {
			seen[val] = true
			result = append(result, val)
		}
	}

	return result
}
