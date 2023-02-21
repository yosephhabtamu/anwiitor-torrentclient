package main

import (
	"crypto/sha1"
	"encoding/base64"
)

func generatePieceHash(piece []byte) (hash string){
	hasher := sha1.New()
	hasher.Write(piece)
	result := base64.URLEncoding.EncodeToString(hasher.Sum(nil))
	return result
}