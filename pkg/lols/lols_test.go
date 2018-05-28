package lols

import "testing"

func TestGetLol(t *testing.T) {

	imageNames = []string{"success.jpg", "success_troll.jpg"}
	n, err := getLol([]string{"succ"})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if n != "success.jpg" {
		t.Errorf("unexpected %s", n)
	}

	n, err = getLol([]string{"succ", "tro"})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if n != "success_troll.jpg" {
		t.Errorf("unexpected %s", n)
	}
}
