package config

import (
	"strings"

	"github.com/fitant/xbin-api/src/types"
	"github.com/spf13/viper"
)

type Config struct {
	App    app
	Svc    Service
	DB     DB
	Http   HTTPServerConfig
	Crypto Crypto
}

var Cfg *Config

func Load() {
	viper.GetViper().AutomaticEnv()
	// App Defaults
	viper.SetDefault("ENV", "dev")
	viper.SetDefault("LOG_LEVEL", "debug")
	viper.SetDefault("CIPHER", "AES")
	// Crypto Defaults
	viper.SetDefault("ARGON2_MEM", 32)
	viper.SetDefault("ARGON2_ROUNDS", 8)
	viper.SetDefault("ARGON2_ID_ROUNDS", 8)
	// DB Defaults
	viper.SetDefault("DB_NAME", "snippets-fitant")
	viper.SetDefault("DB_PORT", "27017")
	viper.SetDefault("DB_EPH_MAX_COUNT", 0)
	viper.SetDefault("DB_TIMEOUT", 10)
	viper.SetDefault("DB_MIGRATIONS", "migrations")
	// Web View Defaults
	viper.SetDefault("HTTP_LISTEN_PORT", 8080)
	viper.SetDefault("HTTP_LISTEN_HOST", "127.0.0.1")
	viper.SetDefault("HTTP_BASE_URL", "http://localhost:8080")
	viper.SetDefault("HTTP_RETURN_FORMAT", "raw")
	viper.SetDefault("HTTP_API_ENDPOINT", "/snippets")
	viper.SetDefault("HTTP_CORS_LIST", "http://localhost:*")
	viper.SetDefault("DB_RSNAME", "rs0")
	cfg := &Config{
		App: app{
			env: viper.GetString("ENV"),
			ll:  viper.GetString("LOG_LEVEL"),
		},
		Svc: Service{
			Overrides: make(map[string]string),
		},
		Crypto: Crypto{
			Salt:         []byte(viper.GetString("SALT")),
			ARGON2Mem:    viper.GetUint32("ARGON2_MEM"),
			ARGON2Rounds: viper.GetUint32("ARGON2_ROUNDS"),
			ARGON2IDRounds: viper.GetUint32("ARGON2_ID_ROUNDS"),
		},
		DB: DB{
			kind:           viper.GetString("DB_TYPE"),
			rsname:         viper.GetString("DB_RSNAME"),
			timeout:        viper.GetInt("DB_TIMEOUT"),
			user:           viper.GetString("DB_USER"),
			pass:           viper.GetString("DB_PASS"),
			host:           viper.GetString("DB_HOST"),
			port:           viper.GetString("DB_PORT"),
			database:       viper.GetString("DB_NAME"),
			collection:     "snippets",
			migrationsPath: viper.GetString("DB_MIGRATIONS"),
			ephemeralCollection: ephemeralCollection{
				Name: "eph_snippets",
			},
		},
		Http: HTTPServerConfig{
			host:         viper.GetString("HTTP_LISTEN_HOST"),
			port:         viper.GetInt("HTTP_LISTEN_PORT"),
			CORS:         viper.GetString("HTTP_CORS_LIST"),
			baseURL:      viper.GetString("HTTP_BASE_URL"),
			Enpoint:      viper.GetString("HTTP_API_ENDPOINT"),
			returnFormat: viper.GetString("HTTP_RETURN_FORMAT"),
		},
	}

	switch viper.GetString("CIPHER") {
	case "SEAT":
		if viper.GetBool("CIPHER_UNTESTED") {
			cfg.Crypto.Cipher = types.SeaTurtle
			break
		}
		fallthrough
	default:
		cfg.Crypto.Cipher = types.AES
	}

	overrides := viper.GetString("OVERRIDES")
	for _, override := range strings.Split(overrides, ",") {
		entry := strings.Split(override, ":")
		if len(entry) == 2 {
			cfg.Svc.Overrides[entry[0]] = entry[1]
		}
	}

	if len(cfg.Crypto.Salt) == 0 {
		panic("config missing: salt not provided")
	}

	Cfg = cfg
}
