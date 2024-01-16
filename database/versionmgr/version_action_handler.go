package versionmgr

import (
	"context"
	"fmt"
	"sort"

	"github.com/didi/gendry/builder"
	"github.com/didi/gendry/scanner"
	"github.com/xxxsen/common/database"
)

type versionActionHandler struct {
	vers   []*VersionAction
	table  string
	ctx    context.Context
	client database.IQueryExecer
}

func newVersionActionHandler(ctx context.Context, client database.IQueryExecer,
	table string, vers []*VersionAction) *versionActionHandler {
	return &versionActionHandler{ctx: ctx, client: client, table: table, vers: vers}
}

func (m *versionActionHandler) ensureVersionValid() error {
	if len(m.vers) == 0 {
		return nil
	}
	sort.Slice(m.vers, func(i, j int) bool {
		left := m.vers[i]
		right := m.vers[j]
		return left.Version < right.Version
	})
	if m.vers[0].LastVersion != 0 || m.vers[0].Version == 0 {
		return fmt.Errorf("invalid init version")
	}
	for i := 1; i < len(m.vers); i++ {
		last := m.vers[i-1]
		cur := m.vers[i]
		if cur.Version == 0 {
			return fmt.Errorf("invalid version info:%+v", *cur)
		}
		if cur.LastVersion != last.Version {
			return fmt.Errorf("missing version detected, current version:%d, last version(missing):%d",
				cur.Version, cur.LastVersion)
		}
	}
	return nil
}

func (m *versionActionHandler) readCurrentVersion() (*Version, bool, error) {
	where := m.versionRecordCondition()
	m.patchSelectLimit(where, 1)
	fields := []string{"id", "version", "ext_info"}
	sql, args, err := builder.BuildSelect(m.table, where, fields)
	if err != nil {
		return nil, false, err
	}
	rows, err := m.client.QueryContext(m.ctx, sql, args...)
	if err != nil {
		return nil, false, err
	}
	rs := make([]*Version, 0, 1)
	if err := scanner.Scan(rows, &rs); err != nil {
		return nil, false, err
	}
	if err := rows.Close(); err != nil {
		return nil, false, err
	}
	if len(rs) == 0 {
		return nil, false, nil
	}
	return rs[0], true, nil
}

func (m *versionActionHandler) ensureVersionTab() error {
	initSQL := fmt.Sprintf(defaultVersionManagementSQL, m.table)
	_, err := m.client.ExecContext(m.ctx, initSQL)
	if err != nil {
		return err
	}
	return nil
}

func (h *versionActionHandler) handle() error {
	if err := h.ensureVersionValid(); err != nil {
		return err
	}
	if err := h.ensureVersionTab(); err != nil {
		return err
	}
	vers, isCreate, err := h.findActionToExec()
	if err != nil {
		return err
	}
	if len(vers) == 0 {
		return nil
	}
	if isCreate {
		return h.doCreate(vers[0])
	}

	return h.doUpgrade(vers)
}

func (m *versionActionHandler) performStepAction(steps []SQLActionStep) error {
	for idx, step := range steps {
		if err := step.Action(m.ctx, m.client); err != nil {
			return fmt.Errorf("perform step fail, idx:%d, name:%s, err:[%w]", idx, step.Name, err)
		}
	}
	return nil
}

func (m *versionActionHandler) doCreate(ver *VersionAction) error {
	if err := m.performStepAction(ver.OnInit.Steps); err != nil {
		return fmt.Errorf("do create failed, err:[%w]", err)
	}
	if err := m.markCurrentVersion(ver.Version, ""); err != nil {
		return fmt.Errorf("do create, mark version failed, err:[%w]", err)
	}
	return nil
}

func (m *versionActionHandler) versionRecordCondition() map[string]interface{} {
	where := map[string]interface{}{
		"id": 1,
	}
	return where
}

func (m *versionActionHandler) patchSelectLimit(mp map[string]interface{}, limit int) {
	mp["_limit"] = []uint{0, uint(limit)}
}

func (m *versionActionHandler) updateVersion(ver uint64, extra string) (int64, error) {
	where := m.versionRecordCondition()
	update := map[string]interface{}{
		"version":  ver,
		"ext_info": extra,
	}
	sql, args, err := builder.BuildUpdate(m.table, where, update)
	if err != nil {
		return 0, err
	}
	rs, err := m.client.ExecContext(m.ctx, sql, args...)
	if err != nil {
		return 0, err
	}
	rows, err := rs.RowsAffected()
	if err != nil {
		return 0, err
	}
	return rows, nil
}

func (m *versionActionHandler) createVersion(ver uint64, extra string) error {
	data := []map[string]interface{}{
		{
			"id":       1,
			"version":  ver,
			"ext_info": extra,
		},
	}
	sql, args, err := builder.BuildInsert(m.table, data)
	if err != nil {
		return err
	}
	if _, err := m.client.ExecContext(m.ctx, sql, args...); err != nil {
		return err
	}
	return nil
}

func (m *versionActionHandler) markCurrentVersion(ver uint64, extra string) error {
	rows, err := m.updateVersion(ver, extra)
	if err != nil {
		return err
	}
	if rows != 0 {
		return nil
	}
	return m.createVersion(ver, extra)
}

func (m *versionActionHandler) doUpgrade(vers []*VersionAction) error {
	for _, ver := range vers {
		if err := m.performStepAction(ver.OnUpgrade.Steps); err != nil {
			return fmt.Errorf("do upgrade failed, err:[%w]", err)
		}
		//mark version
		if err := m.markCurrentVersion(ver.Version, ""); err != nil {
			return fmt.Errorf("do upgrade, mark version failed, err:[%w]", err)
		}
	}
	return nil
}

func (m *versionActionHandler) findActionToExec() ([]*VersionAction, bool, error) {
	stVer, exist, err := m.readCurrentVersion()
	if err != nil {
		return nil, false, err
	}
	if !exist {
		//如果不存在, 那么直接用最新的版本即可。
		return m.vers[len(m.vers)-1:], true, nil
	}
	//记录已经存在了, 那么就从对应版本开始执行升级操作
	for idx := range m.vers {
		ver := m.vers[idx]
		if stVer.Version == ver.Version {
			return m.vers[idx+1:], false, nil
		}
	}
	return nil, false, fmt.Errorf("db version:%d not found in version list", stVer.Version)
}
