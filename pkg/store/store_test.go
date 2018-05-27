package store

import (
	"testing"
)

func TestDownload(t *testing.T) {

	testcorgi := "http://lolmergency.com/dis_gon_be_good.gif"
	newImage := "dis.jpg"

	_, errCh := Put(testcorgi, newImage)
	for err := range errCh {
		if err != nil {
			t.Fatalf("failed to upload: %v", err)
		}
	}
}
