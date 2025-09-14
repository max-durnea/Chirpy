package auth

import (
	"crypto/rand"
	"encoding/hex"
)
func MakeRefreshToken() (string, error){
	key := make([]byte, 32)
	rand.Read(key)

	data_str := hex.EncodeToString(key)
	return data_str, nil
}