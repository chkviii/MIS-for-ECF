// X25519 implementation in JavaScript

// Modular arithmetic for Curve25519
const P = 0x7FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEDn;
const basepoint = new Uint8Array([9, ...new Array(31).fill(0)])

async function mod(a) {
    return ((a % P) + P) % P;
}

async function modInverse(a) {
    // Fermat's little theorem: a^(p-2) ≡ a^(-1) (mod p)
    return modPow(a, P - 2n);
}

async function modPow(base, exp) {
    let result = 1n;
    base = mod(base);
    while (exp > 0n) {
        if (exp & 1n) {
            result = mod(result * base);
        }
        base = mod(base * base);
        exp >>= 1n;
    }
    return result;
}

// Curve25519 scalar multiplication
async function x25519ScalarMult(scalar, u) {
    // Ensure scalar and u are Uint8Arrays of length 32
    if (scalar.length !== 32 || u.length !== 32) {
        throw new Error('Scalar and u-coordinate must be 32 bytes');
    }
    
    // Clamp scalar
    const s = new Uint8Array(scalar);
    s[0] &= 248;
    s[31] &= 127;
    s[31] |= 64;
    
    // Convert u to bigint (little-endian)
    let x = 0n;
    for (let i = 31; i >= 0; i--) {
        x = (x << 8n) + BigInt(u[i]);
    }
    x = mod(x);
    
    // Montgomery ladder
    let x1 = x;
    let x2 = 1n;
    let z2 = 0n;
    let x3 = x;
    let z3 = 1n;
    
    for (let i = 254; i >= 0; i--) {
        const bit = (s[Math.floor(i / 8)] >> (i % 8)) & 1;
        
        // Conditional swap
        if (bit) {
            [x2, x3] = [x3, x2];
            [z2, z3] = [z3, z2];
        }
        
        // Montgomery ladder step
        const A = mod(x2 + z2);
        const AA = mod(A * A);
        const B = mod(x2 - z2);
        const BB = mod(B * B);
        const E = mod(AA - BB);
        const C = mod(x3 + z3);
        const D = mod(x3 - z3);
        const DA = mod(D * A);
        const CB = mod(C * B);
        
        x3 = mod((DA + CB) * (DA + CB));
        z3 = mod(x1 * (DA - CB) * (DA - CB));
        x2 = mod(AA * BB);
        z2 = mod(E * (AA + 121665n * E));
        
        // Conditional swap back
        if (bit) {
            [x2, x3] = [x3, x2];
            [z2, z3] = [z3, z2];
        }
    }
    
    // Final inversion
    const result = mod(x2 * modInverse(z2));
    // check result is not zero
    if (result === 0n) {
        throw new Error('Resulting point is at infinity');
    }

    // Convert back to bytes (little-endian)
    const output = new Uint8Array(32);
    let temp = result;
    for (let i = 0; i < 32; i++) {
        output[i] = Number(temp & 0xFFn);
        temp >>= 8n;
    }
    
    return output;
}

async function x25519ScalarBaseMult(scalar) {
    return x25519ScalarMult(scalar, basepoint);
}

export { x25519ScalarMult, x25519ScalarBaseMult };