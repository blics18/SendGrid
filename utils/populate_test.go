package utils

import "testing"

func TestRandStringRunes(t *testing.T) {
	res := RandStringRunes(10)
	if len(res) != 10 {
		t.Error("Should be equal to 10")
	}
}
