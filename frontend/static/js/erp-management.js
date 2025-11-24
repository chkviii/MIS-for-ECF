// ERP Management System JavaScript
const API_BASE_URL = '/api/v1';

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
document.addEventListener('DOMContentLoaded', function() {
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
    loadEntity('projects');
});

// Entity Configuration
const ENTITY_CONFIG = {
    'projects': {
        title: 'Project Management',
        endpoint: 'projects',
        fields: [
            { name: 'id', label: 'ID', type: 'number', readonly: true, showInTable: true, showInForm: false },
            { name: 'project_id', label: 'Project ID', type: 'text', readonly: true, showInTable: true, showInForm: 'edit' },
            { name: 'name', label: 'Project Name', type: 'text', required: true, showInTable: true, searchable: true },
            { name: 'description', label: 'Description', type: 'textarea', showInTable: false },
            { name: 'project_type', label: 'Project Type', type: 'text', showInTable: true, searchable: true },
            { name: 'budget', label: 'Budget', type: 'number', showInTable: true },
            { name: 'actual_cost', label: 'Actual Cost', type: 'number', showInTable: true },
            { name: 'location_id', label: 'Location ID', type: 'number', showInTable: true },
            { name: 'status', label: 'Status', type: 'select', options: ['planning', 'active', 'completed', 'cancelled'], showInTable: true, searchable: true },
            { name: 'start_date', label: 'Start Date', type: 'date', showInTable: true, searchable: true, dateRange: true },
            { name: 'end_date', label: 'End Date', type: 'date', showInTable: true }
        ]
    },
    'donors': {
        title: 'Donor Management',
        endpoint: 'donors',
        fields: [
            { name: 'id', label: 'ID', type: 'number', readonly: true, showInTable: true, showInForm: false },
            { name: 'donor_id', label: 'Donor ID', type: 'text', readonly: true, showInTable: true, showInForm: 'edit' },
            { name: 'first_name', label: 'First Name', type: 'text', required: true, showInTable: true, searchable: true },
            { name: 'last_name', label: 'Last Name', type: 'text', required: true, showInTable: true, searchable: true },
            { name: 'email', label: 'Email', type: 'email', showInTable: true, searchable: true },
            { name: 'phone', label: 'Phone', type: 'text', showInTable: true, searchable: true },
            { name: 'donor_type', label: 'Type', type: 'select', options: ['individual', 'corporate', 'foundation'], showInTable: true, searchable: true },
            { name: 'status', label: 'Status', type: 'select', options: ['active', 'inactive'], showInTable: true }
        ]
    },
    'donations': {
        title: 'Donation Records',
        endpoint: 'donations',
        fields: [
            { name: 'id', label: 'ID', type: 'number', readonly: true, showInTable: true, showInForm: false },
            { name: 'donation_id', label: 'Donation ID', type: 'text', readonly: true, showInTable: true, showInForm: 'edit' },
            { name: 'donor_id', label: 'Donor ID', type: 'number', required: true, showInTable: true },
            { name: 'project_id', label: 'Project ID', type: 'number', showInTable: true },
            { name: 'amount', label: 'Amount', type: 'number', required: true, showInTable: true },
            { name: 'donation_type', label: 'Donation Type', type: 'text', required: true, showInTable: true, searchable: true },
            { name: 'category', label: 'Category', type: 'text', required: true, showInTable: true, searchable: true },
            { name: 'donation_date', label: 'Donation Date', type: 'date', showInTable: true, searchable: true, dateRange: true }
        ]
    },
    'volunteers': {
        title: 'Volunteer Management',
        endpoint: 'volunteers',
        fields: [
            { name: 'id', label: 'ID', type: 'number', readonly: true, showInTable: true, showInForm: false },
            { name: 'volunteer_id', label: 'Volunteer ID', type: 'text', readonly: true, showInTable: true, showInForm: 'edit' },
            { name: 'first_name', label: 'First Name', type: 'text', required: true, showInTable: true, searchable: true },
            { name: 'last_name', label: 'Last Name', type: 'text', required: true, showInTable: true, searchable: true },
            { name: 'email', label: 'Email', type: 'email', showInTable: true, searchable: true },
            { name: 'phone', label: 'Phone', type: 'text', showInTable: true },
            { name: 'location_id', label: 'Location ID', type: 'number', showInTable: true },
            { name: 'skills', label: 'Skills', type: 'textarea', showInTable: false },
            { name: 'status', label: 'Status', type: 'select', options: ['active', 'inactive'], showInTable: true }
        ]
    },
    'employees': {
        title: 'Employee Management',
        endpoint: 'employees',
        fields: [
            { name: 'id', label: 'ID', type: 'number', readonly: true, showInTable: true, showInForm: false },
            { name: 'employee_id', label: 'Employee ID', type: 'text', readonly: true, showInTable: true, showInForm: 'edit' },
            { name: 'first_name', label: 'First Name', type: 'text', required: true, showInTable: true, searchable: true },
            { name: 'last_name', label: 'Last Name', type: 'text', required: true, showInTable: true, searchable: true },
            { name: 'email', label: 'Email', type: 'email', showInTable: true, searchable: true },
            { name: 'phone', label: 'Phone', type: 'text', showInTable: true },
            { name: 'position', label: 'Position', type: 'text', showInTable: true, searchable: true },
            { name: 'department', label: 'Department', type: 'text', showInTable: true, searchable: true },
            { name: 'salary', label: 'Salary', type: 'number', showInTable: true },
            { name: 'location_id', label: 'Location ID', type: 'number', showInTable: true },
            { name: 'status', label: 'Status', type: 'select', options: ['active', 'inactive'], showInTable: true }
        ]
    },
    'locations': {
        title: 'Location Management',
        endpoint: 'locations',
        fields: [
            { name: 'id', label: 'ID', type: 'number', readonly: true, showInTable: true, showInForm: false },
            { name: 'location_id', label: 'Location ID', type: 'text', readonly: true, showInTable: true, showInForm: 'edit' },
            { name: 'name', label: 'Name', type: 'text', required: true, showInTable: true, searchable: true },
            { name: 'type', label: 'Type', type: 'text', showInTable: true, searchable: true },
            { name: 'address', label: 'Address', type: 'textarea', showInTable: true },
            { name: 'country_code', label: 'Country Code', type: 'text', showInTable: true }
        ]
    },
    'funds': {
        title: 'Fund Management',
        endpoint: 'funds',
        fields: [
            { name: 'id', label: 'ID', type: 'number', readonly: true, showInTable: true, showInForm: false },
            { name: 'fund_id', label: 'Fund ID', type: 'text', readonly: true, showInTable: true, showInForm: 'edit' },
            { name: 'donor_id', label: 'Donor ID', type: 'number', showInTable: true },
            { name: 'project_id', label: 'Project ID', type: 'number', showInTable: true },
            { name: 'transaction_id', label: 'Transaction ID', type: 'number', showInTable: true },
            { name: 'name', label: 'Fund Name', type: 'text', required: true, showInTable: true, searchable: true },
            { name: 'fund_type', label: 'Type', type: 'text', required: true, showInTable: true, searchable: true },
            { name: 'total_amount', label: 'Total Amount', type: 'number', required: true, showInTable: true },
            { name: 'current_balance', label: 'Current Balance', type: 'number', showInTable: true },
            { name: 'status', label: 'Status', type: 'select', options: ['active', 'closed'], showInTable: true }
        ]
    },
    'expenses': {
        title: 'Expense Management',
        endpoint: 'expenses',
        fields: [
            { name: 'id', label: 'ID', type: 'number', readonly: true, showInTable: true, showInForm: false },
            { name: 'expense_id', label: 'Expense ID', type: 'text', readonly: true, showInTable: true, showInForm: 'edit' },
            { name: 'fund_id', label: 'Fund ID', type: 'number', required: true, showInTable: true },
            { name: 'project_id', label: 'Project ID', type: 'number', showInTable: true },
            { name: 'employee_id', label: 'Employee ID', type: 'number', showInTable: true },
            { name: 'transaction_id', label: 'Transaction ID', type: 'number', showInTable: true },
            { name: 'description', label: 'Description', type: 'textarea', required: true, showInTable: true },
            { name: 'amount', label: 'Amount', type: 'number', required: true, showInTable: true },
            { name: 'expense_date', label: 'Expense Date', type: 'date', showInTable: true, searchable: true, dateRange: true },
            { name: 'approval_status', label: 'Approval Status', type: 'select', options: ['pending', 'approved', 'rejected'], showInTable: true }
        ]
    },
    'transactions': {
        title: 'Transaction Records',
        endpoint: 'transactions',
        fields: [
            { name: 'id', label: 'ID', type: 'number', readonly: true, showInTable: true, showInForm: false },
            { name: 'transaction_id', label: 'Transaction ID', type: 'text', readonly: true, showInTable: true, showInForm: 'edit' },
            { name: 'type', label: 'Type', type: 'text', required: true, showInTable: true, searchable: true },
            { name: 'amount', label: 'Amount', type: 'number', required: true, showInTable: true },
            { name: 'from_entity', label: 'From', type: 'text', showInTable: true },
            { name: 'to_entity', label: 'To', type: 'text', showInTable: true },
            { name: 'transaction_date', label: 'Transaction Date', type: 'date', showInTable: true, searchable: true, dateRange: true }
        ]
    },
    'purchases': {
        title: 'Purchase Management',
        endpoint: 'purchases',
        fields: [
            { name: 'id', label: 'ID', type: 'number', readonly: true, showInTable: true, showInForm: false },
            { name: 'purchase_id', label: 'Purchase ID', type: 'text', readonly: true, showInTable: true, showInForm: 'edit' },
            { name: 'transaction_id', label: 'Transaction ID', type: 'number', showInTable: true },
            { name: 'total_spent', label: 'Total Amount', type: 'number', required: true, showInTable: true },
            { name: 'supplier_name', label: 'Supplier', type: 'text', showInTable: true, searchable: true },
            { name: 'purchase_date', label: 'Purchase Date', type: 'date', showInTable: true, searchable: true, dateRange: true },
            { name: 'description', label: 'Description', type: 'textarea', showInTable: false }
        ]
    },
    'payrolls': {
        title: 'Payroll Management',
        endpoint: 'payrolls',
        fields: [
            { name: 'id', label: 'ID', type: 'number', readonly: true, showInTable: true, showInForm: false },
            { name: 'transaction_id', label: 'Transaction ID', type: 'number', required: true, showInTable: true },
            { name: 'employee_id', label: 'Employee ID', type: 'number', required: true, showInTable: true },
            { name: 'amount', label: 'Amount', type: 'number', required: true, showInTable: true },
            { name: 'pay_date', label: 'Pay Date', type: 'date', required: true, showInTable: true, searchable: true, dateRange: true },
            { name: 'deductions', label: 'Deductions', type: 'number', showInTable: true },
            { name: 'bonuses', label: 'Bonuses', type: 'number', showInTable: true }
        ]
    },
    'inventories': {
        title: 'Inventory Management',
        endpoint: 'inventories',
        fields: [
            { name: 'id', label: 'ID', type: 'number', readonly: true, showInTable: true, showInForm: false },
            { name: 'inventory_id', label: 'Inventory ID', type: 'text', readonly: true, showInTable: true, showInForm: 'edit' },
            { name: 'name', label: 'Name', type: 'text', required: true, showInTable: true, searchable: true },
            { name: 'category', label: 'Category', type: 'text', showInTable: true, searchable: true },
            { name: 'purchase_id', label: 'Purchase ID', type: 'number', showInTable: true },
            { name: 'location_id', label: 'Location ID', type: 'number', showInTable: true },
            { name: 'current_stock', label: 'Current Stock', type: 'number', showInTable: true },
            { name: 'unit_cost', label: 'Unit Cost', type: 'number', showInTable: true },
            { name: 'status', label: 'Status', type: 'select', options: ['available', 'out_of_stock', 'discontinued'], showInTable: true }
        ]
    },
    'gift-types': {
        title: 'Gift Types',
        endpoint: 'gift-types',
        fields: [
            { name: 'id', label: 'ID', type: 'number', readonly: true, showInTable: true, showInForm: false },
            { name: 'name', label: 'Name', type: 'text', required: true, showInTable: true, searchable: true },
            { name: 'category', label: 'Category', type: 'text', showInTable: true },
            { name: 'unit_cost', label: 'Unit Cost', type: 'number', showInTable: true },
            { name: 'requires_inventory', label: 'Requires Inventory', type: 'checkbox', showInTable: true }
        ]
    },
    'gifts': {
        title: 'Gift Management',
        endpoint: 'gifts',
        fields: [
            { name: 'id', label: 'ID', type: 'number', readonly: true, showInTable: true, showInForm: false },
            { name: 'gift_id', label: 'Gift ID', type: 'text', readonly: true, showInTable: true, showInForm: 'edit' },
            { name: 'donation_id', label: 'Donation ID', type: 'number', showInTable: true },
            { name: 'delivery_id', label: 'Delivery ID', type: 'number', showInTable: true },
            { name: 'gift_type_id', label: 'Gift Type ID', type: 'number', required: true, showInTable: true },
            { name: 'quantity', label: 'Quantity', type: 'number', showInTable: true },
            { name: 'total_value', label: 'Total Value', type: 'number', showInTable: true },
            { name: 'distribution_status', label: 'Distribution Status', type: 'select', options: ['pending', 'shipped', 'delivered'], showInTable: true }
        ]
    },
    'inventory-transactions': {
        title: 'Inventory Transactions',
        endpoint: 'inventory-transactions',
        fields: [
            { name: 'id', label: 'ID', type: 'number', readonly: true, showInTable: true, showInForm: false },
            { name: 'transaction_id', label: 'Transaction ID', type: 'text', readonly: true, showInTable: true, showInForm: 'edit' },
            { name: 'inventory_id', label: 'Inventory ID', type: 'number', required: true, showInTable: true },
            { name: 'transaction_type', label: 'Transaction Type', type: 'text', required: true, showInTable: true, searchable: true },
            { name: 'quantity_change', label: 'Quantity Change', type: 'number', required: true, showInTable: true },
            { name: 'transaction_date', label: 'Transaction Date', type: 'date', required: true, showInTable: true, searchable: true, dateRange: true }
        ]
    },
    'deliveries': {
        title: 'Delivery Management',
        endpoint: 'deliveries',
        fields: [
            { name: 'id', label: 'ID', type: 'number', readonly: true, showInTable: true, showInForm: false },
            { name: 'delivery_id', label: 'Delivery ID', type: 'text', readonly: true, showInTable: true, showInForm: 'edit' },
            { name: 'inventory_id', label: 'Inventory ID', type: 'number', required: true, showInTable: true },
            { name: 'project_id', label: 'Project ID', type: 'number', showInTable: true },
            { name: 'location_id', label: 'Location ID', type: 'number', showInTable: true },
            { name: 'quantity', label: 'Quantity', type: 'number', required: true, showInTable: true },
            { name: 'recipient_name', label: 'Recipient', type: 'text', showInTable: true, searchable: true },
            { name: 'delivery_date', label: 'Delivery Date', type: 'date', showInTable: true, searchable: true, dateRange: true },
            { name: 'status', label: 'Status', type: 'select', options: ['pending', 'in_transit', 'delivered'], showInTable: true }
        ]
    },
    'volunteer-projects': {
        title: 'Volunteer-Project Assignments',
        endpoint: 'volunteer-projects',
        fields: [
            { name: 'id', label: 'ID', type: 'number', readonly: true, showInTable: true, showInForm: false },
            { name: 'volunteer_id', label: 'Volunteer ID', type: 'number', required: true, showInTable: true },
            { name: 'project_id', label: 'Project ID', type: 'number', required: true, showInTable: true },
            { name: 'role', label: 'Role', type: 'text', showInTable: true, searchable: true },
            { name: 'contract_start', label: 'Contract Start', type: 'date', showInTable: true, searchable: true, dateRange: true },
            { name: 'contract_end', label: 'Contract End', type: 'date', showInTable: true },
            { name: 'work_unit', label: 'Work Unit', type: 'text', showInTable: true },
            { name: 'total_amount', label: 'Total Amount', type: 'number', showInTable: true },
            { name: 'contract_date', label: 'Contract Date', type: 'date', showInTable: true },
            { name: 'contract_detail', label: 'Contract Detail', type: 'textarea', showInTable: false },
            { name: 'status', label: 'Status', type: 'select', options: ['active', 'completed', 'cancelled'], showInTable: true, searchable: true }
        ]
    },
    'employee-projects': {
        title: 'Employee-Project Assignments',
        endpoint: 'employee-projects',
        fields: [
            { name: 'id', label: 'ID', type: 'number', readonly: true, showInTable: true, showInForm: false },
            { name: 'employee_id', label: 'Employee ID', type: 'number', required: true, showInTable: true },
            { name: 'project_id', label: 'Project ID', type: 'number', required: true, showInTable: true },
            { name: 'title', label: 'Title', type: 'text', showInTable: true, searchable: true },
            { name: 'start_date', label: 'Start Date', type: 'date', showInTable: true, searchable: true, dateRange: true },
            { name: 'end_date', label: 'End Date', type: 'date', showInTable: true },
            { name: 'work_unit', label: 'Work Unit', type: 'text', showInTable: true },
            { name: 'allocated_amount', label: 'Allocated Amount', type: 'number', showInTable: true },
            { name: 'last_updated', label: 'Last Updated', type: 'datetime', readonly: true, showInTable: true, showInForm: 'edit' }
        ]
    },
    'fund-projects': {
        title: 'Fund-Project Allocations',
        endpoint: 'fund-projects',
        fields: [
            { name: 'id', label: 'ID', type: 'number', readonly: true, showInTable: true, showInForm: false },
            { name: 'transaction_id', label: 'Transaction ID', type: 'number', required: true, showInTable: true },
            { name: 'project_id', label: 'Project ID', type: 'number', required: true, showInTable: true },
            { name: 'fund_id', label: 'Fund ID', type: 'number', required: true, showInTable: true },
            { name: 'allocated_amount', label: 'Allocated Amount', type: 'number', required: true, showInTable: true },
            { name: 'allocation_date', label: 'Allocation Date', type: 'date', showInTable: true, searchable: true, dateRange: true },
            { name: 'purpose', label: 'Purpose', type: 'textarea', showInTable: false }
        ]
    },
    'donation-inventories': {
        title: 'Donation Inventory (In-Kind)',
        endpoint: 'donation-inventories',
        fields: [
            { name: 'id', label: 'ID', type: 'number', readonly: true, showInTable: true, showInForm: false },
            { name: 'donor_id', label: 'Donor ID', type: 'number', required: true, showInTable: true },
            { name: 'inventory_id', label: 'Inventory ID', type: 'number', required: true, showInTable: true },
            { name: 'donation_date', label: 'Donation Date', type: 'date', showInTable: true, searchable: true, dateRange: true },
            { name: 'project_id', label: 'Project ID', type: 'number', showInTable: true },
            { name: 'quantity', label: 'Quantity', type: 'number', showInTable: true },
            { name: 'estimated_value', label: 'Estimated Value', type: 'number', showInTable: true }
        ]
    },
    'schedules': {
        title: 'Schedule Management',
        endpoint: 'schedules',
        fields: [
            { name: 'id', label: 'ID', type: 'number', readonly: true, showInTable: true, showInForm: false },
            { name: 'schedule_id', label: 'Schedule ID', type: 'text', readonly: true, showInTable: true, showInForm: 'edit' },
            { name: 'person_id', label: 'Person ID', type: 'number', required: true, showInTable: true },
            { name: 'person_type', label: 'Person Type', type: 'select', options: ['volunteer', 'employee'], required: true, showInTable: true, searchable: true },
            { name: 'project_id', label: 'Project ID', type: 'number', showInTable: true },
            { name: 'shift_date', label: 'Shift Date', type: 'date', required: true, showInTable: true, searchable: true, dateRange: true },
            { name: 'start_time', label: 'Start Time', type: 'time', required: true, showInTable: true },
            { name: 'end_time', label: 'End Time', type: 'time', required: true, showInTable: true },
            { name: 'hours_worked', label: 'Hours Worked', type: 'number', showInTable: true },
            { name: 'status', label: 'Status', type: 'select', options: ['scheduled', 'completed', 'cancelled'], showInTable: true, searchable: true },
            { name: 'notes', label: 'Notes', type: 'textarea', showInTable: false }
        ]
    }
};

// Global State
let currentEntity = 'projects';
let currentData = [];
let editingItem = null;

// Initialization
// (initial DOMContentLoaded handler already declared earlier; avoid duplicate initialization)

// Navigation Initialization
function initNavigation() {
    const navItems = document.querySelectorAll('.nav-item');
    navItems.forEach(item => {
        item.addEventListener('click', function(e) {
            e.preventDefault();
            const entity = this.dataset.entity;
            
            navItems.forEach(nav => nav.classList.remove('active'));
            this.classList.add('active');
            
            loadEntity(entity);
        });
    });

    // Button Events
    document.getElementById('btn-add').addEventListener('click', () => openModal());
    document.getElementById('btn-search').addEventListener('click', () => searchData());
    document.getElementById('btn-reset').addEventListener('click', () => resetSearch());
    document.getElementById('btn-save').addEventListener('click', () => saveData());
    document.getElementById('btn-cancel').addEventListener('click', () => closeModal());
    document.getElementById('modal-close').addEventListener('click', () => closeModal());
}

// Sidebar Toggle Initialization
function initSidebarToggle() {
    const sidebar = document.getElementById('sidebar');
    const mainContent = document.getElementById('main-content');
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
    const sidebarWidth = sidebar.offsetWidth;
    mainContent.style.marginLeft = `${sidebarWidth}px`;
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

// Load Entity Data
async function loadEntity(entity) {
    currentEntity = entity;
    const config = ENTITY_CONFIG[entity];
    
    document.getElementById('page-title').textContent = config.title;
    
    // Generate Search Controls
    generateSearchControls(config);
    
    // Generate Table
    generateTable(config);
    
    // Fetch Data
    await fetchData();
}

// Generate Search Controls
function generateSearchControls(config) {
    const searchControls = document.getElementById('search-controls');
    searchControls.innerHTML = '';
    
    config.fields.filter(f => f.searchable).forEach(field => {
        const formGroup = document.createElement('div');
        formGroup.className = 'form-group';
        
        const label = document.createElement('label');
        label.textContent = field.label;
        formGroup.appendChild(label);
        
        if (field.type === 'select') {
            const select = document.createElement('select');
            select.id = `search-${field.name}`;
            select.innerHTML = '<option value="">All</option>';
            field.options.forEach(opt => {
                select.innerHTML += `<option value="${opt}">${opt}</option>`;
            });
            formGroup.appendChild(select);
        } else if (field.dateRange) {
            const startInput = document.createElement('input');
            startInput.type = 'date';
            startInput.id = `search-${field.name}-start`;
            startInput.placeholder = 'Start Date';
            formGroup.appendChild(startInput);
            
            const endInput = document.createElement('input');
            endInput.type = 'date';
            endInput.id = `search-${field.name}-end`;
            endInput.placeholder = 'End Date';
            formGroup.appendChild(endInput);
        } else {
            const input = document.createElement('input');
            input.type = 'text';
            input.id = `search-${field.name}`;
            input.placeholder = `Search ${field.label}`;
            formGroup.appendChild(input);
        }
        
        searchControls.appendChild(formGroup);
    });
}

// Generate Table
function generateTable(config) {
    const tableHead = document.getElementById('table-head');
    const headerRow = document.createElement('tr');
    
    config.fields.filter(f => f.showInTable).forEach(field => {
        const th = document.createElement('th');
        th.textContent = field.label;
        headerRow.appendChild(th);
    });
    
    const actionTh = document.createElement('th');
    actionTh.textContent = 'Actions';
    headerRow.appendChild(actionTh);
    
    tableHead.innerHTML = '';
    tableHead.appendChild(headerRow);
}

// Fetch Data
async function fetchData() {
    try {
        console.log("Fetching data for entity: ", currentEntity);
        const config = ENTITY_CONFIG[currentEntity];
        console.log("current config: ", config);
        const response = await apiRequest(`${API_BASE_URL}/${config.endpoint}`);
        console.log("current response: ", response);
        const result = await response.json();
        console.log("Response Json: ", result);
        
        currentData = result.data || [];
        renderTable();
        
        document.getElementById('data-count').textContent = `Total: ${currentData.length} records`;
    } catch (error) {
        showToast('Failed to load data: ' + error.message, 'error');
    }
}

// Render Table
function renderTable() {
    const config = ENTITY_CONFIG[currentEntity];
    const tableBody = document.getElementById('table-body');
    tableBody.innerHTML = '';
    
    currentData.forEach(item => {
        const row = document.createElement('tr');
        
        config.fields.filter(f => f.showInTable).forEach(field => {
            const td = document.createElement('td');
            td.textContent = item[field.name] || '';
            row.appendChild(td);
        });
        
        const actionTd = document.createElement('td');
        actionTd.innerHTML = `
            <div class="action-buttons">
                <button class="btn btn-small btn-secondary" onclick="editItem(${item.id})">Edit</button>
                <button class="btn btn-small btn-danger" onclick="deleteItem(${item.id})">Delete</button>
            </div>
        `;
        row.appendChild(actionTd);
        
        tableBody.appendChild(row);
    });
}

// Search Data
async function searchData() {
    try {
        const config = ENTITY_CONFIG[currentEntity];
        
        // Build query parameters from search inputs
        const queryParams = new URLSearchParams();
        
        config.fields.filter(f => f.searchable).forEach(field => {
            if (field.dateRange) {
                // Handle date range fields
                const startInput = document.getElementById(`search-${field.name}-start`);
                const endInput = document.getElementById(`search-${field.name}-end`);
                
                if (startInput && startInput.value) {
                    queryParams.append('start_date', startInput.value);
                }
                if (endInput && endInput.value) {
                    queryParams.append('end_date', endInput.value);
                }
            } else {
                // Handle regular fields
                const input = document.getElementById(`search-${field.name}`);
                if (input && input.value) {
                    queryParams.append(field.name, input.value);
                }
            }
        });
        
        // If no search params, just fetch all data
        if (queryParams.toString() === '') {
            await fetchData();
            return;
        }
        
        // Call search endpoint with query params
        const response = await apiRequest(`${API_BASE_URL}/${config.endpoint}/search?${queryParams.toString()}`);
        const result = await response.json();
        
        currentData = result.data || [];
        renderTable();
        
        document.getElementById('data-count').textContent = `Total: ${currentData.length} records`;
        showToast(`Found ${currentData.length} records`, 'success');
    } catch (error) {
        showToast('Search failed: ' + error.message, 'error');
    }
}

// Reset Search
function resetSearch() {
    const searchInputs = document.querySelectorAll('#search-controls input, #search-controls select');
    searchInputs.forEach(input => input.value = '');
    fetchData();
}

// Open Modal
function openModal(item = null) {
    editingItem = item;
    const config = ENTITY_CONFIG[currentEntity];
    const modal = document.getElementById('edit-modal');
    const form = document.getElementById('edit-form');
    
    document.getElementById('modal-title').textContent = item ? 'Edit ' + config.title : 'Add ' + config.title;
    
    form.innerHTML = '';
    config.fields.forEach(field => {
        // Check if field should be shown in form
        if (field.showInForm === false) return;
        if (field.showInForm === 'edit' && !item) return;
        
        const formGroup = document.createElement('div');
        formGroup.className = 'form-group';
        
        const label = document.createElement('label');
        label.textContent = field.label + (field.required ? ' *' : '');
        formGroup.appendChild(label);
        
        let input;
        if (field.type === 'textarea') {
            input = document.createElement('textarea');
        } else if (field.type === 'select') {
            input = document.createElement('select');
            input.innerHTML = '<option value="">Please select</option>';
            field.options.forEach(opt => {
                input.innerHTML += `<option value="${opt}">${opt}</option>`;
            });
        } else if (field.type === 'checkbox') {
            input = document.createElement('input');
            input.type = 'checkbox';
        } else {
            input = document.createElement('input');
            input.type = field.type;
        }
        
        input.id = `edit-${field.name}`;
        input.name = field.name;
        if (field.readonly) input.readOnly = true;
        if (field.required) input.required = true;
        
        if (item && item[field.name] !== undefined) {
            if (field.type === 'checkbox') {
                input.checked = item[field.name];
            } else {
                input.value = item[field.name];
            }
        }
        
        formGroup.appendChild(input);
        form.appendChild(formGroup);
    });
    
    modal.classList.add('show');
}

// Close Modal
function closeModal() {
    document.getElementById('edit-modal').classList.remove('show');
    editingItem = null;
}

// Save Data
async function saveData() {
    const config = ENTITY_CONFIG[currentEntity];
    const form = document.getElementById('edit-form');
    const formData = new FormData(form);
    
    const data = {};
    config.fields.forEach(field => {
        // Skip fields that are not in form
        if (field.showInForm === false) return;
        if (field.showInForm === 'edit' && !editingItem) return;
        
        const input = document.getElementById(`edit-${field.name}`);
        if (input) {
            if (field.type === 'checkbox') {
                data[field.name] = input.checked;
            } else if (field.type === 'number') {
                data[field.name] = input.value ? parseFloat(input.value) : null;
            } else if (field.type === 'date' && input.value) {
                // Convert date string to ISO 8601 format with local timezone offset
                data[field.name] = formatDateWithTimezone(input.value);
            } else {
                data[field.name] = input.value || null;
            }
        }
    });
    
    try {
        const url = editingItem 
            ? `${API_BASE_URL}/${config.endpoint}/${editingItem.id}`
            : `${API_BASE_URL}/${config.endpoint}`;
        
        const method = editingItem ? 'PUT' : 'POST';
        
        const response = await apiRequest(url, {
            method: method,
            body: JSON.stringify(data)
        });
        
        const result = await response.json();
        
        if (response.ok) {
            showToast(result.message || 'Saved successfully', 'success');
            closeModal();
            fetchData();
        } else {
            showToast(result.error || 'Save failed', 'error');
        }
    } catch (error) {
        showToast('Save failed: ' + error.message, 'error');
    }
}

// Format date with timezone offset
function formatDateWithTimezone(dateString) {
    // Create date object from input (YYYY-MM-DD)
    const date = new Date(dateString + 'T00:00:00');
    
    // Get timezone offset in minutes
    const timezoneOffset = -date.getTimezoneOffset();
    
    // Convert offset to hours and minutes
    const offsetHours = Math.floor(Math.abs(timezoneOffset) / 60);
    const offsetMinutes = Math.abs(timezoneOffset) % 60;
    
    // Format offset as +HH:MM or -HH:MM
    const offsetSign = timezoneOffset >= 0 ? '+' : '-';
    const offsetString = `${offsetSign}${String(offsetHours).padStart(2, '0')}:${String(offsetMinutes).padStart(2, '0')}`;
    
    // Format date components
    const year = date.getFullYear();
    const month = String(date.getMonth() + 1).padStart(2, '0');
    const day = String(date.getDate()).padStart(2, '0');
    const hours = String(date.getHours()).padStart(2, '0');
    const minutes = String(date.getMinutes()).padStart(2, '0');
    const seconds = String(date.getSeconds()).padStart(2, '0');
    
    // Return ISO 8601 format with timezone offset
    return `${year}-${month}-${day}T${hours}:${minutes}:${seconds}${offsetString}`;
}

// Edit Item
function editItem(id) {
    const item = currentData.find(i => i.id === id);
    if (item) {
        openModal(item);
    }
}

// Delete Item
async function deleteItem(id) {
    if (!confirm('Are you sure you want to delete this record?')) return;
    
    try {
        const config = ENTITY_CONFIG[currentEntity];
        const response = await apiRequest(`${API_BASE_URL}/${config.endpoint}/${id}`, {
            method: 'DELETE'
        });
        
        const result = await response.json();
        
        if (response.ok) {
            showToast(result.message || 'Deleted successfully', 'success');
            fetchData();
        } else {
            showToast(result.error || 'Delete failed', 'error');
        }
    } catch (error) {
        showToast('Delete failed: ' + error.message, 'error');
    }
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
