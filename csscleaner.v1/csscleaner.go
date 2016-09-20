package csscleaner

import (
	"lib"
)

// Formfield overrides the default form field
func Formfield(field string) lib.Option {
	return func(c *lib.Config) {
		c.Formfield = field
	}
}

// FileDownloadCheckbox overrides the default form checkbox to activate file download
func FileDownloadCheckbox(field string) lib.Option {
	return func(c *lib.Config) {
		c.FileDownloadCheckbox = field
	}
}

// PostURL overrides the default post url
func PostURL(url string) lib.Option {
	return func(c *lib.Config) {
		c.PostURL = url
	}
}

// ResultElementID overrides the default result element id
func ResultElementID(id string) lib.Option {
	return func(c *lib.Config) {
		c.ResultElementID = id
	}
}

// Cleanup cleans up the given css via the service www.codebeautifier.com
func Cleanup(css string, options ...lib.Option) (string, error) {
	return lib.NewCodeBeautifier(options...).Cleanup(css)
}
