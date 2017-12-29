package main

import (
	"fmt"
	"os"

	"github.com/detectify/go-test-formatter/formatters"
	_ "github.com/detectify/go-test-formatter/formatters/teamcity"
	"github.com/detectify/go-test-formatter/parser"
)

func main() {
	parser := parser.New()

	suites, parserErr := parser.Parse(os.Stdin)

	if parserErr != nil {
		fmt.Println(parserErr.Error())
	}
	if suites == nil {
		os.Exit(1)
	}

	formatter, err := formatters.Find("teamcity")

	if err != nil {
		fmt.Println(err.Error())

		os.Exit(3)
	}

	err = formatter.Format(suites, os.Stdout)

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(4)
	}

	if parserErr != nil {
		os.Exit(1)
	}

	for _, suite := range suites {
		for _, test := range suite.Tests {
			if test.Failed {
				os.Exit(2)
			}
		}
	}
}
