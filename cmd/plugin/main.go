package main

import (
	"net/url"
	"strings"

	"github.com/thorfour/lols/pkg/lols"
)

func init() {
	lols.Sync() // sync the initial internal cache
}

// Handler is the plugin handler
func Handler(v url.Values) (string, error) {
	text := strings.Split(v["text"][0], " ")
	var args []string
	if len(text) > 1 {
		args = text[1:]
	}
	return lols.Handle(text[0], args)
}
