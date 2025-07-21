package clickhouse

import (
	"context"
	"fmt"
)

func (c *ClickhouseDB) TruncateTables(ctx context.Context, tableNames... string) error {
	for _, table := range tableNames {
		query := fmt.Sprintf("TRUNCATE TABLE %s", table)
		if err := c.conn.Exec(ctx, query); err != nil {
			return err
		}
	}
	return nil
}
