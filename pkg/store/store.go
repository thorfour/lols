package store

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/minio/minio-go"
)

const (
	version   = "v1"
	ssl       = true
	spacesURL = "" // TODO
)

var (
	endpoint  string
	accessKey string
	secretKey string
	name      string
)

func init() {
	accessKey = os.Getenv("SPACES_KEY")
	secretKey = os.Getenv("SPACES_SECRET")
	endpoint = os.Getenv("SPACES_ENDPOINT")
	name = "lols"
}

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
	client, err := minio.New(endpoint, accessKey, secretKey, ssl)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to datastore: %v", err)
	}

	var objects []string
	for o := range client.ListObjects(version, "", false, nil) {
		objects = append(objects, o.Key)
	}

	return objects, nil
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
	client, err := minio.New(endpoint, accessKey, secretKey, ssl)
	if err != nil {
		return "", fmt.Errorf("failed to connect to datastore: %v", err)
	}

	if _, err := client.PutObject(version, name, img, -1, minio.PutObjectOptions{}); err != nil {
		return "", fmt.Errorf("failed to store image: %v", err)
	}

	return name, nil
}
