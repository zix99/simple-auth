package db_test

import "simple-auth/pkg/db"

var sadb db.SADB

func init() {
	sadb = db.New("sqlite3", "file::memory:?cache=shared")
}
