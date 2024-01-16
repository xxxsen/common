package versionmgr

import (
	"context"
	"database/sql"
	"os"
	"testing"

	_ "github.com/mattn/go-sqlite3"

	"github.com/stretchr/testify/assert"
)

func mustCreateSqliteForTest() *sql.DB {
	dbPath := os.TempDir() + "/test.db"
	_ = os.Remove(dbPath)
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		panic(err)
	}
	return db
}

func TestVersionManager(t *testing.T) {
	mgr, err := New()
	assert.NoError(t, err)
	mgr.AddVersionAction(&VersionAction{
		OnInit: SQLAction{Steps: []SQLActionStep{
			{
				Name: "create table",
				Action: SimpleSQLExecAction(`
						create table if not exists abc (
							a int unsigned not null,
							b int unsigned not null,
							primary key(a)
						)
					`),
			},
		}},
		OnUpgrade:   SQLAction{},
		Version:     20210101,
		LastVersion: 0,
	})
	mgr.AddVersionAction(&VersionAction{
		OnInit: SQLAction{Steps: []SQLActionStep{
			{
				Name: "create table",
				Action: SimpleSQLExecAction(`
						create table if not exists abc (
							a int unsigned not null,
							b int unsigned not null,
							c varchar(256) not null default '',
							primary key(a)
						)
					`),
			},
		}},
		OnUpgrade: SQLAction{
			Steps: []SQLActionStep{
				{
					Name: "create c field",
					Action: SimpleSQLExecAction(`
						alter table abc add column c varchar(256) not null default ''
					`),
				},
			},
		},
		Version:     20210102,
		LastVersion: 20210101,
	})
	db := mustCreateSqliteForTest()
	h := mgr.createHandler(context.Background(), db)
	err = h.ensureVersionValid()
	assert.NoError(t, err)
	err = h.ensureVersionTab()
	assert.NoError(t, err)
	acts, isCreate, err := h.findActionToExec()
	assert.NoError(t, err)
	assert.True(t, isCreate)
	assert.True(t, len(acts) == 1)
	act := acts[0]
	assert.Equal(t, uint64(20210101), act.LastVersion)
	assert.Equal(t, uint64(20210102), act.Version)
	err = h.performStepAction(act.OnInit.Steps)
	assert.NoError(t, err)
	//retry should not return err
	err = h.performStepAction(act.OnInit.Steps)
	assert.NoError(t, err)
	err = mgr.ProcessVersionAction(context.Background(), db)
	assert.NoError(t, err)
}

func TestVersionManagerUpdate(t *testing.T) {
	mgr, err := New()
	assert.NoError(t, err)
	mgr.AddVersionAction(&VersionAction{
		OnInit: SQLAction{Steps: []SQLActionStep{
			{
				Name: "create table",
				Action: SimpleSQLExecAction(`
						create table if not exists abc (
							a int unsigned not null,
							b int unsigned not null,
							primary key(a)
						)
					`),
			},
		}},
		OnUpgrade:   SQLAction{},
		Version:     20210101,
		LastVersion: 0,
	})
	mgr.AddVersionAction(&VersionAction{
		OnInit: SQLAction{Steps: []SQLActionStep{
			{
				Name: "create table",
				Action: SimpleSQLExecAction(`
						create table if not exists abc (
							a int unsigned not null,
							b int unsigned not null,
							c varchar(256) not null default '',
							primary key(a)
						)
					`),
			},
		}},
		OnUpgrade: SQLAction{
			Steps: []SQLActionStep{
				{
					Name: "create c field",
					Action: SimpleSQLExecAction(`
						alter table abc add column c varchar(256) not null default ''
					`),
				},
			},
		},
		Version:     20210102,
		LastVersion: 20210101,
	})

	db := mustCreateSqliteForTest()
	err = mgr.ProcessVersionAction(context.Background(), db)
	assert.NoError(t, err)
	mgr.AddVersionAction(&VersionAction{
		OnInit: SQLAction{Steps: []SQLActionStep{
			{
				Name: "create table",
				Action: SimpleSQLExecAction(`
						create table if not exists abc (
							a int unsigned not null,
							b int unsigned not null,
							c varchar(256) not null default '',
							d varchar(512) not null default 'abc',
							primary key(a)
						)
					`),
			},
		}},
		OnUpgrade: SQLAction{
			Steps: []SQLActionStep{
				{
					Name: "create c field",
					Action: SimpleSQLExecAction(`
						alter table abc add column d varchar(512) not null default 'abc'
					`),
				},
			},
		},
		Version:     20210103,
		LastVersion: 20210102,
	})
	err = mgr.ProcessVersionAction(context.Background(), db)
	assert.NoError(t, err)
}
