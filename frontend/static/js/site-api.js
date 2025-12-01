// Lightweight API helper used by simple pages (donor/volunteer/employee)
const API_BASE_URL = '/api/v1';

async function apiRequestSimple(url, options = {}) {
    const token = localStorage.getItem('token');
    if (!token) {
        window.location.href = '/login';
        throw new Error('No authentication token found');
    }
    const headers = {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${token}`,
        ...options.headers
    };
    const requestOptions = {
        ...options,
        headers,
        credentials: 'include'
    };
    const resp = await fetch(url, requestOptions);
    if (resp.status === 401) {
        localStorage.removeItem('token');
        localStorage.removeItem('user');
        window.location.href = '/login';
        throw new Error('Unauthorized');
    }
    return resp;
}

function checkAuthSimple() {
    const token = localStorage.getItem('token');
    const user = localStorage.getItem('user');
    if (!token || !user) {
        window.location.href = '/login';
        return false;
    }
    return true;
}

function logoutSimple() {
    localStorage.removeItem('token');
    localStorage.removeItem('user');
    window.location.href = '/login';
}

function decodeJWT(token) {
    try {
        const parts = token.split('.');
        if (parts.length !== 3) return null;
        const payload = parts[1];
        const decoded = atob(payload.replace(/-/g, '+').replace(/_/g, '/'));
        return JSON.parse(decodeURIComponent(escape(decoded)));
    } catch (e) {
        return null;
    }
}

export { apiRequestSimple, checkAuthSimple, logoutSimple, decodeJWT, API_BASE_URL };
