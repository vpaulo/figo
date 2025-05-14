package figo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
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

func (f *Figma) Pages(file figma.File) []figma.Node {
	var pages []figma.Node

	for _, page := range file.Document.Children {
		pages = append(pages, page)
	}

	return pages
}

func (f *Figma) Normalise(file figma.File) {
	fmt.Printf("[FILE] : %v \n\n", file)
	// componentSets := file.ComponentSets
	// component := file.Components
	// styles := file.Styles
	// tokens := make(map[string]figma.Token)
	pages := f.Pages(file)

	fmt.Printf("[PAGES] : %+v \n\n", pages)
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
		// fmt.Printf("[tokens] : %+v \n\n", tokens)
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
			tk[token.Theme] = rules
		} else if _, exists := tk[token.ClassName]; !exists && token.ClassName != "" {
			tk[token.ClassName] = strings.Split(token.Value, "|")
			// fmt.Printf("[PORRA] : %+v \n\n", token)
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
