package sqlitestore

import (
	"database/sql"

	sq "github.com/Masterminds/squirrel"
	"github.com/august-kuhfuss/hotdamn/store"

	_ "modernc.org/sqlite"
)

type sqlite struct {
	DB *sq.StatementBuilderType
}

var conn *sql.DB

func NewStore(dataSource string) (store.Store, error) {
	var err error
	conn, err = sql.Open("sqlite", dataSource)
	if err != nil {
		panic(err)
	}

	cache := sq.NewStmtCache(conn)
	db := sq.StatementBuilder.RunWith(cache).PlaceholderFormat(sq.Question)

	return &sqlite{
		DB: &db,
	}, nil
}

func Ping() error {
	return conn.Ping()
}

func Close() error {
	return conn.Close()
}

func MigrateUp() error {
	_, err := conn.Exec(`
		CREATE TABLE IF NOT EXISTS measurements (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			temperature REAL,
			humidity REAL,
			created_at DATETIME
		);
	`)
	return err
}
