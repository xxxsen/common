package tls

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDedup(t *testing.T) {
	lst := []string{"a", "b", "a"}
	lst2 := dedup(lst)
	sort.Strings(lst2)
	assert.Equal(t, 2, len(lst2))
	assert.Equal(t, "a", lst2[0])
	assert.Equal(t, "b", lst2[1])
}
