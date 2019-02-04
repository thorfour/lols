package main

import (
	"flag"
	"fmt"
	"net/url"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/thorfour/lols/pkg/lols"
	"github.com/thorfour/sillyputty/pkg/sillyputty"
)

var port = flag.Int("p", 80, "port to serve on")

func init() {
	flag.Parse()
}

func main() {
	logrus.Info("Starting lols server")

	// Sync internal cache on start
	if err := lols.Sync(); err != nil {
		logrus.Fatalf("Failed to sync: %v", err)
	}

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

	s.Port = *port
	s.Run()
}
