package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/ckotzbauer/chekr/tools/api-lifecycle-gen/internal"
)

func main() {
	argsWithoutProg := os.Args[1:]

	if len(argsWithoutProg) != 2 {
		panic("Pass single folder where k8s.io/api is located and the resulting json file-location")
	}

	groups, err := internal.ListFiles(argsWithoutProg[0])

	if err != nil {
		panic(err)
	}

	deprecatedTypes, err := internal.ParseGroups(groups)

	if err != nil {
		panic(err)
	}

	buf, err := json.MarshalIndent(deprecatedTypes, "", "  ")

	if err != nil {
		panic(err)
	}

	err = internal.WriteFile(argsWithoutProg[1], string(buf))

	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully parsed")
}
