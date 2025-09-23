package main

import (
	"crypto/ed25519"
	"crypto/sha256"
	"fmt"

	"golang.org/x/crypto/curve25519"
)

func main() {
	pubkey, privkey, err := ed25519.GenerateKey(nil)
	if err != nil {
		fmt.Println("Error generating key pair:", err)
		return
	}
	fmt.Println("Public Key:", pubkey)
	fmt.Println("Private Key:", privkey)

	fmt.Println([]byte(`{"pubkey":"` + string(pubkey) + `"}`))

	seedHash := sha256.Sum256([]byte("my secret seed"))

	// 为不同用途派生不同的种子
	signSeed := sha256.Sum256(append(seedHash[:], []byte("signing")...))
	encSeed := signSeed

	// 创建签名密钥对
	signPrivKey := ed25519.NewKeyFromSeed(signSeed[:])
	signPubKey := signPrivKey.Public().(ed25519.PublicKey)

	// 创建加密密钥对
	var encPrivKey, encPubKey [32]byte
	copy(encPrivKey[:], encSeed[:])
	curve25519.ScalarBaseMult(&encPubKey, &encPrivKey)

	fmt.Println("Seed:", signSeed)
	fmt.Println("Sign Public Key:", signPubKey)
	fmt.Println("Sign Private Key:", signPrivKey)
	fmt.Println("Enc Public Key:", encPubKey)
	fmt.Println("Enc Private Key:", encPrivKey)

}
