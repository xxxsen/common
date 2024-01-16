package versionmgr

import (
	"context"

	"github.com/xxxsen/common/database"
)

type ActionFunc func(ctx context.Context, client database.IQueryExecer) error

func SimpleSQLExecAction(execSQL string) ActionFunc {
	return func(ctx context.Context, client database.IQueryExecer) error {
		_, err := client.ExecContext(ctx, execSQL)
		if err != nil {
			return err
		}
		return nil
	}
}
