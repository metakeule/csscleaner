package lib

import (
	"fmt"
	"golang.org/x/net/html"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
)

// using the service http://www.codebeautifier.com

type Config struct {
	Formfield            string
	FileDownloadCheckbox string
	PostURL              string
	ResultElementID      string
}

type Option func(*Config)

func NewCodeBeautifier(options ...Option) Config {
	c := Config{
		Formfield:            "css_text",
		PostURL:              "http://www.codebeautifier.com",
		ResultElementID:      "result",
		FileDownloadCheckbox: "file_output",
	}

	for _, o := range options {
		o(&c)
	}

	return c
}

func (c Config) post(css string) (*http.Response, error) {
	resp, err := http.PostForm(
		c.PostURL,
		url.Values{
			c.Formfield:            {css},
			c.FileDownloadCheckbox: {c.FileDownloadCheckbox},
		},
	)

	if err != nil {
		return nil, err
	}

	switch resp.StatusCode {
	case http.StatusOK:
		return resp, nil
	case http.StatusMovedPermanently:
		resp.Body.Close()
		return http.Get(resp.Header.Get("Location"))
	case http.StatusTemporaryRedirect:
		resp.Body.Close()
		return http.Get(resp.Header.Get("Location"))
	default:
		resp.Body.Close()
		return nil, fmt.Errorf("wrong status code: %v", resp.StatusCode)
	}

}

var findCSSDownloadLink = regexp.MustCompile(`\.(css|CSS)$`)

func (c Config) getResultElement(parent *html.Node) (*html.Node, error) {
	node := getElement(parent, byID(c.ResultElementID))
	if node == nil {
		return nil, fmt.Errorf("could not find element with id %#v on result page", c.ResultElementID)
	}
	return node, nil
}

func (c Config) getDownloadLinkElement(parent *html.Node) (*html.Node, error) {
	node := getElement(parent, and(byTag("a"), byAttrMatch("href", findCSSDownloadLink)))
	if node == nil {
		return nil, fmt.Errorf("could not find download link")
	}
	return node, nil
}

func (c Config) getHrefAttribute(node *html.Node) (*html.Attribute, error) {
	href := getAttribute(node, "href")
	if href == nil {
		return nil, fmt.Errorf("could not get href of download link")
	}
	return href, nil
}

func (c Config) Cleanup(in string) (string, error) {
	var (
		err      error
		resp     *http.Response
		node     *html.Node
		href     *html.Attribute
		baseURL  *url.URL
		resBytes []byte
	)

steps:
	for jump := 1; err == nil; jump++ {
		switch jump - 1 {
		default:
			break steps
		// count a number up for each following step
		case 0:
			resp, err = c.post(in)
		case 1:
			node, err = html.Parse(resp.Body)
			resp.Body.Close()
		case 2:
			node, err = c.getResultElement(node)
		case 3:
			node, err = c.getDownloadLinkElement(node)
		case 4:
			href, err = c.getHrefAttribute(node)
		case 5:
			baseURL, err = url.Parse(c.PostURL)
		case 6:
			var newURL url.URL
			newURL.Scheme = baseURL.Scheme
			newURL.Host = baseURL.Host
			newURL.Path = href.Val
			resp, err = http.Get(newURL.String())
		case 7:
			resBytes, err = ioutil.ReadAll(resp.Body)
			resp.Body.Close()
		}
	}

	return string(resBytes), err
}
