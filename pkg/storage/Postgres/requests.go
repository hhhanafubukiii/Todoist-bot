package postgres

import "database/sql"

type TokenStorage struct {
	db   *sql.DB
	conn *sql.Conn
}
