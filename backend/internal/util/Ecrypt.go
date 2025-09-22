package util

import (
	"crypto/ed25519"
	"crypto/sha256"

	"mypage-backend/internal/config"
)

var privkey ed25519.PrivateKey

func KeyInit() ([]byte, ed25519.PublicKey, ed25519.PrivateKey) {
	//Load key string in config
	seed := []byte(config.GlobalConfig.Ecrypt_Seed)
	//hash the seed to 32 bytes
	seedHash := sha256.Sum256(seed)
	aeskey := seedHash[:]

	// Generate a new Ed25519 key pair from the hashed seed
	privkey = ed25519.NewKeyFromSeed(seedHash[:])
	pubkey := privkey.Public().(ed25519.PublicKey)

	return aeskey, pubkey, privkey
}

func SrvPubKey() ed25519.PublicKey {
	return privkey.Public().(ed25519.PublicKey)
}

func SrvPrivKey() ed25519.PrivateKey {
	return privkey
}
