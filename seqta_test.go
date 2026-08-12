package seqta

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewDSN(t *testing.T) {
	host := "localhost"
	port := uint16(5432)
	db := "seqta"
	user := "admin"
	pwd := "adminpwd"

	expected := DSN{
		host: "localhost",
		port: 5432,
		db:   "seqta",
		user: "admin",
		pwd:  "adminpwd",
	}

	actual := *NewDSN(host, port, db, user, pwd)

	assert.Equal(t, expected, actual)
}

func TestDSNString(t *testing.T) {
	host := "localhost"
	port := uint16(5432)
	db := "seqta"
	user := "admin"
	pwd := "adminpwd"

	expected := "postgres://admin:adminpwd@localhost:5432/seqta"

	actual := NewDSN(host, port, db, user, pwd).String()

	assert.Equal(t, expected, actual)
}
