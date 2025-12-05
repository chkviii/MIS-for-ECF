// ERP Management System JavaScript
const API_BASE_URL = '/api/v1/dbms';

import { apiRequest, showToast } from '/static/js/site-api.js';


// Add logout button event on page load
document.addEventListener('DOMContentLoaded', async function() {

    // Ensure entity config is loaded before initializing UI that depends on it
    await ensureEntityConfig();

    // Initialize navigation and table view
    initTableView();

    const headertitle = document.getElementById('page-title');
    const initEntity = headertitle && headertitle.getAttribute('init');
    if (initEntity) {
        loadEntity(initEntity);
    }

});

// Lazy-load entity config from JSON when needed (no global/window binding)
let ENTITY_CONFIG = null;

async function ensureEntityConfig() {
    if (ENTITY_CONFIG) return;
    try {
        const resp = await fetch('/static/js/entity-config.json', { cache: 'no-store' });
        if (!resp.ok) throw new Error('Failed to fetch entity-config.json: ' + resp.status);
        ENTITY_CONFIG = await resp.json();
    } catch (err) {
        console.error('Unable to load entity configuration:', err);
        // Provide empty fallback to avoid runtime exceptions — UI will show errors when config missing
        ENTITY_CONFIG = {};
    }
}

// Global State
let currentEntity = '';
let currentData = [];
let editingItem = null;
// Currently-loaded entity config (defensive accessor stores last loaded config)
let currentConfig = null;

// Defensive accessor: ensure config JSON is loaded and return a safe config object
async function getEntityConfig(entity) {
    if (!ENTITY_CONFIG) await ensureEntityConfig();
    const cfg = ENTITY_CONFIG && ENTITY_CONFIG[entity];
    if (!cfg) {
        console.warn(`Entity config for "${entity}" not found, using fallback.`);
        return { title: entity, endpoint: entity, fields: [] };
    }
    return cfg;
}

// Initialization
// (initial DOMContentLoaded handler already declared earlier; avoid duplicate initialization)

// Navigation Initialization
function initTableView() {
    const navItems = document.querySelectorAll('.nav-item');
    navItems.forEach(item => {
        if (item.getAttribute('data-entity') !== null) {
            item.addEventListener('click', function(e) {
                e.preventDefault();
                const entity = this.dataset.entity;
                
                navItems.forEach(nav => nav.classList.remove('active'));
                this.classList.add('active');
                
                loadEntity(entity);
            });
        }
    });

    // Button Events
    document.getElementById('btn-add').addEventListener('click', () => openModal());
    document.getElementById('btn-search').addEventListener('click', () => searchData());
    document.getElementById('btn-reset').addEventListener('click', () => resetSearch());
    document.getElementById('btn-save').addEventListener('click', () => saveData());
    document.getElementById('btn-cancel').addEventListener('click', () => closeModal());
    document.getElementById('modal-close').addEventListener('click', () => closeModal());
}

// Load Entity Data
async function loadEntity(entity) {
    currentEntity = entity;
    currentConfig = await getEntityConfig(entity);

    document.getElementById('page-title').textContent = currentConfig.title || entity;

    // Generate Search Controls
    if (document.getElementById('search-controls')) {
        generateSearchControls(currentConfig);
    }

    // Generate Table
    generateTable(currentConfig);

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
        } else if (field.type === 'date') {
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
        } else if (field.type === 'number') {
            // Render min/max inputs for numeric range searches
            const minInput = document.createElement('input');
            minInput.type = 'number';
            minInput.id = `search-${field.name}-min`;
            minInput.placeholder = `${field.label} (min)`;
            formGroup.appendChild(minInput);

            const maxInput = document.createElement('input');
            maxInput.type = 'number';
            maxInput.id = `search-${field.name}-max`;
            maxInput.placeholder = `${field.label} (max)`;
            formGroup.appendChild(maxInput);
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

    const actionTh = document.createElement('th');
    actionTh.textContent = 'Actions';
    headerRow.appendChild(actionTh);

    config.fields.filter(f => f.showInTable).forEach(field => {
        const th = document.createElement('th');
        th.textContent = field.label;
        headerRow.appendChild(th);
    });
    
    
    
    
    tableHead.innerHTML = '';
    tableHead.appendChild(headerRow);
}

// Fetch Data
async function fetchData() {
    try {
        console.log("Fetching data for entity: ", currentEntity);
        const config = currentConfig || await getEntityConfig(currentEntity);
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
    const config = currentConfig || (ENTITY_CONFIG && ENTITY_CONFIG[currentEntity]) || { fields: [] };
    const tableBody = document.getElementById('table-body');
    tableBody.innerHTML = '';
    
    currentData.forEach(item => {

        const actionTd = document.createElement('td');
        const row = document.createElement('tr');

        actionTd.innerHTML = `
            <div class="action-buttons">
                <button class="btn btn-small btn-secondary" id="edit${item.id}">Edit</button>
                <button class="btn btn-small btn-danger" id="delete${item.id}">Delete</button>
            </div>
        `;
        actionTd.querySelector(`#edit${item.id}`).addEventListener('click', () => editItem(item.id));
        actionTd.querySelector(`#delete${item.id}`).addEventListener('click', () => deleteItem(item.id));

        row.appendChild(actionTd);

        config.fields.filter(f => f.showInTable).forEach(field => {
            const td = document.createElement('td');
            td.textContent = item[field.name] || '';
            row.appendChild(td);
        });
        
        
        
        tableBody.appendChild(row);
    });
}

// Search Data
async function searchData() {
    try {
        const config = currentConfig || (ENTITY_CONFIG && ENTITY_CONFIG[currentEntity]) || { fields: [] };
        // Build three maps: query, number_range, date_range
        const query = {};
        const number_range = {};
        const date_range = {};

        config.fields.filter(f => f.searchable).forEach(field => {
            if (field.type === 'date') {
                const startInput = document.getElementById(`search-${field.name}-start`);
                const endInput = document.getElementById(`search-${field.name}-end`);
                const start = startInput && startInput.value ? startInput.value : '';
                const end = endInput && endInput.value ? endInput.value : '';
                if (start || end) {
                    date_range[field.name] = [start, end];
                }
            } else if (field.type === 'number') {
                const minInput = document.getElementById(`search-${field.name}-min`);
                const maxInput = document.getElementById(`search-${field.name}-max`);
                const min = minInput && minInput.value ? minInput.value : '';
                const max = maxInput && maxInput.value ? maxInput.value : '';
                if (min || max) {
                    number_range[field.name] = [min, max];
                }
            } else {
                const input = document.getElementById(`search-${field.name}`);
                if (input && input.value) {
                    query[field.name] = input.value;
                }
            }
        });

        // If nothing provided, fetch all
        if (Object.keys(query).length === 0 && Object.keys(number_range).length === 0 && Object.keys(date_range).length === 0) {
            await fetchData();
            return;
        }

        // Encode maps as JSON in query string (server will decode)
        const qs = new URLSearchParams();
        if (Object.keys(query).length) qs.append('query', JSON.stringify(query));
        if (Object.keys(number_range).length) qs.append('number_range', JSON.stringify(number_range));
        if (Object.keys(date_range).length) qs.append('date_range', JSON.stringify(date_range));

        const response = await apiRequest(`${API_BASE_URL}/${config.endpoint}/search?` + qs.toString());
        const result = await response.json();
        if (!response.ok) {
            throw new Error('Search request failed: ' + result.error);
        }
        
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
    const config = currentConfig || (ENTITY_CONFIG && ENTITY_CONFIG[currentEntity]) || { fields: [], title: currentEntity };
    const modal = document.getElementById('edit-modal');
    const form = document.getElementById('edit-form');
    
    document.getElementById('modal-title').textContent = item ? 'Edit ' + config.title : 'Add ' + config.title;
    
    form.innerHTML = '';
    config.fields.forEach(field => {
        // Do not show timestamps in add/edit forms
        if (field.name === 'created_at' || field.name === 'updated_at') return;
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
            } else if (field.type === 'date') {
                // If the existing value is an ISO datetime, extract the date part for <input type="date">
                input.value = formatDateForInput(item[field.name]);
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
    const config = currentConfig || (ENTITY_CONFIG && ENTITY_CONFIG[currentEntity]) || { fields: [] };
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
    // Preserve the original created_at when editing so backend receives the original timestamp
    if (editingItem && editingItem.created_at !== undefined) {
        data.created_at = editingItem.created_at;
    }
    
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

// Format an ISO datetime or date string to YYYY-MM-DD for <input type="date">
function formatDateForInput(value) {
    if (!value) return '';
    try {
        // If already a plain YYYY-MM-DD, return as-is
        if (/^\d{4}-\d{2}-\d{2}$/.test(value)) return value;

        const d = new Date(value);
        if (isNaN(d.getTime())) return '';
        const year = d.getFullYear();
        const month = String(d.getMonth() + 1).padStart(2, '0');
        const day = String(d.getDate()).padStart(2, '0');
        return `${year}-${month}-${day}`;
    } catch (e) {
        return '';
    }
}

// Edit Item
export function editItem(id) {
    const item = currentData.find(i => i.id === id);
    if (item) {
        openModal(item);
    }
}

// Delete Item
export async function deleteItem(id) {
    if (!confirm('Are you sure you want to delete this record?')) return;
    
    try {
        const config = currentConfig || (ENTITY_CONFIG && ENTITY_CONFIG[currentEntity]) || { fields: [] };
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
// function showToast(message, type = 'success') {
//     const toast = document.getElementById('toast');
//     const toastMessage = document.getElementById('toast-message');
    
//     toastMessage.textContent = message;
//     toast.className = `toast ${type} show`;
    
//     setTimeout(() => {
//         toast.classList.remove('show');
//     }, 3000);
// }
