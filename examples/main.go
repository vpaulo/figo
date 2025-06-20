package main

import (
	"fmt"
	"maps"

	"github.com/vpaulo/figo"
)

func main() {
	figma := figo.Figma{
		Prefix: "vp",
	}

	// Get Figma data from JSON file(API response saved in a JSON file)
	file, err := figma.GetDataFromFile("./tmp/original_output.json")
	if err != nil {
		fmt.Println("Error fetching Figma file:", err)
		return
	}

	// Get Figma Variables from JSON file(API response saved in a JSON file)
	variables, err := figma.GetVariablesFromFile("./tmp/variable_output.json")
	if err != nil {
		fmt.Println("Error fetching Figma Variables:", err)
		return
	}

	// Tokens from Variables API
	variableTokens := figma.ParseVariables(variables)

	// Tokens from File API
	tokens := figma.ParseTokens(file)

	// Merge both tokens
	maps.Copy(tokens, variableTokens)

	components := figma.ParseComponents(file, tokens)

	componentsCSS, err := figma.GenerateComponentsCSS(components)
	if err != nil {
		fmt.Println("Error generating components CSS:", err)
		return
	}

	componentsHTML, err := figma.GenerateComponentsHTML(components)
	if err != nil {
		fmt.Println("Error generating components HTML:", err)
		return
	}

	fmt.Printf("[componentsCSS] : %+v \n\n", componentsCSS)
	fmt.Printf("[componentsHTML] : \n%+v \n\n", componentsHTML)
}
