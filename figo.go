package figo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"time"

	"github.com/vpaulo/figo/figma"
)

type Figma struct {
	FILE_KEY string
	API_KEY  string
}

func (figma *Figma) getUri() (string, error) {
	component_url := `https://api.figma.com/v1/files/{{.FILE_KEY}}?geometry=paths`
	t, parsingFailure := template.New("figma_uri").Parse(component_url)
	if parsingFailure != nil {
		return "", parsingFailure
	}

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
		Timeout: 10 * time.Second, // may need lnger timeout as figma files tend to get big
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

func (f *Figma) Pages(file figma.File) []figma.Node {
	var pages []figma.Node

	for _, page := range file.Document.Children {
		pages = append(pages, page)
	}

	fmt.Printf("[PAGES] : %v \n\n", pages)

	return pages
}

func (f *Figma) Normalise(file figma.File) {
	fmt.Printf("[FILE] : %v \n\n", file)
	// componentSets := file.ComponentSets
	// component := file.Components
	// styles := file.Styles
	pages := f.Pages(file)

	fmt.Printf("[PAGES] : %v \n\n", pages)
}
