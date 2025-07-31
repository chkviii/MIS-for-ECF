class App{
    constructor() {
        this.isSidebarOpen = true;
        this.init();
    }

    init() {
        this.initSidebar();
    }

    initSidebar() {
        this.toggleButton = document.getElementById('toggleSidebar');
        this.sidebar = document.getElementById('sidebar');
        this.main = document.querySelector('main');
        this.footer = document.querySelector('footer');

        console.log(1, this.main, this.footer);

        
        if (!this.toggleButton || !this.sidebar) {
            console.warn('Sidebar elements not found');
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
        console.log(2, this.main, this.footer);
        
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
}

document.addEventListener('DOMContentLoaded', () => {
    new App();
    console.log('App initialized');
});