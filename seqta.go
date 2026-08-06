// Package setqa implements queries for the SEQTA database.
package seqta

import (
	"fmt"

	_ "github.com/doug-martin/goqu/v9/dialect/postgres"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

// DSN contains the fields needed to create a data source name.
type DSN struct {
	host string
	port uint16
	db   string
	user string
	pwd  string
}

// NewDSN allocates and returns a new [DSN].
func NewDSN(host string, port uint16, db string, user string, pwd string) *DSN {
	return &DSN{host, port, db, user, pwd}
}

// String returns the formatted DSN.
func (a *DSN) String() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%v/%s", a.user, a.pwd, a.host, a.port, a.db)
}

func Open(dsn *DSN) (*sqlx.DB, error) {
	return sqlx.Open("pgx", dsn.String())
}
