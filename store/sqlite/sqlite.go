package sqlite

import (
	"database/sql"
	"embed"

	sq "github.com/Masterminds/squirrel"
	"github.com/august-kuhfuss/hotdamn/store"
	"github.com/pressly/goose/v3"

	_ "modernc.org/sqlite"
)

type sqlite struct {
	DB *sq.StatementBuilderType
}

var (
	driver = "sqlite"
	dbURL  string
	conn   *sql.DB
)

func NewStore(dataSource string) (store.Store, error) {
	dbURL = dataSource

	var err error
	conn, err = sql.Open(driver, dbURL)
	if err != nil {
		return nil, err
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

//go:embed migrations/*.sql
var migrations embed.FS

func init() {
	goose.SetDialect(driver)
	goose.SetBaseFS(migrations)
}

func MigrateUp() error {
	db, err := sql.Open(driver, dbURL)
	if err != nil {
		return err
	}
	defer db.Close()
	if err := goose.Up(db, "migrations"); err != nil {
		return err
	}
	return nil
}

func MigrateDown() error {
	db, err := sql.Open(driver, dbURL)
	if err != nil {
		return err
	}
	defer db.Close()
	if err := goose.Down(db, "migrations"); err != nil {
		return err
	}

	return nil
}
