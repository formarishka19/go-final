package db

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

const SCHEMA = `CREATE TABLE IF NOT EXISTS scheduler (
    id      INTEGER PRIMARY KEY AUTOINCREMENT,
    date    CHAR(8) NOT NULL DEFAULT "",
    title   VARCHAR(128) NOT NULL,
    comment TEXT,
    repeat  VARCHAR(128)
);
CREATE INDEX IF NOT EXISTS idx_date ON scheduler (date);`

var DB *sql.DB

func Init(dbFile string) error {
	_, err := DB.Exec(SCHEMA)
	if err != nil {
		return fmt.Errorf("error generating db schema %w", err)
	}
	return nil
}
