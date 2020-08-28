package main

import (
	"simple-auth/pkg/config"
	"simple-auth/pkg/db"
)

func getDB() db.SADB {
	config := config.Load(false)
	return db.New(config.Db.Driver, config.Db.URL)
}
