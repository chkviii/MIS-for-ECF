// API Helper function - automatically add Authorization header
export async function apiRequest(url, options = {}) {
    const token = localStorage.getItem('token');
    
    if (!token) {
        window.location.href = '/login';
        throw new Error('No authentication token found');
    }
    
    // Set default headers
    const headers = {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${token}`,
        ...options.headers
    };
    
    // Merge options
    const requestOptions = {
        ...options,
        headers
    };
    
    try {
        const response = await fetch(url, requestOptions);
        
        // If unauthorized, redirect to login
        if (response.status === 401) {
            localStorage.removeItem('token');
            localStorage.removeItem('user');
            localStorage.removeItem('profile');
            window.location.href = '/login';
            throw new Error('Authentication failed');
        }
        
        return response;
    } catch (error) {
        throw error;
    }
}

// Check if user is logged in
function checkAuth() {
    const token = localStorage.getItem('token');
    const user = localStorage.getItem('user');
    
    if (!token || !user) {
        window.location.href = '/login';
        return false;
    }
    return true;
}

// Logout function
function logout() {
    localStorage.removeItem('token');
    localStorage.removeItem('user');
    localStorage.removeItem('profile');
    window.location.href = '/login';
}

// Add logout button event on page load
document.addEventListener('DOMContentLoaded', async function() {
    // Check authentication
    checkAuth();



    // Add logout button to header if it exists
    const header = document.querySelector('.header');
    if (header) {
        const user = JSON.parse(localStorage.getItem('user') || '{}');
        const profile = JSON.parse(localStorage.getItem('profile') || '{}');

        const userInfo = document.createElement('div');
        userInfo.className = 'user-info';
        userInfo.style.cssText = 'position: absolute; right: 20px; top: 50%; transform: translateY(-50%); display: flex; align-items: center; gap: 15px;';
        userInfo.innerHTML = `
            <span style="color: #666;">Welcome, ${profile.first_name || user.username}</span>
            <button id="logout-btn" class="btn btn-small btn-secondary">Logout</button>
        `;
        header.appendChild(userInfo);
    }

    const logoutBtn = document.getElementById('logout-btn');
    if (logoutBtn) {
        logoutBtn.addEventListener('click', function() {
            logout();
        }
    );
    }

    const sidebar = document.getElementById('sidebar');
    if (sidebar) {
        // Adjust main content margin based on sidebar width 
        // Initialize navigation
        initNavigation();
        // Initialize sidebar toggle
        initSidebarToggle();
    }
});



// Initialization
// (initial DOMContentLoaded handler already declared earlier; avoid duplicate initialization)

// Navigation Initialization
function initNavigation() {
    return; // later
}

// Sidebar Toggle Initialization
function initSidebarToggle() {
    const sidebar = document.getElementById('sidebar');
    const toggleBtn = document.getElementById('sidebar-toggle');
    
    // Load saved state from localStorage
    const sidebarState = localStorage.getItem('sidebarState') || 'normal';
    applySidebarState(sidebarState);
    
    // Initial adjustment
    updateMainContentMargin();
    
    toggleBtn.addEventListener('click', function() {
        let currentState = 'normal';
        
        // Toggle between normal and collapsed only
        if (sidebar.classList.contains('collapsed')) {
            currentState = 'normal';
        } else {
            currentState = 'collapsed';
        }
        
        applySidebarState(currentState);
        localStorage.setItem('sidebarState', currentState);
    });
    
    // Update main content margin when sidebar width changes
    const resizeObserver = new ResizeObserver(entries => {
        updateMainContentMargin();
    });
    
    resizeObserver.observe(sidebar);
    
    // Also update on window resize
    window.addEventListener('resize', () => {
        updateMainContentMargin();
    });
}

// Update main content margin based on sidebar width
function updateMainContentMargin() {
    const sidebar = document.getElementById('sidebar');
    const mainContent = document.getElementById('main-content');
    const header = document.querySelector('.header');
    const sidebarWidth = sidebar.offsetWidth;
    mainContent.style.marginLeft = `${sidebarWidth}px`;
    header.style.marginLeft = `${sidebarWidth}px`;
}

// Apply Sidebar State
function applySidebarState(state) {
    const sidebar = document.getElementById('sidebar');
    const mainContent = document.getElementById('main-content');
    
    // Remove all state classes
    sidebar.classList.remove('collapsed');
    mainContent.classList.remove('sidebar-collapsed');
    
    if (state === 'collapsed') {
        sidebar.classList.add('collapsed');
        mainContent.classList.add('sidebar-collapsed');
    }
    
    // Wait for CSS transition to complete, then update margin
    setTimeout(() => {
        updateMainContentMargin();
    }, 50);
}

// Show Toast
export function showToast(message, type = 'success') {
    const toast = document.getElementById('toast');
    const toastMessage = document.getElementById('toast-message');
    
    toastMessage.textContent = message;
    toast.className = `toast ${type} show`;
    
    setTimeout(() => {
        toast.classList.remove('show');
    }, 3000);
}