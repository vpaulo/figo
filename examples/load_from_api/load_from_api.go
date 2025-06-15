package main

import (
	"fmt"
	"maps"

	"github.com/vpaulo/figo"
)

func main() {
	figma := figo.Figma{
		FILE_KEY: "FIGMA_FILE",
		API_KEY:  "YOUR_API_KEY",
		Prefix:   "vp",
	}

	// Get Figma data from API
	file, err := figma.GetData()
	if err != nil {
		fmt.Println("Error fetching Figma file:", err)
		return
	}

	// Get Figma Variables from API
	variables, err := figma.GetVariablesData()
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

	// Get CSS from tokens
	tokensCSS, err := figma.GenerateTokensCSS(tokens)
	if err != nil {
		panic(err)
	}

	// Parse components
	components := figma.ParseComponents(file, tokens)

	// Get CSS from components
	componentsCSS, err := figma.GenerateComponentsCSS(components)
	if err != nil {
		panic(err)
	}

	fmt.Printf("[tokensCSS] : %+v \n\n", tokensCSS)
	fmt.Printf("[componentsCSS] : %+v \n\n", componentsCSS)
}
