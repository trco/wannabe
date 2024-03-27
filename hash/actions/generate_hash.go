package actions

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

func GenerateHash(curl string) (string, error) {
	hash := sha256.New()
	_, err := hash.Write([]byte(curl))
	if err != nil {
		return "", fmt.Errorf("GenerateHash: failed writing hash: %v", err)
	}

	encodedHash := hash.Sum(nil)

	return hex.EncodeToString(encodedHash), nil
}
