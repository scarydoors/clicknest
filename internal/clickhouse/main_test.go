package clickhouse_test

import (
	"context"
	"database/sql"
	"log"
	"strings"
	"testing"

	_ "github.com/ClickHouse/clickhouse-go/v2"
	"github.com/pressly/goose/v3"
	"github.com/scarydoors/clicknest/internal/clickhouse"
	"github.com/testcontainers/testcontainers-go"
	clickhousetc "github.com/testcontainers/testcontainers-go/modules/clickhouse"
)

var (
	clickhouseDB *clickhouse.ClickhouseDB
)

func TestMain(m *testing.M) {
	ctx := context.Background()
	user := "default"
	password := ""
	dbname := "default"

	clickhouseContainer, err := clickhousetc.Run(ctx,

		"clickhouse/clickhouse-server:23.3.8.21-alpine",
		clickhousetc.WithUsername(user),
		clickhousetc.WithPassword(password),
		clickhousetc.WithDatabase(dbname),
	)

	defer func() {
		if err := testcontainers.TerminateContainer(clickhouseContainer); err != nil {
			log.Printf("failed to terminate container: %s", err)
		}
	}()

	if err != nil {
		log.Printf("failed to start container: %s", err)
		return
	}

	host, err := clickhouseContainer.ConnectionHost(ctx)
	if err != nil {
		log.Printf("connection host: %s", err)
	}
	split := strings.Split(host, ":")
	
	config := clickhouse.ClickhouseDBConfig{
		Host: split[0],
		Port: split[1],
		Database: dbname,
		Username: user, 
		Password: password,
	}

	dsn, err := clickhouseContainer.ConnectionString(ctx)
	if err != nil {
		log.Fatalf("what %s", err)
	}

	db, err := sql.Open("clickhouse", dsn)
	if err != nil {
		log.Fatalf("what %s", err)
	}

	// Run migrations
	goose.SetDialect("clickhouse")
	if err := goose.Up(db, "../../migrations"); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	clickhouseDB, err = clickhouse.NewClickhouseDB(ctx, config)

	if err != nil {
		log.Printf("new clickhouse db: %s", err)
	}
	defer clickhouseDB.Close()

	m.Run()
}
