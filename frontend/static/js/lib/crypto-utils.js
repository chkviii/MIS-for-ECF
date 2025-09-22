// 从种子生成Ed25519密钥对
export async function generateKeyPairFromSeed(seed) {
    // 确保种子是32字节
    let seedBytes;
    if (seed.length >= 32) {
        seedBytes = seed.slice(0, 32);
    } else {
        console.error('Seed must be at least 32 bytes');
        return null;
    }

    // 使用Web Crypto API生成Ed25519密钥对
    try {
        // 导入种子作为原始密钥材料
        const keyMaterial = await crypto.subtle.importKey(
            'raw',
            seedBytes,
            { name: 'Ed25519' },
            false,
            ['sign']
        );

        // 生成密钥对
        const keyPair = await crypto.subtle.generateKey(
            {
                name: 'Ed25519',
            },
            true,
            ['sign', 'verify']
        );

        // 导出密钥为原始字节
        const privateKeyBytes = await crypto.subtle.exportKey('pkcs8', keyPair.privateKey);
        const publicKeyBytes = await crypto.subtle.exportKey('spki', keyPair.publicKey);

        return {
            privateKey: new Uint8Array(privateKeyBytes),
            publicKey: new Uint8Array(publicKeyBytes)
        };

    } catch (error) {
        console.error('Key generation error:', error);
        return null;
        // 如果Web Crypto API不支持Ed25519，使用备用方案
        // return generateKeyPairFallback(seedBytes);
    }
}

// // 备用密钥生成方案（如果浏览器不支持Ed25519）
// async function generateKeyPairFallback(seed) {
//     // 这里可以使用第三方库如 @noble/ed25519
//     // 或者回退到其他支持的算法
//     console.warn('Using fallback key generation method');
    
//     // 使用ECDSA作为备用方案
//     const keyPair = await crypto.subtle.generateKey(
//         {
//             name: 'ECDSA',
//             namedCurve: 'P-256'
//         },
//         true,
//         ['sign', 'verify']
//     );

//     const privateKeyBytes = await crypto.subtle.exportKey('pkcs8', keyPair.privateKey);
//     const publicKeyBytes = await crypto.subtle.exportKey('spki', keyPair.publicKey);

//     return {
//         privateKey: new Uint8Array(privateKeyBytes),
//         publicKey: new Uint8Array(publicKeyBytes),
//         algorithm: 'ECDSA' // 标记使用的算法
//     };
// }

// 使用私钥签名数据
export async function sign(privateKey, data) {
    try {
        // 导入私钥
        const key = await crypto.subtle.importKey(
            'pkcs8',
            privateKey,
            { name: 'Ed25519' },
            false,
            ['sign']
        );

        // 签名数据
        const signature = await crypto.subtle.sign(
            'Ed25519',
            key,
            data
        );

        return new Uint8Array(signature);

    } catch (error) {
        console.error('Signing error:', error);
        return null;
        // 使用备用签名方法
        // return signFallback(privateKey, data);
    }
}

// // 备用签名方案
// async function signFallback(privateKey, data) {
//     try {
//         const key = await crypto.subtle.importKey(
//             'pkcs8',
//             privateKey,
//             {
//                 name: 'ECDSA',
//                 namedCurve: 'P-256'
//             },
//             false,
//             ['sign']
//         );

//         const signature = await crypto.subtle.sign(
//             {
//                 name: 'ECDSA',
//                 hash: 'SHA-256'
//             },
//             key,
//             data
//         );

//         return new Uint8Array(signature);
//     } catch (error) {
//         console.error('Fallback signing failed:', error);
//         throw error;
//     }
// }

// 验证签名
export async function verify(publicKey, signature, data) {
    try {
        const key = await crypto.subtle.importKey(
            'spki',
            publicKey,
            { name: 'Ed25519' },
            false,
            ['verify']
        );

        return await crypto.subtle.verify(
            'Ed25519',
            key,
            signature,
            data
        );

    } catch (error) {
        console.error('Verification error:', error);
        return false;
    }
}