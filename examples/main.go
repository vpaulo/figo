package main

import (
	"fmt"

	"github.com/vpaulo/figo"
)

func main() {
	figma := figo.Figma{
		FILE_KEY: "FIGMA_FILE",
		API_KEY:  "YOUR_API_KEY",
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

	// TODO: GetVariables

	tokens := figma.ParseTokens(file)

	// TODO: ParseVariablesTokens

	tokensCSS, err := figma.GenerateTokensCSS(tokens)
	if err != nil {
		panic(err)
	}
	fmt.Printf("[tokensCSS] : %+v \n\n", tokensCSS)
}
