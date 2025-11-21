package main

import (
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/hex"
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

	//test point RFC 7748
	hexStr := "77076d0a7318a57d3c16c17251b26645df4c2f87ebc0992ab177fba51db92c2a"
	scalarBytes, err := hex.DecodeString(hexStr)
	if err != nil {
		panic(err)
	}

	var scalar, result [32]byte
	copy(scalar[:], scalarBytes)
	curve25519.ScalarBaseMult(&result, &scalar)
	fmt.Printf("Scalar:   %x\n", scalar)
	fmt.Printf("Resulting: %x\n", result)

	var derivedKey1, derivedKey2 [32]byte
	curve25519.ScalarMult(&derivedKey1, &scalar, &encPubKey)
	curve25519.ScalarMult(&derivedKey2, &encPrivKey, &result)
	fmt.Printf("Derived Key 1: %x\n", derivedKey1)
	fmt.Printf("Derived Key 2: %x\n", derivedKey2)

	var result1 [32]byte
	var dummy, scalar2 [32]byte

	for i := range scalar2 {
		scalar2[i] = byte(255)
	}

	for i := range 256 {
		dummy[31] = byte(i)

		curve25519.ScalarMult(&result1, &scalar2, &dummy)
		fmt.Printf("Result with dummy %d: %x\n", i, result1)
	}

}
