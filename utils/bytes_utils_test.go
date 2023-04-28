package utils

import (
	"encoding/hex"
	"log"
	"testing"
)

func TestRandBytes(t *testing.T) {
	data := RandBytes(10, 20)
	log.Printf("%s", hex.EncodeToString(data))
}
