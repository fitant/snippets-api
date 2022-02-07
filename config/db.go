package config

import (
	"fmt"
)

type DB struct {
	kind                string
	user                string
	pass                string
	host                string
	port                string
	migrationsPath      string
	database            string
	collection          string
	rsname              string
	timeout             int
	ephemeralCollection ephemeralCollection
}

type ephemeralCollection struct {
	Name string
}

func (db *DB) Address() string {
	switch db.kind {
	case "dnsseed":
		return fmt.Sprintf("mongodb+srv://%s:%s@%s/", db.user, db.pass, db.host)
	case "replicaset":
		return fmt.Sprintf("mongodb://%s:%s@%s/?replicaSet=%s", db.user, db.pass, db.host, db.rsname)
	default: // Use "standalone" as default fallback
		return fmt.Sprintf("mongodb://%s:%s@%s:%s/", db.user, db.pass, db.host, db.port)
	}
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
