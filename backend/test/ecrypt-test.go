package main

import (
	"crypto/ed25519"
	"fmt"
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

}
