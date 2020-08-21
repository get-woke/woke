package main

import (
	"fmt"

	"github.com/caitlinelfring/woke/pkg/config"
)

func main() {
	c, _ := config.NewConfig("default.yaml")
	results, err := c.Parse("README.md")
	if err != nil {
		panic(err)
	}
	for _, result := range results {
		fmt.Println(result)
	}
}
