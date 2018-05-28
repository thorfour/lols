package store

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

const (
	ssl = true
)

var (
	endpoint string
	bucket   string
	s3Config *aws.Config
)

func init() {
	accessKey := os.Getenv("SPACES_KEY")
	secretKey := os.Getenv("SPACES_SECRET")
	endpoint = os.Getenv("SPACES_ENDPOINT")
	bucket = os.Getenv("SPACES_BUCKET")

	s3Config = &aws.Config{
		Credentials: credentials.NewStaticCredentials(accessKey, secretKey, ""),
		Endpoint:    aws.String(endpoint),
		Region:      aws.String("us-east-1"),
	}
}

// Put takes an image url, downloads that image and uploads it to s3/spaces location and returns the new url string
func Put(url, newname string) (string, <-chan error) {
	// Download the image to s3/spaces
	return downloadImage(url, newname)
}

// List returns a list of all the file names in the s3/spaces bucket
// This function is meant to be called once on startup, and the results cached
func List() (<-chan string, error) {
	s := session.New(s3Config)
	client := s3.New(s)

	o, err := client.ListObjects(&s3.ListObjectsInput{
		Bucket: &bucket,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list objects: %v", err)
	}

	out := make(chan string, len(o.Contents))
	go func() {
		defer close(out)
		for _, r := range o.Contents {
			out <- *r.Key
		}
	}()

	return out, nil
}

// downloadImage downloads an image from a given url and uploads it to s3/spaces
func downloadImage(url, newfilename string) (string, <-chan error) {
	loc := getNewFilePath(url, newfilename)

	// Download and re-upload the image in the background
	errCh := make(chan error, 1)
	go func() {
		defer close(errCh)
		resp, err := http.Get(url)
		if err != nil {
			errCh <- fmt.Errorf("failed to download image: %v", err)
			return
		}
		defer resp.Body.Close()

		contentType := resp.Header.Get("Content-Type")
		if err := storeImage(resp.Body, resp.ContentLength, contentType, filepath.Base(loc)); err != nil {
			errCh <- err
			return
		}
	}()

	return loc, errCh
}

// storeImage copies data to a s3/spaces location
func storeImage(img io.Reader, size int64, contentType string, name string) error {
	s := session.New(s3Config)
	uploader := s3manager.NewUploader(s)

	acl := "public-read"
	if contentType == "" {
		contentType = "image/jpeg"
	}
	if _, err := uploader.Upload(&s3manager.UploadInput{
		ACL:         &acl,
		Bucket:      &bucket,
		Key:         &name,
		ContentType: &contentType,
		Body:        img,
	}); err != nil {
		return fmt.Errorf("failed to store image: %v", err)
	}

	return nil
}

// getNewFilePath returns the new file location for the given old and new name
func getNewFilePath(url, newfilename string) string {
	ext := filepath.Ext(url)
	if filepath.Ext(newfilename) == "" { // Copy the same file extension over
		newfilename = fmt.Sprintf("%s%s", newfilename, ext)
	}
	loc := fmt.Sprintf("%s.%s/%s", bucket, endpoint, newfilename)

	return loc
}
