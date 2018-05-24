package store

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/minio/minio-go"
)

// Put takes an image url, downloads that image and uploads it to s3/spaces location and returns the new url string
func Put(url, newname string) (string, error) {
	// Download the image
	loc, err := downloadImage(url, newname)
	defer deleteImage(loc)
	if err != nil {
		return "", err
	}

	// Upload the image
	// TODO
	client, err := minio.New("", "", "", "")
	if err != nil {
		return "", fmt.Errorf("failed to connect to datastore: ", err)
	}

	return "", nil
}

// List returns a list of all the file names in the s3/spaces bucket
// This function is meant to be called once on startup, and the results cached
func List() ([]string, error) {
	client, err := minio.New("", "", "", "")
	if err != nil {
		return nil, fmt.Errorf("failed to connect to datastore: ", err)
	}

	for o := range client.ListObjects("bucketname", "", false, nil) {
	}

	return nil, nil
}

// downloadImage downloads an image from a given url with and locally stores it with given newfilename
// returns the location of the new file
func downloadImage(url, newfilename string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to download image: %v", err)
	}
	defer resp.Body.Close()

	f, err := ioutil.TempFile("", newfilename)
	if err != nil {
		return "", fmt.Errorf("failed to create local copy: %v", err)
	}
	defer f.Close()

	if _, err = io.Copy(f, resp.Body); err != nil {
		return f.Name(), fmt.Errorf("copy failed: %v", err)
	}

	return f.Name(), nil
}

// deleteImage deletes a local image from disk
func deleteImage(loc string) {
}
