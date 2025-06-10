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
	Prefix   string // Prefix for components tag
}

func (figma *Figma) getUri() (string, error) {
	// TODO: download maybe become very large with geometry=path param
	component_url := `https://api.figma.com/v1/files/{{.FILE_KEY}}?geometry=paths`

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

	figma.SetDefaults(&file)

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

	figma.SetDefaults(&variables)

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

	figma.SetDefaults(&file)

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

	figma.SetDefaults(&variables)

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
				case "grid": // TODO: get styles for grid
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

func (f *Figma) ParseComponents(file figma.File, tokens map[string]figma.Token) map[string]fg.Element {
	pages := f.Pages(file)
	components := f.initElementData(file)
	// var elements []fg.Element

	for _, page := range pages {
		children := page.Children

		for _, node := range children {
			if node.IsComponentOrSet() {
				element := components[node.ID]
				components[node.ID] = f.generateComponent(node.ID, node, figma.Node{}, "", &components, element, &tokens)

				// fmt.Printf("[yyy] : %+v \n\n", components[node.ID])
			}
		}
	}

	return components
}

func (f *Figma) initElementData(file figma.File) map[string]fg.Element {
	components := make(map[string]fg.Element)
	cmpSets := file.ComponentSets
	cmp := file.Components

	for key, set := range cmpSets {
		components[key] = fg.Element{
			Name: fmt.Sprintf("%v", fg.ToKebabCase(f.Prefix+" "+set.Name)),
		}
	}

	for key, c := range cmp {
		if c.ComponentSetId == "" {
			components[key] = fg.Element{
				Name: fmt.Sprintf("%v", fg.ToKebabCase(f.Prefix+" "+c.Name)),
			}
		} else {
			// components[key] = components[c.ComponentSetId]
		}
	}

	return components
}

func (f *Figma) generateComponent(id string, node figma.Node, parent figma.Node, parentClasses string, components *map[string]fg.Element, element fg.Element, tokens *map[string]figma.Token) fg.Element {
	// if id != "505:17" {
	// 	return element
	// }

	isMainComponent := false

	if element.Name != "" {
		isMainComponent = true
	}

	if isMainComponent {
		fmt.Printf("[IS COMPONENT] : %+v \n\n", element.Name)
	} else {
		fmt.Printf("[NOT COMPONENT] : %+v \n\n", node.Name)
		element.Name = fmt.Sprintf("%v", fg.ToKebabCase(node.Name))
	}

	element.Selectors = fmt.Sprintf("%v %v", parentClasses, node.Classes(f.Prefix, isMainComponent))

	if node.IsComponentSet() {
		fmt.Printf("[COMPONENT_SET] : %+v \n\n", (*components)[node.ID].Name)
		element.Variants = node.Variants()
	}
	if node.IsInstance() {
		fmt.Printf("[INSTANCE] : %+v \n\n", (*components)[node.ID].Name)
		return element
	}
	if node.IsComponent() {
		fmt.Printf("[COMPONENT] : %+v \n\n", (*components)[node.ID].Name)
		if parent.IsComponentSet() {
			fmt.Printf("[PARENT IS SET] : %+v \n\n", parent.Name)
			element.Name = fmt.Sprintf("%v", fg.ToKebabCase(f.Prefix+" "+parent.Name))
			element.Selectors = fmt.Sprintf("%v%v", parentClasses, node.Classes(f.Prefix, true))
		}
	}
	if !node.IsComponentSet() && !node.IsInstance() && !node.IsComponent() {
		fmt.Printf("[FRAME] : %+v %+v \n\n", node.Name, node.Type)
	}

	if !node.IsComponentSet() && !node.IsInstance() && !node.IsText() && !node.IsVector() {
		element.Styles = node.Css(parent)
	}

	if node.IsText() {
		// TODO: Update node.Font() to get the rest of the styles for a Text element and
		// and token if exists
		element.Styles = node.TextCss()
	}
	// TODO: vector styles
	// fmt.Printf("[STYLES] : %+v \n\n", el.Styles)
	//
	fmt.Printf("[ELEMENT] : %+v \n\n", element)

	for _, child := range node.Children {
		elem := f.generateComponent(id, child, node, element.Selectors, components, fg.Element{}, tokens)

		element.Children = append(element.Children, elem)
	}

	// fmt.Printf("[+++] : %+v \n\n", element)
	return element
}

func (f *Figma) GenerateComponentsCSS(components map[string]fg.Element) (string, error) {
	var styles []string

	for _, component := range components {
		if /*id == "505:17" &&*/ component.Selectors != "" {
			style, err := f.ComponentCSS(component)
			if err != nil {
				return "", err
			}
			styles = append(styles, strings.TrimSpace(style))
		}
	}

	return strings.Join(styles, "\n"), nil
}

func (f *Figma) ComponentCSS(component fg.Element) (string, error) {
	var out bytes.Buffer
	tmp := figma.CreateTmpl("components", figma.CssComponentsTemplate)
	err := tmp.Execute(&out, component)
	if err != nil {
		return "", err
	}

	return out.String(), nil
}
