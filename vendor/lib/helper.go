package lib

import (
	"golang.org/x/net/html"
	"regexp"
	"strings"
)

func getAttribute(n *html.Node, attrname string) *html.Attribute {

	for _, attr := range n.Attr {
		if strings.ToLower(attr.Key) == strings.ToLower(attrname) {
			return &attr
		}
	}

	return nil
}

func matchAttribute(n *html.Node, attrname string, matcher *regexp.Regexp) bool {
	attr := getAttribute(n, attrname)
	if attr == nil {
		return false
	}

	return matcher.MatchString(attr.Val)
}

type by func(*html.Node) bool

func byID(id string) by {
	return func(n *html.Node) bool {
		if n.Type != html.ElementNode {
			return false
		}
		pid := getAttribute(n, "id")
		return pid != nil && pid.Val == id
	}
}

func byTag(tag string) by {
	return func(n *html.Node) bool {
		if n.Type != html.ElementNode {
			return false
		}
		return strings.ToLower(n.Data) == strings.ToLower(tag)
	}
}

func byAttrMatch(attr string, matcher *regexp.Regexp) by {
	return func(n *html.Node) bool {
		if n.Type != html.ElementNode {
			return false
		}
		return matchAttribute(n, attr, matcher)
	}
}

func and(bs ...by) by {
	return func(n *html.Node) bool {
		for _, b := range bs {
			if !b(n) {
				return false
			}
		}
		return true
	}
}

func getElement(parent *html.Node, b by) *html.Node {
	if b(parent) {
		return parent
	}

	for c := parent.FirstChild; c != nil; c = c.NextSibling {
		if n := getElement(c, b); n != nil {
			return n
		}
	}

	return nil
}
