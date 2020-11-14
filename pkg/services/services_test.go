package services

import "simple-auth/pkg/db"

var testDB db.SADB

func getDB() db.SADB {
	if testDB == nil {
		testDB = db.New("sqlite3", "file::memory:?cache=shared")
	}
	return testDB
}
