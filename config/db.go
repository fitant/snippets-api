package config

import "fmt"


type DB struct {
	env                 string
	user                string
	pass                string
	host                string
	port                string
	migrationsPath      string
	database            string
	collection          string
	timeout             int
	ephemeralCollection ephemeralCollection
}

type ephemeralCollection struct {
	Name    string
}

func (db *DB) Address() string {
	if db.env == Environments["dev"] {
		return fmt.Sprintf("%s://%s:%s@%s:%s/", "mongodb", db.user, db.pass, db.host, db.port)
	}
	return fmt.Sprintf("%s://%s:%s@%s/?retryWrites=true&w=majority", "mongodb+srv", db.user, db.pass, db.host)
}

func (db *DB) Collection() string {
	return db.collection
}

func (db *DB) MigrationsPath() string {
	return db.migrationsPath
}

func (db *DB) EphemeralCollection() ephemeralCollection {
	return db.ephemeralCollection
}

func (db *DB) Database() string {
	return db.database
}

func (db *DB) TimeoutInSec() int {
	return db.timeout
}
