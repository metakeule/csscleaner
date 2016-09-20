package main

import (
	"fmt"
	"github.com/metakeule/config"
	"io/ioutil"
	"lib"
	"os"
)

var (
	cb      = lib.NewCodeBeautifier()
	cfg     = config.MustNew("csscleaner", "1.0.0", "cleans a given css file by using www.codebeautifier.com service")
	fileArg = cfg.NewString(
		"file",
		"file containing the css",
		config.Required,
	)
	fieldArg = cfg.NewString(
		"field",
		"overrides the form field where the css is entered on the page",
		config.Default(cb.Formfield),
	)
	urlArg = cfg.NewString(
		"url",
		"overrides the url to which the form is posted",
		config.Default(cb.PostURL),
	)
	elementArg = cfg.NewString(
		"element",
		"overrides the id of the element on the result page, containing the code",
		config.Default(cb.ResultElementID),
	)
	downloadArg = cfg.NewString(
		"download",
		"overrides the name of download checkbox",
		config.Default(cb.FileDownloadCheckbox),
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
			configure(&cb)
			result, err = cb.Cleanup(string(css))
		}
	}

	return
}

func configure(cb *lib.Config) {
	if fieldArg.IsSet() {
		cb.Formfield = fieldArg.Get()
	}
	if urlArg.IsSet() {
		cb.PostURL = urlArg.Get()
	}
	if elementArg.IsSet() {
		cb.ResultElementID = elementArg.Get()
	}
	if downloadArg.IsSet() {
		cb.FileDownloadCheckbox = downloadArg.Get()
	}
}
