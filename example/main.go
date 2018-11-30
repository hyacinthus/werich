package main

import (
	"fmt"
	"io/ioutil"

	"github.com/hyacinthus/werich"
)

func main() {
	input, err := ioutil.ReadFile("../test.md")
	if err != nil {
		panic(err)
	}
	md := werich.NewMD(input)
	// meta
	var meta = make(map[string]interface{})
	err = md.Meta(meta)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%v", meta)
	// rich
	rich := md.Rich()
	fmt.Println(string(rich))
}
