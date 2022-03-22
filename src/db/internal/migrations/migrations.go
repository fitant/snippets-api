package migrations

import (
	"fmt"

	"github.com/fitant/xbin-api/config"
	"github.com/golang-migrate/migrate/v4"
	mongoMigrate "github.com/golang-migrate/migrate/v4/database/mongodb"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"go.mongodb.org/mongo-driver/mongo"
)

func Migrate(db *mongo.Client, cfg *config.Config) {
	migInstance, err := mongoMigrate.WithInstance(db, &mongoMigrate.Config{
		DatabaseName: cfg.DB.Database(),
	})
	if err != nil {
		panic(fmt.Sprintf("%s : %v", "[DB] [Migrate] [WithInstance]", err))
	}

	migrations, err := migrate.NewWithDatabaseInstance(fmt.Sprintf("file://%s", cfg.DB.MigrationsPath()),
		cfg.DB.Database(), migInstance)
	if err != nil {
		panic(fmt.Sprintf("%s : %v", "[DB] [Migrate] [NewWithDatabaseInstance]", err))
	}
	migrations.Up()

}
