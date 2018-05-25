package store

import (
	"testing"
)

func TestDownload(t *testing.T) {

	testcorgi := "https://www.pets4homes.co.uk/images/breeds/50/d248d59954bb644e4437cce1758a9ce2.jpg"
	newImage := "testcorgi.jpg"

	_, errCh := Put(testcorgi, newImage)
	for err := range errCh {
		if err != nil {
			t.Fatalf("failed to upload: %v", err)
		}
	}
}
