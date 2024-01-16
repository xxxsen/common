package versionmgr

import (
	"context"

	"github.com/xxxsen/common/database"
)

type VersionManager struct {
	c    *config
	vers []*VersionAction
}

func New(opts ...Option) (*VersionManager, error) {
	c := &config{
		tabName: defaultVersionManagementTab,
	}
	for _, opt := range opts {
		opt(c)
	}
	return &VersionManager{c: c}, nil
}

func (m *VersionManager) AddVersionAction(ver *VersionAction) {
	m.vers = append(m.vers, ver)
}

func (m *VersionManager) ProcessVersionAction(ctx context.Context, client database.IQueryExecer) error {
	if len(m.vers) == 0 {
		return nil
	}
	h := m.createHandler(ctx, client)
	return h.handle()
}

func (m *VersionManager) createHandler(ctx context.Context, client database.IQueryExecer) *versionActionHandler {
	return newVersionActionHandler(ctx, client, m.c.tabName, m.vers)
}
