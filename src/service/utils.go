package service

import (
	"strings"

	"github.com/leonklingele/passphrase"
)

func generateID(n int) string {
	passphrase.Separator = "-"
	id, _ := passphrase.Generate(n)

	idFields := strings.Split(id, "-")
	for i := range idFields {
		idFields[i] = strings.Title(idFields[i])
	}
	return strings.Join(idFields, "")
}
