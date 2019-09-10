package lols

import (
	"fmt"
	"path/filepath"
	"strings"
	"sync"

	"github.com/thorfour/lols/pkg/store"
)

var commands = map[string]func([]string) (string, error){
	"new":  newLol,
	"add":  newLol,
	"help": helpLol,
}

var (
	// internal cache of image names
	imageNames []string
	imagesLock *sync.RWMutex
)

func init() {
	imagesLock = new(sync.RWMutex)
}

// Handle responds to given commands
func Handle(cmd string, args []string) (string, error) {
	f, ok := commands[cmd]
	if !ok { // if the command isn't found, treat it as a descriptor
		return getLol(append(args, cmd))
	}

	return f(args)
}

// Sync the internal cache of imageNames, should be run in a separate go routine
func Sync() error {
	out, err := store.List()
	if err != nil {
		return fmt.Errorf("failed to list: %v", err)
	}

	imagesLock.Lock()
	defer imagesLock.Unlock()
	imageNames = make([]string, len(imageNames))
	for n := range out {
		imageNames = append(imageNames, n)
	}

	return nil
}

func newLol(args []string) (string, error) {
	if len(args) != 2 {
		return "", fmt.Errorf("invalid number of arguments")
	}

	img := args[0]
	newname := args[1]

	// TODO async handle these errors
	loc, _ := store.Put(img, newname)

	// Add the image to the cache
	imagesLock.Lock()
	defer imagesLock.Unlock()
	imageNames = append(imageNames, loc)

	// wrap it as an error to be returned as an ephemeral message
	return "", fmt.Errorf("Added")
}

type result struct {
	name  string
	match int
}

// getLol returns the url for the best matching image
func getLol(args []string) (string, error) {

	// Look through all the images for a matching filename
	imagesLock.RLock()
	defer imagesLock.RUnlock()
	results := make(chan result, len(imageNames))
	wg := new(sync.WaitGroup)
	wg.Add(len(imageNames))
	go func() { // Close results after all have been received
		wg.Wait()
		close(results)
	}()
	for _, n := range imageNames {
		go match(n, args, results, wg)
	}

	// Check all the results
	max := 0
	bestMatch := ""
	for r := range results {
		if r.match > max {
			bestMatch = r.name
			max = r.match
		}
		// Pick the shortest matching file name
		if r.match == max && len(r.name) < len(bestMatch) {
			bestMatch = r.name
		}
	}

	if bestMatch == "" {
		return "", fmt.Errorf("That's not an emergency")
	}

	return bestMatch, nil
}

// match checks strings against a given name, and send the number of those substrings that are found in the name into the result channel
func match(n string, args []string, results chan<- result, wg *sync.WaitGroup) {
	defer wg.Done()

	count := 0
	name := filepath.Base(n)
	for _, a := range args {
		if strings.Contains(strings.ToLower(name), strings.ToLower(a)) {
			count++
		}
	}

	results <- result{name: n, match: count}
}

func helpLol(_ []string) (string, error) {
	return "", fmt.Errorf("'new <image_url> <image name>' to store a new lolmergency")
}
