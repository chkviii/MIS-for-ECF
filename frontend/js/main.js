class BlogApp {
    constructor() {
        this.currentUser = null;
        this.currentPage = 1;
        this.commentsPerPage = 10;
        this.init();
    }

    init() {
        this.bindEvents();
        this.checkLoginStatus();
        this.loadComments();
    }

    bindEvents() {
        // 模态框控制
        document.getElementById('loginBtn').addEventListener('click', () => this.openModal('loginModal'));
        document.getElementById('registerBtn').addEventListener('click', () => this.openModal('registerModal'));
        document.getElementById('logoutBtn').addEventListener('click', () => this.logout());

        // 关闭模态框
        document.querySelectorAll('.close, [data-modal]').forEach(btn => {
            btn.addEventListener('click', (e) => {
                if (e.target.dataset.modal) {
                    this.closeModal(e.target.dataset.modal);
                } else {
                    this.closeModal(e.target.closest('.modal').id);
                }
            });
        });

        // 表单提交
        document.getElementById('loginForm').addEventListener('submit', (e) => this.handleLogin(e));
        document.getElementById('registerForm').addEventListener('submit', (e) => this.handleRegister(e));
        document.getElementById('addCommentForm').addEventListener('submit', (e) => this.handleAddComment(e));

        // 分页
        document.getElementById('prevPage').addEventListener('click', () => this.prevPage());
        document.getElementById('nextPage').addEventListener('click', () => this.nextPage());

        // 点击模态框外部关闭
        document.querySelectorAll('.modal').forEach(modal => {
            modal.addEventListener('click', (e) => {
                if (e.target === modal) {
                    this.closeModal(modal.id);
                }
            });
        });
    }

    openModal(modalId) {
        document.getElementById(modalId).classList.remove('hidden');
    }

    closeModal(modalId) {
        document.getElementById(modalId).classList.add('hidden');
        // 清空表单
        const form = document.querySelector(`#${modalId} form`);
        if (form) form.reset();
    }

    showToast(message, isError = false) {
        const toast = document.getElementById('toast');
        const toastMessage = document.getElementById('toastMessage');
        
        toastMessage.textContent = message;
        toast.className = `toast ${isError ? 'error' : ''}`;
        toast.classList.remove('hidden');

        setTimeout(() => {
            toast.classList.add('hidden');
        }, 3000);
    }

    async handleLogin(e) {
        e.preventDefault();
        const formData = new FormData(e.target);
        const loginData = {
            username: formData.get('username'),
            password: formData.get('password')
        };

        try {
            const response = await fetch('/api/auth/login', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(loginData)
            });

            const result = await response.json();

            if (response.ok) {
                localStorage.setItem('token', result.token);
                this.currentUser = result.user;
                this.updateUserInterface();
                this.closeModal('loginModal');
                this.showToast('登录成功！');
            } else {
                this.showToast(result.error || '登录失败', true);
            }
        } catch (error) {
            console.error('登录错误:', error);
            this.showToast('网络错误，请稍后重试', true);
        }
    }

    async handleRegister(e) {
        e.preventDefault();
        const formData = new FormData(e.target);
        const password = formData.get('password');
        const confirmPassword = formData.get('confirmPassword');

        if (password !== confirmPassword) {
            this.showToast('两次输入的密码不一致', true);
            return;
        }

        const registerData = {
            username: formData.get('username'),
            password: password
        };

        try {
            const response = await fetch('/api/auth/register', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(registerData)
            });

            const result = await response.json();

            if (response.ok) {
                this.closeModal('registerModal');
                this.showToast('注册成功！请登录');
                this.openModal('loginModal');
            } else {
                this.showToast(result.error || '注册失败', true);
            }
        } catch (error) {
            console.error('注册错误:', error);
            this.showToast('网络错误，请稍后重试', true);
        }
    }

    async handleAddComment(e) {
        e.preventDefault();
        const formData = new FormData(e.target);
        const commentData = {
            content: formData.get('content'),
            article_id: 1 // 示例文章ID
        };

        try {
            const response = await fetch('/api/comments', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${localStorage.getItem('token')}`
                },
                body: JSON.stringify(commentData)
            });

            const result = await response.json();

            if (response.ok) {
                e.target.reset();
                this.showToast('评论发表成功！');
                this.loadComments();
            } else {
                this.showToast(result.error || '评论发表失败', true);
            }
        } catch (error) {
            console.error('评论错误:', error);
            this.showToast('网络错误，请稍后重试', true);
        }
    }

    logout() {
        localStorage.removeItem('token');
        this.currentUser = null;
        this.updateUserInterface();
        this.showToast('已退出登录');
    }

    checkLoginStatus() {
        const token = localStorage.getItem('token');
        if (token) {
            // 验证token有效性
            this.verifyToken(token);
        }
    }

    async verifyToken(token) {
        try {
            const response = await fetch('/api/auth/verify', {
                headers: {
                    'Authorization': `Bearer ${token}`
                }
            });

            if (response.ok) {
                const result = await response.json();
                this.currentUser = result.user;
                this.updateUserInterface();
            } else {
                localStorage.removeItem('token');
            }
        } catch (error) {
            console.error('验证token错误:', error);
            localStorage.removeItem('token');
        }
    }

    updateUserInterface() {
        const loginBtn = document.getElementById('loginBtn');
        const registerBtn = document.getElementById('registerBtn');
        const userMenu = document.getElementById('userMenu');
        const username = document.getElementById('username');
        const commentForm = document.getElementById('commentForm');

        if (this.currentUser) {
            loginBtn.classList.add('hidden');
            registerBtn.classList.add('hidden');
            userMenu.classList.remove('hidden');
            username.textContent = this.currentUser.username;
            commentForm.classList.remove('hidden');
        } else {
            loginBtn.classList.remove('hidden');
            registerBtn.classList.remove('hidden');
            userMenu.classList.add('hidden');
            commentForm.classList.add('hidden');
        }
    }

    async loadComments() {
        const commentsList = document.getElementById('commentsList');
        commentsList.innerHTML = '<div class="loading">加载评论中...</div>';

        try {
            const response = await fetch(`/api/comments?page=${this.currentPage}&limit=${this.commentsPerPage}&article_id=1`);
            const result = await response.json();

            if (response.ok) {
                this.renderComments(result.comments);
                this.updatePagination(result.total, result.page, result.totalPages);
            } else {
                commentsList.innerHTML = '<div class="loading">加载评论失败</div>';
            }
        } catch (error) {
            console.error('加载评论错误:', error);
            commentsList.innerHTML = '<div class="loading">网络错误，请稍后重试</div>';
        }
    }

    renderComments(comments) {
        const commentsList = document.getElementById('commentsList');
        
        if (comments.length === 0) {
            commentsList.innerHTML = '<div class="loading">暂无评论，快来发表第一条评论吧！</div>';
            return;
        }

        const commentsHtml = comments.map(comment => `
            <div class="comment">
                <div class="comment-header">
                    <span class="comment-author">${comment.username}</span>
                    <span class="comment-date">${new Date(comment.created_at).toLocaleString('zh-CN')}</span>
                </div>
                <div class="comment-content">${comment.content}</div>
            </div>
        `).join('');

        commentsList.innerHTML = commentsHtml;
    }

    updatePagination(total, currentPage, totalPages) {
        const pagination = document.getElementById('pagination');
        const prevBtn = document.getElementById('prevPage');
        const nextBtn = document.getElementById('nextPage');
        const pageInfo = document.getElementById('pageInfo');

        if (totalPages <= 1) {
            pagination.classList.add('hidden');
            return;
        }

        pagination.classList.remove('hidden');
        pageInfo.textContent = `第 ${currentPage} 页，共 ${totalPages} 页`;
        
        prevBtn.disabled = currentPage <= 1;
        nextBtn.disabled = currentPage >= totalPages;
    }

    prevPage() {
        if (this.currentPage > 1) {
            this.currentPage--;
            this.loadComments();
        }
    }

    nextPage() {
        this.currentPage++;
        this.loadComments();
    }
}

// 启动应用
document.addEventListener('DOMContentLoaded', () => {
    new BlogApp();
});
