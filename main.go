package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/hatsu38/postman-tui/ui"
)

func usage() {
	format := `
	                 _                               _         _
	_ __   ___  ___| |_ _ __ ___   __ _ _ __       | |_ _   _(_)
       | '_ \ / _ \/ __| __| '_ ' _ \ / _' | '_ \ _____| __| | | | |
       | |_) | (_) \__ \ |_| | | | | | (_| | | | |_____| |_| |_| | |
       | .__/ \___/|___/\__|_| |_| |_|\__,_|_| |_|      \__|\__,_|_|
       |_|

	Usage:
	  postman-tui
	Author:
	  hatsu38<hajiwata0308@gmail.com>
	`
	fmt.Fprintln(os.Stderr, fmt.Sprintf(format))
}
func run() int {
	var i interface{}
	if err := ui.New().Run(i); err != nil {
		log.Println(err)
		return 1
	}
	return 0
}

func main() {
	flag.Usage = usage
	flag.Parse()

	run()
}
