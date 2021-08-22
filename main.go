package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/hatsu38/postman-tui/ui"
)

var (
	Version  string
	Revision string
)

func usage() {
	format :=
`		_                                _         _
_ __   ___  ___| |_ _ __ ___   __ _ _ __        | |_ _   _(_)
| '_ \ / _ \/ __| __| '_ ' _ \ / _' | '_ \ _____| __| | | | |
| |_) | (_) \__ \ |_| | | | | | (_| | | | |_____| |_| |_| | |
| .__/ \___/|___/\__|_| |_| |_|\__,_|_| |_|      \__|\__,_|_|
|_|  Version: %s-%s

Usage:
  postman-tui

  Flags:
    -version: Print the current version

Author:
  hatsu38<hajiwata0308@gmail.com>
`
	fmt.Fprintln(os.Stderr, fmt.Sprintf(format, Version, Revision))
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
	var withVersion bool
	flag.BoolVar(&withVersion, "version", false, "Print the current version")

	flag.Usage = usage
	flag.Parse()

	if withVersion {
		fmt.Println("Version ", Version)
		return
	}

	run()
}
