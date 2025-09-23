package util

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/sha256"

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
	seedHash := sha256.Sum256(seed)

	// Generate a new Ed25519 key pair from the hashed seed
	privkey := ed25519.NewKeyFromSeed(seedHash[:])
	var encPrivKey, encPubKey [32]byte
	copy(encPrivKey[:], seedHash[:])
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

// sha256 hash
func Sha256(data []byte) []byte {
	hash := sha256.Sum256(data)
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
