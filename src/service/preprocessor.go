package service

import (
	"encoding/base64"

	"github.com/fitant/xbin-api/src/types"
	"github.com/fitant/xbin-api/src/utils"
)

var encryptionKeys chan types.EncryptionStack

func populateEncryptionStack(idSize int) {
	for {
		x := types.EncryptionStack{
			ID:   utils.GenerateID(idSize),
			Salt: utils.GenSalt(),
		}
		id := []byte(x.ID)
		x.Hash = base64.StdEncoding.EncodeToString(utils.HashID(id))
		x.Key = utils.GenKey(id, x.Salt[:])
		encryptionKeys <- x
	}
}
