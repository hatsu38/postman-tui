package main

import (
	"log"

	"github.com/hatsu38/postman-tui/ui"
)

func run() int {
	var i interface{}
	if err := ui.New().Run(i); err != nil {
		log.Println(err)
		return 1
	}
	return 0
}
func main() {
	run()
}
