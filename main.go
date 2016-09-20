package main

import (
	"fmt"
	"github.com/metakeule/config"
	"io/ioutil"
	"lib"
	"os"
)

var (
	cfg     = config.MustNew("csscleaner", "2.0.0", "cleans a given css file by using csstidy")
	fileArg = cfg.NewString(
		"file",
		"file containing the css",
		config.Required,
	)
)

func main() {
	res, err := run()

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
		os.Exit(1)
	}

	fmt.Fprintln(os.Stdout, res)
	os.Exit(0)
}

func run() (result string, err error) {
	var (
		css []byte
	)

steps:
	for jump := 1; err == nil; jump++ {
		switch jump - 1 {
		default:
			break steps
		// count a number up for each following step
		case 0:
			err = cfg.Run()
		case 1:
			css, err = ioutil.ReadFile(fileArg.Get())
		case 2:
			result, err = lib.Cleanup(string(css))
		}
	}

	return
}
