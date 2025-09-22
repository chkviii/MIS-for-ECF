class ModalManager {
    constructor() {
        this.init();
    }

    init() {
        this.bindEvents();
    }

    bindEvents() {
        // 关闭按钮事件
        document.addEventListener('click', (e) => {
            if (e.target.matches('.close[data-modal]')) {
                const modalId = e.target.getAttribute('data-modal');
                this.hide(modalId);
            }

            // 取消按钮事件
            if (e.target.matches('.btn[data-modal]') && 
                (e.target.textContent.includes('取消') || e.target.textContent.includes('Cancel'))) {
                const modalId = e.target.getAttribute('data-modal');
                this.hide(modalId);
            }

            // 背景点击关闭
            if (e.target.matches('.modal')) {
                this.hide(e.target.id);
            }
        });

        // ESC 键关闭
        document.addEventListener('keydown', (e) => {
            if (e.key === 'Escape') {
                const visibleModal = document.querySelector('.modal:not(.hidden)');
                if (visibleModal) {
                    this.hide(visibleModal.id);
                }
            }
        });
    }

    show(modalId) {
        const modal = document.getElementById(modalId);
        if (modal) {
            modal.classList.remove('hidden');
            document.body.style.overflow = 'hidden';
            // 聚焦到模态框内的第一个可聚焦元素
            const firstInput = modal.querySelector('input, button, textarea, select');
            if (firstInput) {
                setTimeout(() => firstInput.focus(), 100);
            }
        }
    }

    hide(modalId) {
        const modal = document.getElementById(modalId);
        if (modal) {
            modal.classList.add('hidden');
            document.body.style.overflow = '';
        }
    }

    toggle(modalId) {
        const modal = document.getElementById(modalId);
        if (modal) {
            if (modal.classList.contains('hidden')) {
                this.show(modalId);
            } else {
                this.hide(modalId);
            }
        }
    }
}

// 导出供其他模块使用
export { ModalManager };