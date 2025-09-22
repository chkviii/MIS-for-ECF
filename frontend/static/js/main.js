import { ModalManager } from './modal.js';

class App{
    constructor() {
        this.isSidebarOpen = true;
        this.isLoggedIn = false;
        this.modalManager = new ModalManager();
        this.init();
    }

    init() {
        this.loggedIn()
        this.initSidebar();
        this.initAuth(); 
    }

    // Initialize authentication elements and event listeners
    initAuth() {
        // Bind buttons
        this.loginBtn = document.getElementById('loginBtn');
        this.registerBtn = document.getElementById('registerBtn');
        this.logoutBtn = document.getElementById('logoutBtn');
        // TODO: 用户中心

        // Bind modal elements
        this.loginModal = document.getElementById('loginModal');

        // Event listeners
        this.loginBtn?.addEventListener('click', () => {
            const response = fetch('/api/v0/newloginsession', {
                method: 'GET',
                credentials: 'include'
            }).then(res => res.json());

            // Assume response contains a salt for signing
            this.loginSalt = response.salt;
            this.modalManager.show('loginModal');
        });

        loginForm.addEventListener('submit', (e) => {
            e.preventDefault(); // 阻止默认提交
            this.handleLogin(e.target);
        });

        this.registerBtn?.addEventListener('click', () => {
            //go to register page
            window.location.href = '/register';
        });
        this.logoutBtn?.addEventListener('click', () => {
            // Handle logout logic here
        });
    }

    // Check if user is logged in and update UI accordingly
    loggedIn() {
        const jwt = document.cookie.split('; ').find(row => row.startsWith('jwt='));
        if (jwt) {
            const exptime = jwt ? parseInt(atob(jwt.split('.')[1]).exp) : 0;
            if (Date.now() < exptime * 1000) {
                this.isLoggedIn = true;
            }
        }
    }

    SwithLoginState() {
        this.isLoggedIn = !this.isLoggedIn;
        if (this.isLoggedIn) {
            this.loginBtn?.classList.add('hidden');
            this.registerBtn?.classList.add('hidden');
            this.logoutBtn?.classList.remove('hidden');
        } else {
            this.loginBtn?.classList.remove('hidden');
            this.registerBtn?.classList.remove('hidden');
            this.logoutBtn?.classList.add('hidden');
        }


    }

    initSidebar() {
        this.toggleButton = document.getElementById('toggleSidebar');
        this.sidebar = document.getElementById('sidebar');
        this.main = document.getElementById('main');
        this.footer = document.getElementById('footer');

        // console.log(1, this.main, this.footer);

        
        if (!this.toggleButton || !this.sidebar) {
            this.main?.classList.add('main-expanded');
            this.footer?.classList.add('footer-expanded');
            // console.warn('Sidebar elements not found');
            return;
        }

        // 直接调用方法，不需要传参
        this.toggleButton.addEventListener('click', () => {
            this.toggleSidebar();
        });
    }

    // 切换侧边栏状态 - 使用实例属性，无需参数
    toggleSidebar() {
        this.isSidebarOpen = !this.isSidebarOpen;
        // console.log(2, this.main, this.footer);
        
        if (this.isSidebarOpen) {
            // 显示侧边栏
            this.sidebar.classList.remove('sidebar-hidden');
            this.main?.classList.remove('main-expanded');
            this.footer?.classList.remove('footer-expanded');
            this.toggleButton.textContent = 'X';
            this.toggleButton.setAttribute('aria-expanded', 'true');
        } else {
            // 隐藏侧边栏
            this.sidebar.classList.add('sidebar-hidden');
            this.main?.classList.add('main-expanded');
            this.footer?.classList.add('footer-expanded');
            this.toggleButton.textContent = '≡';
            this.toggleButton.setAttribute('aria-expanded', 'false');
        }
    }

    async handleLogin(form) {
        const formData = new FormData(form);
        const username = formData.get('username');
        const password = formData.get('password');
        const salt = this.loginSalt;

        if (!username || !password) {
            alert('Please enter both username and password.');
            return;
        }

        const auth = await import('./auth.js');
        const authInstance = new auth.Auth();
        const success = await authInstance.login(username, password, salt);
        if (success) {
            this.SwithLoginState();
            this.modalManager.hide('loginModal');
            window.location.reload(); // 刷新页面以加载用户数据
        } else {
            alert('Login failed. Please check your credentials.');
        }
    }


}

document.addEventListener('DOMContentLoaded', () => {
    new App();
    // console.log('Page initialized');
});