package config

import "github.com/spf13/viper"

type Config struct {
	App  app
	DB   DB
	Http HTTPServerConfig
}

func Load() *Config {
	viper.GetViper().AutomaticEnv()
	// App Defaults
	viper.SetDefault("ENV", "dev")
	viper.SetDefault("LOG_LEVEL", "debug")
	// DB Defaults
	viper.SetDefault("DB_ENV", "dev")
	viper.SetDefault("DB_EPH_MAX_COUNT", 0)
	viper.SetDefault("DB_TIMEOUT", 10)
	viper.SetDefault("DB_EPH_CAPPED", true)
	viper.SetDefault("DB_MIGRATIONS", "migrations")
	// Web View Defaults
	viper.SetDefault("WEB_LISTEN_PORT", 8080)
	viper.SetDefault("WEB_LISTEN_HOST", "127.0.0.1")
	return &Config{
		App: app{
			env: viper.GetString("ENV"),
			ll:  viper.GetString("LOG_LEVEL"),
		},
		DB: DB{
			env:            viper.GetString("ENV"),
			timeout:        viper.GetInt("DB_TIMEOUT"),
			user:           viper.GetString("DB_USER"),
			pass:           viper.GetString("DB_PASS"),
			host:           viper.GetString("DB_HOST"),
			port:           viper.GetString("DB_PORT"),
			database:       viper.GetString("DB_NAME"),
			collection:     viper.GetString("DB_COLLECTION"),
			migrationsPath: viper.GetString("DB_MIGRATIONS"),
			ephemeralCollection: ephemeralCollection{
				Name: viper.GetString("DB_EPH_COLLECTION"),
			},
		},
		Http: HTTPServerConfig{
			host: viper.GetString("WEB_LISTEN_HOST"),
			port: viper.GetInt("WEB_LISTEN_PORT"),
		},
	}
}