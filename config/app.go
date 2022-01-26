package config


type app struct {
	env string
	ll  string
}

func (a *app) GetEnv() string {
	return a.env
}

func (a *app) GetLogLevel() string {
	return a.ll
}
