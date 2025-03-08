package bolt

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetSet(t *testing.T) {
	file := filepath.Join(os.TempDir(), uuid.NewString())
	defer os.RemoveAll(file)
	tab := "tmp_tab"
	db, err := New(file, tab)
	assert.NoError(t, err)
	defer db.Close()
	ctx := context.Background()
	err = db.Set(ctx, tab, "hello", []byte("world"))
	assert.NoError(t, err)
	val, ok, err := db.Get(ctx, tab, "hello")
	assert.NoError(t, err)
	assert.True(t, ok)
	assert.Equal(t, []byte("world"), val)
	_, ok, err = db.Get(ctx, tab, "test_not_found")
	assert.NoError(t, err)
	assert.False(t, ok)
	//bucket not found
	_, _, err = db.Get(ctx, "test", "111")
	assert.Error(t, err)
}

func TestIter(t *testing.T) {
	file := filepath.Join(os.TempDir(), uuid.NewString())
	defer os.RemoveAll(file)
	tab := "tmp_tab"
	db, err := New(file, tab)
	assert.NoError(t, err)
	defer db.Close()
	ctx := context.Background()

	err = db.MultiSet(ctx, tab, map[string][]byte{
		"k-a": []byte("v-a"),
		"k-b": []byte("v-b"),
		"k-c": []byte("v-c"),
		"k-d": []byte("v-d"),
		"k-e": []byte("v-e"),
		"k-f": []byte("v-f"),
		"k-g": []byte("v-g"),
	})
	assert.NoError(t, err)
	rs, err := db.MultiGet(ctx, tab, []string{"k-a", "k-b", "k-c", "k-d", "k-e", "k-f", "k-g"})
	assert.NoError(t, err)
	assert.Equal(t, []byte("v-a"), rs["k-a"])
	assert.Equal(t, []byte("v-b"), rs["k-b"])
	assert.Equal(t, []byte("v-c"), rs["k-c"])
	assert.Equal(t, []byte("v-d"), rs["k-d"])
	assert.Equal(t, []byte("v-e"), rs["k-e"])
	assert.Equal(t, []byte("v-f"), rs["k-f"])
	assert.Equal(t, []byte("v-g"), rs["k-g"])
	err = db.Iter(ctx, tab, "k-", func(ctx context.Context, key string, value []byte) (bool, error) {
		t.Logf("read k:%s, v:%s", key, string(value))
		return true, nil
	})
	assert.NoError(t, err)
}
