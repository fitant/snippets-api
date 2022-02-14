package config

import "github.com/fitant/xbin-api/src/types"

type Service struct {
	Salt []byte
	Overrides map[string]string
	Cipher types.CipherSelection
}
