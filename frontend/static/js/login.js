document.addEventListener('DOMContentLoaded', function() {
    const loginForm = document.getElementById('login-form');
    const messageDiv = document.getElementById('message');

    loginForm.addEventListener('submit', async function(e) {
        e.preventDefault();

        const username = document.getElementById('username').value;
        const password = document.getElementById('password').value;

        try {
            const response = await fetch('/api/v1/auth/login', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    username: username,
                    password: password
                })
            });

            const data = await response.json();

            if (data.success) {
                // Save token and user info to localStorage
                localStorage.setItem('token', data.data.token);
                localStorage.setItem('user', JSON.stringify({
                    username: username,
                    user_type: data.data.user_type,
                    user_id: data.data.user_id
                }));

                messageDiv.textContent = 'Login successful! Redirecting...';
                messageDiv.className = 'message success';
                messageDiv.style.display = 'block';

                // Redirect based on user type
                setTimeout(() => {
                    const ut = data.data.user_type || '';
                    if (ut === 'donor') {
                        window.location.href = '/donor';
                    } else if (ut === 'volunteer') {
                        window.location.href = '/volunteer';
                    } else if (ut === 'employee') {
                        window.location.href = '/employee';
                    } else {
                        window.location.href = '/erp-management';
                    }
                }, 800);
            } else {
                messageDiv.textContent = data.message || 'Login failed';
                messageDiv.className = 'message error';
                messageDiv.style.display = 'block';
            }
        } catch (error) {
            console.error('Login error:', error);
            messageDiv.textContent = 'An error occurred. Please try again.';
            messageDiv.className = 'message error';
            messageDiv.style.display = 'block';
        }
    });
});
