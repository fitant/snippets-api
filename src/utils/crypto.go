package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
	"runtime"

	"github.com/Sid-Sun/seaturtle"
	"github.com/fitant/xbin-api/src/types"
	"golang.org/x/crypto/argon2"
)

// Hardcoded ARGON2 Parameters
// 16 iterations for ID generatioon
// Higher as SALT is not unique per app instance
const idHashIterations = 16

// 12 iterations for encryption key
// Unique Salt is used for keys
const keyHashIterations = 12

// 64 MB Memory Space
const hashMemorySpace = 64

// Output hash length
const hashLength = 32

func HashID(id []byte, salt []byte) []byte {
	gid := argon2.IDKey(id, salt, idHashIterations, hashMemorySpace*1024, uint8(runtime.NumCPU()), hashLength)
	return gid
}

func getCipher(selection types.CipherSelection, key []byte) cipher.Block {
	var c cipher.Block
	switch selection {
	case types.SeaTurtle:
		c, _ = seaturtle.NewCipher(key)
	default:
		c, _ = aes.NewCipher(key)
	}
	return c
}

func Encrypt(data []byte, key []byte, cphSel types.CipherSelection) []byte {
	salt := make([]byte, hashLength)
	// Read random values from crypto/rand for key salt
	// Error can be safely ignored
	io.ReadFull(rand.Reader, salt)

	// Derive Key for encryption from ID using Argon2
	key = argon2.IDKey(key, salt, keyHashIterations, hashMemorySpace*1024, uint8(runtime.NumCPU()), hashLength)

	// Generate cipher
	c := getCipher(cphSel, key)

	// use CFB to encrypt full data
	data = cfbEncrypt(data, c)

	// Append salt to the end of data
	data = append(data, salt...)
	return data
}

func Decrypt(data []byte, key []byte, cphSel types.CipherSelection) []byte {
	// Read the salt from end of data
	salt := data[len(data)-hashLength:]

	// Derive Key for decryption from ID using Argon2
	key = argon2.IDKey(key, salt, keyHashIterations, hashMemorySpace*1024, uint8(runtime.NumCPU()), hashLength)

	// Generate cipher
	c := getCipher(cphSel, key)

	// Send IV and data bits to decrypt via CFB
	data = cfbDecrypt(data[:len(data)-hashLength], c)

	// data does not have salt
	return data
}

func cfbEncrypt(data []byte, blockCipher cipher.Block) []byte {
	// Create dst with length of cipher blocksize + data length
	// And initialize first BlockSize bytes pseudorandom for IV
	dst := make([]byte, blockCipher.BlockSize()+len(data))

	// Read random values from crypto/rand for CFB initialization vector
	// Error can be safely ignored
	io.ReadFull(rand.Reader, dst[:blockCipher.BlockSize()])

	// dst from 0 to blockSize is the IV
	cfb := cipher.NewCFBEncrypter(blockCipher, dst[:blockCipher.BlockSize()])
	cfb.XORKeyStream(dst[blockCipher.BlockSize():], data)
	return dst
}

func cfbDecrypt(data []byte, blockCipher cipher.Block) []byte {
	// Create CFB Decrypter with cipher, instantiating with IV (first blockSize blocks of data)
	cfb := cipher.NewCFBDecrypter(blockCipher, data[:blockCipher.BlockSize()])
	// Create variable for storing decrypted note of shorter length taking into account IV
	decrypted := make([]byte, len(data)-blockCipher.BlockSize())
	// Decrypt data starting from blockSize to decrypted
	cfb.XORKeyStream(decrypted, data[blockCipher.BlockSize():])
	return decrypted
}
