package db

import (
	"database/sql"
	"os"

	_ "modernc.org/sqlite"
)

var db *sql.DB

const schema = `
CREATE TABLE scheduler (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date CHAR(8) NOT NULL DEFAULT '',
    title VARCHAR(255) NOT NULL DEFAULT '',
    comment TEXT,
    repeat VARCHAR(128) NOT NULL DEFAULT ''
);

CREATE INDEX idx_scheduler_date ON scheduler(date);
`

func Init(dbFile string) error {
	// файл базы
	var install bool
	_, err := os.Stat(dbFile)
	if err != nil {
		install = true
	}

	// db
	db, err = sql.Open("sqlite", dbFile)
	if err != nil {
		return err
	}

	// схема
	if install {
		_, err = db.Exec(schema)
		if err != nil {
			return err
		}
	}

	return nil
}
