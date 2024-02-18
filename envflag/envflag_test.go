package envflag

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	os.Setenv("A", "1")
	os.Setenv("B", "2.5")
	os.Setenv("C", "true")
	os.Setenv("D", "hello")
	var a int
	var b float64
	var c bool
	f := DefaultParser
	f.IntVar(&a, "a", 0, "")
	f.Float64Var(&b, "b", 0, "")
	f.BoolVar(&c, "c", true, "")
	d := f.String("d", "", "")
	f.Parse()
	assert.Equal(t, 1, a)
	assert.Equal(t, 2.5, b)
	assert.Equal(t, true, c)
	assert.Equal(t, "hello", *d)
}
