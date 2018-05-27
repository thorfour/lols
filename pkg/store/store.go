package store

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

const (
	version = "v0.0.0"
	ssl     = true
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
func List() ([]string, error) {
	s := session.New(s3Config)
	client := s3.New(s)

	var objects []string
	bucket := version
	o, err := client.ListObjects(&s3.ListObjectsInput{
		Bucket: &bucket,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list objects: %v", err)
	}

	for _, r := range o.Contents {
		objects = append(objects, *r.Key)
	}

	return objects, nil
}

// downloadImage downloads an image from a given url and uploads it to s3/spaces
func downloadImage(url, newfilename string) (string, <-chan error) {
	loc := fmt.Sprintf("%s/%s/%s", endpoint, version, newfilename)

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

		if err := storeImage(resp.Body, resp.ContentLength, newfilename); err != nil {
			errCh <- err
			return
		}
	}()

	return loc, errCh
}

// storeImage copies data to a s3/spaces location
func storeImage(img io.Reader, size int64, name string) error {
	s := session.New(s3Config)
	uploader := s3manager.NewUploader(s)

	acl := "public-read"
	imgtype := "image/jpeg"
	if _, err := uploader.Upload(&s3manager.UploadInput{
		ACL:         &acl,
		Bucket:      &bucket,
		Key:         &name,
		ContentType: &imgtype,
		Body:        img,
	}); err != nil {
		return fmt.Errorf("failed to store image: %v", err)
	}

	return nil
}
