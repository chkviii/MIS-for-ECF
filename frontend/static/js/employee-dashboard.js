// ERP Management System JavaScript
const API_BASE_URL = '/api/v1';

let ReportModulePromise = null;
let chart = null;

// API Helper function - automatically add Authorization header
async function apiRequest(url, options = {}) {
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
            <button onclick="logout()" class="btn btn-small btn-secondary">Logout</button>
        `;
        header.appendChild(userInfo);
    }

    

    initNavigation();
    initSidebarToggle();

    isAnyReportNavItemActive();
    addReportNavEventListeners();
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
function showToast(message, type = 'success') {
    const toast = document.getElementById('toast');
    const toastMessage = document.getElementById('toast-message');
    
    toastMessage.textContent = message;
    toast.className = `toast ${type} show`;
    
    setTimeout(() => {
        toast.classList.remove('show');
    }, 3000);
}



// Find if any report nav item in sidebar is active
function isAnyReportNavItemActive() {
    const reportContainer = document.getElementById('report-container');
    if (reportContainer) reportContainer.style.display = 'none';

    const reportNavItems = document.querySelectorAll('.nav-item[data-report]');
    if (!reportNavItems || reportNavItems.length === 0) {
        // no report nav items -> keep hidden
        return;
    }

    // Preload report.js asynchronously and cache the promise
    if (!ReportModulePromise) {
        ReportModulePromise = import('/static/js/report.js')
            .catch(err => {
                console.error('Failed to preload report module:', err);
                // clear cache so future attempts can retry
                ReportModulePromise = null;
            });
    }
}

// Add report event listeners
function addReportNavEventListeners() {
    const reportNavItems = document.querySelectorAll('.nav-item[data-report]');
    // console.log('Report nav items:', reportNavItems);
    if (!reportNavItems || reportNavItems.length === 0) return;

    reportNavItems.forEach(item => {
        item.addEventListener('click', async function(e) {
            const mainContent = document.getElementById('main-content');
            for (const child of mainContent.children) {
                child.style.display = 'none';
            }
            const reportContainer = document.getElementById('report-container');
            if (reportContainer) reportContainer.style.display = 'block';

            try {
                // ensure module loaded (use cached promise if present)
                const mod = ReportModulePromise ? await ReportModulePromise : await import('/static/js/report.js');
                const ReportChart = mod && (mod.default || mod.ReportChart);
                if (!ReportChart) throw new Error('ReportChart class not found in module');

                // reuse a single instance to avoid duplicates; destroy previous if exists
                if (window._currentReportChart) {
                    try { window._currentReportChart.destroy(); } catch(e){/*ignore*/ }
                    window._currentReportChart = null;
                }

                // instantiate and update
                window._currentReportChart = new ReportChart();
                await window._currentReportChart.update();
            } catch (err) {
                console.error('Failed to load or initialize report chart:', err);
                if (typeof showToast === 'function') showToast('无法加载报表: ' + (err.message || err), 'error');
            }
        });
    });
}