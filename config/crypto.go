package config

import "github.com/fitant/xbin-api/src/types"

type Crypto struct {
	Salt           []byte
	Cipher         types.CipherSelection
	ARGON2Mem      uint32
	ARGON2Rounds   uint32
	ARGON2IDRounds uint32
}
