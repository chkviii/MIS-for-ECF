package util

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/sha512"
	"errors"

	"mypage-backend/internal/config"

	"golang.org/x/crypto/curve25519"
)

type Ecrypt struct {
	SignPubKey  ed25519.PublicKey
	SignPrivKey ed25519.PrivateKey
	EncPubKey   []byte
	EncPrivKey  []byte
}

var keys Ecrypt

func KeyInit() {
	//Load key string in config
	seed := []byte(config.GlobalConfig.Ecrypt_Seed)
	//hash the seed to 32 bytes
	seedHash := SHA512(seed)

	//clamp the private key acording to RFC 7748
	seedHash[0] &= 248
	seedHash[31] &= 127
	seedHash[31] |= 64

	// Generate a new Ed25519 key pair from the hashed seed
	privkey := ed25519.NewKeyFromSeed(seedHash[:32])
	var encPrivKey, encPubKey [32]byte
	copy(encPrivKey[:], seedHash[:32])
	curve25519.ScalarBaseMult(&encPubKey, &encPrivKey)

	keys = Ecrypt{
		SignPubKey:  privkey.Public().(ed25519.PublicKey),
		SignPrivKey: privkey,
		EncPubKey:   encPubKey[:],
		EncPrivKey:  encPrivKey[:],
	}
}

func SrvSignPubKey() ed25519.PublicKey {
	return keys.SignPubKey
}

func SrvSignPrivKey() ed25519.PrivateKey {
	return keys.SignPrivKey
}

func SrvEncPubKey() []byte {
	return keys.EncPubKey
}

func SrvEncPrivKey() []byte {
	return keys.EncPrivKey
}

// sha512 hash
func SHA512(data []byte) []byte {
	hash := sha512.Sum512(data)
	return hash[:]
}

// Encrypt message with recipient's public key using X25519 and AES-GCM
func ECDHEncrypt(pubKey, message []byte) ([]byte, error) {
	pubKeyArr := make([]byte, 32)
	privKeyArr := make([]byte, 32)
	copy(pubKeyArr[:], pubKey)
	copy(privKeyArr[:], keys.EncPrivKey)
	copy(privKeyArr, keys.EncPrivKey)
	var sharedKey []byte
	var err error
	sharedKey, err = curve25519.X25519(privKeyArr, pubKeyArr)
	if err != nil {
		return nil, err
	}

	//check if sharedKey is all zero
	isAllZero := true
	for _, b := range sharedKey {
		if b != 0 {
			isAllZero = false
			break
		}
	}

	if isAllZero {
		return nil, errors.New("failed to derive shared key: invalid public key")
	}

	ciphertext, err := AESEncrypt(sharedKey, message)
	if err != nil {
		return nil, err
	}

	return ciphertext, nil
}

// AES-GCM encrypt message with given key
func AESEncrypt(aeskey []byte, message []byte) ([]byte, error) {
	block, err := aes.NewCipher(aeskey)
	if err != nil {
		return nil, err
	}
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, aesGCM.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return nil, err
	}
	ciphertext := aesGCM.Seal(nonce, nonce, message, nil)
	return ciphertext, nil
}

// AES-GCM decrypt message with given key
func AESDecrypt(aeskey []byte, ciphertext []byte) ([]byte, error) {
	block, err := aes.NewCipher(aeskey)
	if err != nil {
		return nil, err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := aesGCM.NonceSize()
	// if len(ciphertext) < nonceSize {
	// 	return nil, err
	// }

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)

	if err != nil {
		return nil, err
	}

	return plaintext, nil
}
