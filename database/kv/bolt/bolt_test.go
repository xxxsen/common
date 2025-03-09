package bolt

import (
	"context"
	"fmt"
	"math/rand/v2"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/xxxsen/common/database/kv"
)

type testSt struct {
	A int    `json:"int"`
	B bool   `json:"b"`
	C string `json:"c"`
}

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

func TestTx(t *testing.T) {
	file := filepath.Join(os.TempDir(), uuid.NewString())
	defer os.RemoveAll(file)
	tab := "tmp_tab"
	db, err := New(file, tab)
	assert.NoError(t, err)
	defer db.Close()
	ctx := context.Background()
	go func() {
		time.Sleep(100 * time.Millisecond)
		start := time.Now()
		_, ok, err := db.Get(ctx, tab, "aaa")
		assert.NoError(t, err)
		assert.False(t, ok)
		t.Logf("thread cost:%dms", time.Since(start).Milliseconds())
	}()
	err = db.OnTranscation(ctx, func(ctx context.Context, db kv.IKvQueryExecutor) error {
		if err := db.Set(ctx, tab, "aaa", []byte("bbb")); err != nil {
			return err
		}
		res, ok, err := db.Get(ctx, tab, "aaa")
		if err != nil {
			return err
		}
		if !ok {
			return fmt.Errorf("not found in transation")
		}
		assert.Equal(t, []byte("bbb"), res)
		time.Sleep(500 * time.Millisecond)
		return nil
	})
	{
		res, ok, err := db.Get(ctx, tab, "aaa")
		assert.NoError(t, err)
		assert.True(t, ok)
		assert.Equal(t, []byte("bbb"), res)
	}
	assert.NoError(t, err)
}

func TestKvGetSetObj(t *testing.T) {
	st := &testSt{
		A: 1,
		B: true,
		C: "test",
	}
	file := filepath.Join(os.TempDir(), uuid.NewString())
	defer os.RemoveAll(file)
	tab := "tmp_tab"
	db, err := New(file, tab)
	assert.NoError(t, err)
	defer db.Close()
	ctx := context.Background()
	err = kv.MultiSetJsonObject(ctx, db, tab, map[string]*testSt{
		"aaa": st,
	})
	assert.NoError(t, err)
	m, err := kv.MultiGetJsonObject[testSt](ctx, db, tab, []string{"aaa"})
	assert.NoError(t, err)
	v, ok := m["aaa"]
	assert.True(t, ok)
	assert.Equal(t, st, v)
}

func TestSelectForUpdate(t *testing.T) {
	st := &testSt{
		A: 1,
		B: true,
		C: "test",
	}
	file := filepath.Join(os.TempDir(), uuid.NewString())
	defer os.RemoveAll(file)
	tab := "tmp_tab"
	db, err := New(file, tab)
	assert.NoError(t, err)
	defer db.Close()
	ctx := context.Background()
	err = kv.SetJsonObject(ctx, db, tab, "aaa", st)
	assert.NoError(t, err)
	err = kv.OnGetJsonKeyForUpdate[testSt](ctx, db, tab, "aaa", func(ctx context.Context, key string, val *testSt) (*testSt, bool, error) {
		val.C = "this is a test"
		return val, true, nil
	})
	assert.NoError(t, err)
	obj, ok, err := kv.GetJsonObject[testSt](ctx, db, tab, "aaa")
	assert.NoError(t, err)
	assert.True(t, ok)
	assert.Equal(t, "this is a test", obj.C)
}

func BenchmarkGet(b *testing.B) {
	file := filepath.Join(os.TempDir(), uuid.NewString())
	defer os.RemoveAll(file)
	tab := "tmp_tab"
	db, err := New(file, tab)
	if err != nil {
		b.Fatalf("failed to create db: %v", err)
	}
	defer db.Close()
	ctx := context.Background()
	err = db.Set(ctx, tab, "hello", []byte("world"))
	if err != nil {
		b.Fatalf("failed to set value: %v", err)
	}

	b.ResetTimer()
	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			_, _, err := db.Get(ctx, tab, "hello")
			if err != nil {
				b.Fatalf("failed to get value: %v", err)
			}
		}
	})

}

func BenchmarkSet(b *testing.B) {
	file := filepath.Join(os.TempDir(), uuid.NewString())
	defer os.RemoveAll(file)
	tab := "tmp_tab"
	db, err := New(file, tab)
	if err != nil {
		b.Fatalf("failed to create db: %v", err)
	}
	defer db.Close()
	ctx := context.Background()

	b.ResetTimer()
	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			err := db.Set(ctx, tab, strconv.FormatUint(rand.Uint64(), 10), []byte("value"))
			if err != nil {
				b.Fatalf("failed to set value: %v", err)
			}
		}
	})
}
