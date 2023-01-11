package utils

import "testing"

func TestRandString(t *testing.T) {
	t.Logf("%s", RandString(5, 5))
	t.Logf("%s", RandString(5, 7))
	t.Logf("%s", RandString(0, 7))
	t.Logf("%s", RandString(2, 7))
	t.Logf("%s", RandString(3, 7))
}
