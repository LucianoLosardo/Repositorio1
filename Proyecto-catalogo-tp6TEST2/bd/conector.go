package bd

import (
	"database/sql"

	_ "github.com/lib/pq"
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://grupoWeb:CatalogoGrupoWeb@database:5432/catalogoDB?sslmode=disable"
)

func ConnectDB() (*sql.DB, error) {
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		return nil, err
	}

	if err := conn.Ping(); err != nil {
		return nil, err
	}
	return conn, nil
}
