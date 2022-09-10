package testutil

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func OpenDBForTest(t *testing.T) *sqlx.DB {
	t.Helper()

	port := 55432
	if _, defined := os.LookupEnv("CI"); defined {
		port = 5432
	}
	db, err := sql.Open(
		"postgres",
		fmt.Sprintf(
			"host=localhost port=%d dbname=todotest user=todo password=todo sslmode=disable",
			port,
		),
	)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(
		func() { _ = db.Close() },
	)
	return sqlx.NewDb(db, "postgres")
}
