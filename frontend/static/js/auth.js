import { generateKeyPairFromSeed, sign} from './lib/crypto-utils.js';

// 服务端设置HttpOnly Cookie
// Set-Cookie: jwt=xxx; HttpOnly; Secure; SameSite=Strict; Max-Age=3600

// 客户端处理
class AuthManager {
    constructor() {
        this.csrfToken = null;
        this.privateKey = null;
        this.publicKey = null;
    }

    async prelogin(username) {
        const response = await fetch(`/api/v0/prelogin`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            credentials: 'include',
            body: JSON.stringify({ 'username': username })
        });

        if (response.ok) {
            const data = await response.json();
            if (data.salt){
                return data.salt;
            } else if (data.error) {
                throw new Error(data.error);
            } else {
                throw new Error('Invalid prelogin response');
            }
        } else {
            throw new Error('Prelogin failed');
        }
    }

    // 计算密码哈希值
    async hashPassword(password) {
        const encoder = new TextEncoder();
        const data = encoder.encode(password);
        const hashBuffer = await crypto.subtle.digest('SHA-256', data);
        return new Uint8Array(hashBuffer);
    }

    // 从哈希生成Ed25519密钥对
    async generateKeyPairFromPasswordHash(passwordHash) {
        // 使用密码哈希作为种子生成Ed25519密钥对
        const keyPair = await generateKeyPairFromSeed(passwordHash);
        return keyPair;
    }

    // 保存私钥到localStorage（可选择其他存储方式）
    savePrivateKeyLocally(privateKey) {
        // 将私钥转换为base64存储
        const privateKeyBase64 = btoa(String.fromCharCode(...privateKey));
        localStorage.setItem('user_private_key', privateKeyBase64);
    }

    // 从本地加载私钥
    loadPrivateKeyFromLocal() {
        const privateKeyBase64 = localStorage.getItem('user_private_key');
        if (privateKeyBase64) {
            const privateKeyBytes = Uint8Array.from(atob(privateKeyBase64), c => c.charCodeAt(0));
            return privateKeyBytes;
        }
        return null;
    }

    async login(username, password, salt) {
        try {
            // 1. 计算password的哈希值
            const pwdHash = await this.hashPassword(password);
            console.log('Password hash generated');

            // 2. 用pwdHash生成Ed25519的密钥对
            const keyPair = await this.generateKeyPairFromPasswordHash(pwdHash);
            this.privateKey = keyPair.privateKey;
            console.log('Ed25519 key pair generated');



            // 3. 用私钥签名salt
            const saltBytes = new TextEncoder().encode(salt);
            const signature = await sign(this.privateKey, saltBytes);
            console.log('Salt signed with private key');

            // 4. 将私钥保存在本地
            this.savePrivateKeyLocally(this.privateKey);
            console.log('Private key saved locally');

            // 5. 发送签名和用户名给服务器
            const response = await fetch('/api/v0/login', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                credentials: 'include',
                body: JSON.stringify({ 
                    username: username,
                    signature: Array.from(signature), // 转换为数组发送
                })
            });

            if (response.ok) {
                const data = await response.json();
                this.csrfToken = data.csrfToken;
                console.log('Login successful');
                return true;
            } else {
                console.error('Login failed:', response.statusText);
                return false;
            }

        } catch (error) {
            console.error('Login error:', error);
            return false;
        }
    }
    
    async makeAuthenticatedRequest(url, options = {}) {
        return fetch(url, {
            ...options,
            credentials: 'include',
            headers: {
                ...options.headers,
                'X-CSRF-Token': this.csrfToken
            }
        });
    }

    // 清除本地存储的私钥
    clearLocalKeys() {
        localStorage.removeItem('user_private_key');
        this.privateKey = null;
        this.publicKey = null;
    }

    // 登出时清除密钥
    async logout() {
        try {
            const response = await fetch('/api/v0/logout', {
                method: 'POST',
                credentials: 'include'
            });
            
            this.clearLocalKeys();
            this.csrfToken = null;
            
            return response.ok;
        } catch (error) {
            console.error('Logout error:', error);
            return false;
        }
    }
}

export { AuthManager };