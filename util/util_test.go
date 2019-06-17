package util

import "testing"

func TestGenShortId(t *testing.T) {
	shortId, err := GenShortId()
	if shortId == "" || err == nil {
		t.Error("GenShortId failed!")
	}
	t.Log("GenShortId test pass")
}
