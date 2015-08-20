package main

import (
	"fmt"
	"github.com/detectify/go-test-formatter/formatters"
	_ "github.com/detectify/go-test-formatter/formatters/teamcity"
	"github.com/detectify/go-test-formatter/parser"
	"os"
)

func main() {
	parser := parser.New()

	suites, err := parser.Parse(os.Stdin)

	if err != nil {
		fmt.Println(err.Error())

		return
	}

	formatter, err := formatters.Find("teamcity")

	if err != nil {
		fmt.Println(err.Error())

		return
	}

	err = formatter.Format(suites, os.Stdout)

	if err != nil {
		fmt.Println(err.Error())
	}
}
