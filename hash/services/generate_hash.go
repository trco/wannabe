package services

import (
	"wannabe/hash/actions"
)

func GenerateHash(curl string) (string, error) {
	return actions.GenerateHash(curl)
}
