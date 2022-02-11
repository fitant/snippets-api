package config

import "github.com/fitant/xbin-api/src/types"


type app struct {
	env string
	Salt string
	Cipher types.CipherSelection
	cipher string
	ll  string
}

func (a *app) GetEnv() string {
	return a.env
}

func (a *app) GetLogLevel() string {
	return a.ll
}
