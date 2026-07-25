package sqlite

import (
	"context"
	"fmt"
	"path/filepath"
	"testing"

	"github.com/didi/gendry/builder"
	"github.com/stretchr/testify/assert"
	"github.com/xxxsen/common/database"
	"github.com/xxxsen/common/database/dbkit"
)

const table = "test_tab"

func newTestDatabase(t *testing.T) database.IDatabase {
	t.Helper()
	ctx := t.Context()
	db, err := New(filepath.Join(t.TempDir(), "sqlite.db"), func(db database.IDatabase) error {
		if _, err := db.ExecContext(ctx, "create table test_tab(name text, age integer);"); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_ = db.Close()
	})
	return db
}

type testSt struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func TestReadWrite(t *testing.T) {
	db := newTestDatabase(t)
	data := []map[string]interface{}{
		{
			"name": "hello",
			"age":  123,
		},
		{
			"name": "world",
			"age":  321,
		},
	}
	sql, args, err := builder.BuildInsert(table, data)
	assert.NoError(t, err)
	ctx := context.Background()
	_, err = db.ExecContext(ctx, sql, args...)
	assert.NoError(t, err)

	where := map[string]interface{}{
		"_limit": []uint{0, 10},
	}
	rs := make([]*testSt, 0, 10)
	err = dbkit.SimpleQuery(ctx, db, table, where, &rs, dbkit.ScanWithTagName("json"))
	assert.NoError(t, err)
	assert.Equal(t, 2, len(rs))
	assert.Equal(t, "hello", rs[0].Name)
	assert.Equal(t, 123, rs[0].Age)
	assert.Equal(t, "world", rs[1].Name)
	assert.Equal(t, 321, rs[1].Age)
	t.Logf("data:%+v", rs)
}

func TestTxNoCommit(t *testing.T) {
	db := newTestDatabase(t)
	ctx := context.Background()
	err := db.OnTransation(ctx, func(ctx context.Context, qe database.IQueryExecer) error {
		data := []map[string]interface{}{
			{
				"name": "hello",
				"age":  123,
			},
			{
				"name": "world",
				"age":  321,
			},
		}
		sql, args, err := builder.BuildInsert(table, data)
		assert.NoError(t, err)
		_, err = qe.ExecContext(ctx, sql, args...)
		assert.NoError(t, err)
		return fmt.Errorf("skip commit")
	})
	assert.Error(t, err)

	where := map[string]interface{}{
		"_limit": []uint{0, 10},
	}
	rs := make([]*testSt, 0, 10)
	err = dbkit.SimpleQuery(ctx, db, table, where, &rs, dbkit.ScanWithTagName("json"))
	assert.NoError(t, err)
	assert.Equal(t, 0, len(rs))
}

func TestTxCommit(t *testing.T) {
	db := newTestDatabase(t)
	ctx := context.Background()
	err := db.OnTransation(ctx, func(ctx context.Context, qe database.IQueryExecer) error {
		data := []map[string]interface{}{
			{
				"name": "hello",
				"age":  123,
			},
			{
				"name": "world",
				"age":  321,
			},
		}
		sql, args, err := builder.BuildInsert(table, data)
		assert.NoError(t, err)
		_, err = qe.ExecContext(ctx, sql, args...)
		assert.NoError(t, err)
		return nil
	})
	assert.NoError(t, err)

	where := map[string]interface{}{
		"_limit": []uint{0, 10},
	}
	rs := make([]*testSt, 0, 10)
	err = dbkit.SimpleQuery(ctx, db, table, where, &rs, dbkit.ScanWithTagName("json"))
	assert.NoError(t, err)
	assert.Equal(t, 2, len(rs))
}
