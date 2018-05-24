package main

import (
	"net/url"
	"strings"

	"github.com/thorfour/lols/pkg/store"
)

// Handler is the plugin handler
func Handler(v url.Values) (string, error) {
	text := strings.Split(v["text"][0], " ")
	return store.Handle(text)
}
