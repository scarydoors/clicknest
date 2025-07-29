package clickhouse

import (
	"context"
	"fmt"

	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
)

func TruncateTables(ctx context.Context, conn driver.Conn, tableNames... string) error {
	for _, table := range tableNames {
		query := fmt.Sprintf("TRUNCATE TABLE %s", table)
		if err := conn.Exec(ctx, query); err != nil {
			return err
		}
	}
	return nil
}
