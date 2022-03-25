package types

type CipherSelection string

const (
	AES       CipherSelection = "aes"
	SeaTurtle CipherSelection = "seaturtle"
)

type EncryptionStack struct {
	ID   string
	Key  []byte
	Hash string
	Salt *[32]byte
}
