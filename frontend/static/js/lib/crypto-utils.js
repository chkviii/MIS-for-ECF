// TODO: ALL OF THIS IS WRONG See 'https://developer.mozilla.org/zh-CN/docs/Web/API/SubtleCrypto' to proceed

import { x25519ScalarMult ,x25519ScalarBaseMult } from "./curve25519.js";

// 从种子生成Curve25519密钥对
export async function generateX25519KeyPairFromSeed(seed) {
    // Hash the seed with SHA-512
    const hash = await crypto.subtle.digest('SHA-512', seed);
    const hashArray = new Uint8Array(hash);

    // Clamp the first 32 bytes to form the private scalar
    hashArray[0] &= 248;
    hashArray[31] &= 63;
    hashArray[31] |= 64;
    const privateKey = hashArray.slice(0, 32);

    // Derive the public key using X25519 scalar multiplication
    const publicKey = await x25519ScalarBaseMult(privateKey);
    return { privateKey, publicKey };
}

export async function deriveSharedSecret(privateKey, publicKey) {
    // Perform X25519 scalar multiplication to derive the shared secret
    const sharedSecret = await x25519ScalarMult(privateKey, publicKey);
    return sharedSecret;
}


//从种子生成Ed25519密钥对
export async function generateEd25519KeyPairFromSeed(seed) {
    // Hash the seed with SHA-512
    const hash = await crypto.subtle.digest('SHA-512', seed); 
    const hashArray = new Uint8Array(hash);

    // Clamp the first 32 bytes to form the private scalar
    hashArray[0] &= 248;
    hashArray[31] &= 63;
    hashArray[31] |= 64;
    const privateKey = hashArray.slice(0, 32);

