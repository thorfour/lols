package main

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/thorfour/lols/pkg/lols"
	"github.com/thorfour/sillyputty/pkg/sillyputty"
)

func main() {
	logrus.Info("Starting lols server")

	s := sillyputty.New("/v1",
		sillyputty.HandlerOpt("/lols", func(v url.Values) (string, error) {
			if v == nil {
				return "", fmt.Errorf("not enough arguments")
			}

			text := strings.Split(v["text"][0], " ")
			var args []string
			if len(text) > 1 { // remove first arg as the command
				args = text[1:]
			}

			return lols.Handle(text[0], args)
		}),
	)

	s.Run()
}
