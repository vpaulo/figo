package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/vpaulo/figo"
	fg "github.com/vpaulo/figo/figma"
)

func main() {
	figma := figo.Figma{
		FILE_KEY: "",
		API_KEY:  "",
	}

	// file, err := figma.GetData()
	// if err != nil {
	// 	fmt.Println("Error fetching Figma file:", err)
	// 	return
	// }

	var file fg.File
	dat, err := os.ReadFile("./tmp/original_output.json")
	if err != nil {
		fmt.Println("Error Reading Figma file:", err)
		return
	}

	if unmarshallingError := json.Unmarshal(dat, &file); unmarshallingError != nil {
		fmt.Println("Error Unmarshal Figma file:", unmarshallingError)
		return
	}

	figma.Normalise(file)
}
