package utils

import (
	"math/rand"
	"time"
)

var rander *rand.Rand

func init() {
	rander = rand.New(rand.NewSource(time.Now().Unix()))
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func RandString(l, r int) string {
	if l > r {
		l = r
	}
	length := rander.Intn(r-l+1) + l
	b := make([]rune, length)
	for i := range b {
		b[i] = letters[rander.Intn(len(letters))]
	}
	return string(b)
}
