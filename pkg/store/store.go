package store

import (
	"fmt"
	"io"
	"net/http"

	"github.com/minio/minio-go"
)

var (
	spacesURL string
)

// Put takes an image url, downloads that image and uploads it to s3/spaces location and returns the new url string
func Put(url, newname string) (string, error) {
	// Download the image to s3/spaces
	loc, err := downloadImage(url, newname)
	if err != nil {
		return "", err
	}

	return loc, nil
}

// List returns a list of all the file names in the s3/spaces bucket
// This function is meant to be called once on startup, and the results cached
func List() ([]string, error) {
	client, err := minio.New("", "", "", false)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to datastore: ", err)
	}

	for o := range client.ListObjects("bucketname", "", false, nil) {
		fmt.Println(o)
	}

	return nil, nil
}

// downloadImage downloads an image from a given url and uploads it to s3/spaces
func downloadImage(url, newfilename string) (string, error) {
	loc := fmt.Sprintf("%s/%s", spacesURL, newfilename)

	// Download and re-upload the image in the background
	go func() {
		resp, err := http.Get(url)
		if err != nil {
			return
		}
		defer resp.Body.Close()

		storeImage(resp.Body, newfilename)
	}()

	return loc, nil
}

// storeImage copies data to a s3/spaces location
func storeImage(img io.Reader, name string) (string, error) {
	// Upload the image
	// TODO
	client, err := minio.New("", "", "", false)
	if err != nil {
		return "", fmt.Errorf("failed to connect to datastore: ", err)
	}
	return "", nil
}
