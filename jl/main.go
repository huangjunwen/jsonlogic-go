package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/huangjunwen/jsonlogic-go"
)

func outputErrorAndExit(err error) {
	fmt.Fprint(os.Stderr, err)
	os.Exit(1)
}

func main() {
	var (
		logic, data interface{}
	)
	decoder := json.NewDecoder(os.Stdin)
	if err := decoder.Decode(&logic); err != nil {
		outputErrorAndExit(err)
		return
	}
	if err := decoder.Decode(&data); err != nil {
		if err != io.EOF {
			outputErrorAndExit(err)
			return
		}
	}

	result, err := jsonlogic.New().Apply(logic, data)
	if err != nil {
		outputErrorAndExit(err)
		return
	}

	b, err := json.Marshal(result)
	if err != nil {
		outputErrorAndExit(err)
		return
	}

	fmt.Printf("%s\n", b)
}
