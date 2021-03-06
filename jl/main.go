package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/huangjunwen/jsonlogic-go"
	"github.com/huangjunwen/jsonlogic-go/ext"
)

func outputErrorAndExit(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}

func init() {
	// Add extensions.
	ext.AddOpRange(jsonlogic.DefaultJSONLogic)
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

	result, err := jsonlogic.Apply(logic, data)
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
