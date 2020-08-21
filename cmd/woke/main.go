package main

import (
	"fmt"
	"os"

	"github.com/caitlinelfring/woke/pkg/config"
	"github.com/caitlinelfring/woke/pkg/parser"
)

func main() {
	c, _ := config.NewConfig("default.yaml")
	p := parser.Parser{Rules: c.Rules}
	results, err := p.Parse("README.md")
	if err != nil {
		panic(err)
	}
	for _, result := range results {
		fmt.Println(result)
	}
	if len(results) > 0 {
		os.Exit(1)
	}
}
