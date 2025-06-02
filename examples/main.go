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
	// file, err := figma.GetData()
	// if err != nil {
	// 	fmt.Println("Error fetching Figma file:", err)
	// 	return
	// }

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

	// tokensCSS, err := figma.GenerateTokensCSS(tokens)
	// if err != nil {
	// 	panic(err)
	// }

	// css := figma.ParseCSS(file, tokens)
	components := figma.ParseComponents(file, tokens)
	fmt.Printf("[COMPONENTS] : %+v \n\n", components)
	// fmt.Printf("[tokensCSS] : %+v \n\n", tokensCSS)
	// fmt.Printf("[CSS] : %+v \n\n", css)
}
