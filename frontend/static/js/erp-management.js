// ERP管理系统JavaScript
const API_BASE_URL = '/api/v1';

// 实体配置
const ENTITY_CONFIG = {
    'projects': {
        title: '项目管理',
        endpoint: 'projects',
        fields: [
            { name: 'id', label: 'ID', type: 'number', readonly: true, showInTable: true },
            { name: 'project_id', label: '项目编号', type: 'text', readonly: true, showInTable: true },
            { name: 'name', label: '项目名称', type: 'text', required: true, showInTable: true, searchable: true },
            { name: 'description', label: '描述', type: 'textarea', showInTable: false },
            { name: 'project_type', label: '项目类型', type: 'text', showInTable: true, searchable: true },
            { name: 'budget', label: '预算', type: 'number', showInTable: true },
            { name: 'actual_cost', label: '实际成本', type: 'number', showInTable: true },
            { name: 'status', label: '状态', type: 'select', options: ['planning', 'active', 'completed', 'cancelled'], showInTable: true, searchable: true },
            { name: 'start_date', label: '开始日期', type: 'date', showInTable: true, searchable: true, dateRange: true },
            { name: 'end_date', label: '结束日期', type: 'date', showInTable: true }
        ]
    },
    'donors': {
        title: '捐赠者管理',
        endpoint: 'donors',
        fields: [
            { name: 'id', label: 'ID', type: 'number', readonly: true, showInTable: true },
            { name: 'donor_id', label: '捐赠者编号', type: 'text', readonly: true, showInTable: true },
            { name: 'first_name', label: '名', type: 'text', required: true, showInTable: true, searchable: true },
            { name: 'last_name', label: '姓', type: 'text', required: true, showInTable: true, searchable: true },
            { name: 'email', label: '邮箱', type: 'email', showInTable: true, searchable: true },
            { name: 'phone', label: '电话', type: 'text', showInTable: true, searchable: true },
            { name: 'donor_type', label: '类型', type: 'select', options: ['individual', 'corporate', 'foundation'], showInTable: true, searchable: true },
            { name: 'status', label: '状态', type: 'select', options: ['active', 'inactive'], showInTable: true }
        ]
    },
    'donations': {
        title: '捐赠记录',
        endpoint: 'donations',
        fields: [
            { name: 'id', label: 'ID', type: 'number', readonly: true, showInTable: true },
            { name: 'donation_id', label: '捐赠编号', type: 'text', readonly: true, showInTable: true },
            { name: 'donor_id', label: '捐赠者ID', type: 'number', required: true, showInTable: true },
            { name: 'amount', label: '金额', type: 'number', required: true, showInTable: true },
            { name: 'donation_type', label: '捐赠类型', type: 'text', required: true, showInTable: true, searchable: true },
            { name: 'category', label: '类别', type: 'text', required: true, showInTable: true, searchable: true },
            { name: 'donation_date', label: '捐赠日期', type: 'date', showInTable: true, searchable: true, dateRange: true }
        ]
    },
    'volunteers': {
        title: '志愿者管理',
        endpoint: 'volunteers',
        fields: [
            { name: 'id', label: 'ID', type: 'number', readonly: true, showInTable: true },
            { name: 'volunteer_id', label: '志愿者编号', type: 'text', readonly: true, showInTable: true },
            { name: 'first_name', label: '名', type: 'text', required: true, showInTable: true, searchable: true },
            { name: 'last_name', label: '姓', type: 'text', required: true, showInTable: true, searchable: true },
            { name: 'email', label: '邮箱', type: 'email', showInTable: true, searchable: true },
            { name: 'phone', label: '电话', type: 'text', showInTable: true },
            { name: 'skills', label: '技能', type: 'textarea', showInTable: false },
            { name: 'status', label: '状态', type: 'select', options: ['active', 'inactive'], showInTable: true }
        ]
    },
    'employees': {
        title: '员工管理',
        endpoint: 'employees',
        fields: [
            { name: 'id', label: 'ID', type: 'number', readonly: true, showInTable: true },
            { name: 'employee_id', label: '员工编号', type: 'text', readonly: true, showInTable: true },
            { name: 'first_name', label: '名', type: 'text', required: true, showInTable: true, searchable: true },
            { name: 'last_name', label: '姓', type: 'text', required: true, showInTable: true, searchable: true },
            { name: 'email', label: '邮箱', type: 'email', showInTable: true, searchable: true },
            { name: 'phone', label: '电话', type: 'text', showInTable: true },
            { name: 'position', label: '职位', type: 'text', showInTable: true, searchable: true },
            { name: 'department', label: '部门', type: 'text', showInTable: true, searchable: true },
            { name: 'salary', label: '薪资', type: 'number', showInTable: true },
            { name: 'status', label: '状态', type: 'select', options: ['active', 'inactive'], showInTable: true }
        ]
    },
    'locations': {
        title: '地点管理',
        endpoint: 'locations',
        fields: [
            { name: 'id', label: 'ID', type: 'number', readonly: true, showInTable: true },
            { name: 'location_id', label: '地点编号', type: 'text', readonly: true, showInTable: true },
            { name: 'name', label: '名称', type: 'text', required: true, showInTable: true, searchable: true },
            { name: 'type', label: '类型', type: 'text', showInTable: true, searchable: true },
            { name: 'address', label: '地址', type: 'textarea', showInTable: true },
            { name: 'country_code', label: '国家代码', type: 'text', showInTable: true }
        ]
    },
    'funds': {
        title: '基金管理',
        endpoint: 'funds',
        fields: [
            { name: 'id', label: 'ID', type: 'number', readonly: true, showInTable: true },
            { name: 'fund_id', label: '基金编号', type: 'text', readonly: true, showInTable: true },
            { name: 'name', label: '基金名称', type: 'text', required: true, showInTable: true, searchable: true },
            { name: 'fund_type', label: '类型', type: 'text', required: true, showInTable: true, searchable: true },
            { name: 'total_amount', label: '总金额', type: 'number', required: true, showInTable: true },
            { name: 'current_balance', label: '当前余额', type: 'number', showInTable: true },
            { name: 'status', label: '状态', type: 'select', options: ['active', 'closed'], showInTable: true }
        ]
    },
    'expenses': {
        title: '支出管理',
        endpoint: 'expenses',
        fields: [
            { name: 'id', label: 'ID', type: 'number', readonly: true, showInTable: true },
            { name: 'expense_id', label: '支出编号', type: 'text', readonly: true, showInTable: true },
            { name: 'fund_id', label: '基金ID', type: 'number', required: true, showInTable: true },
            { name: 'description', label: '描述', type: 'textarea', required: true, showInTable: true },
            { name: 'amount', label: '金额', type: 'number', required: true, showInTable: true },
            { name: 'expense_date', label: '支出日期', type: 'date', showInTable: true, searchable: true, dateRange: true },
            { name: 'approval_status', label: '审批状态', type: 'select', options: ['pending', 'approved', 'rejected'], showInTable: true }
        ]
    },
    'transactions': {
        title: '交易记录',
        endpoint: 'transactions',
        fields: [
            { name: 'id', label: 'ID', type: 'number', readonly: true, showInTable: true },
            { name: 'transaction_id', label: '交易编号', type: 'text', readonly: true, showInTable: true },
            { name: 'type', label: '类型', type: 'text', required: true, showInTable: true, searchable: true },
            { name: 'amount', label: '金额', type: 'number', required: true, showInTable: true },
            { name: 'from_entity', label: '来源', type: 'text', showInTable: true },
            { name: 'to_entity', label: '目标', type: 'text', showInTable: true },
            { name: 'transaction_date', label: '交易日期', type: 'date', showInTable: true, searchable: true, dateRange: true }
        ]
    },
    'purchases': {
        title: '采购管理',
        endpoint: 'purchases',
        fields: [
            { name: 'id', label: 'ID', type: 'number', readonly: true, showInTable: true },
            { name: 'purchase_id', label: '采购编号', type: 'text', readonly: true, showInTable: true },
            { name: 'total_spent', label: '总金额', type: 'number', required: true, showInTable: true },
            { name: 'supplier_name', label: '供应商', type: 'text', showInTable: true, searchable: true },
            { name: 'purchase_date', label: '采购日期', type: 'date', showInTable: true, searchable: true, dateRange: true },
            { name: 'description', label: '描述', type: 'textarea', showInTable: false }
        ]
    },
    'payrolls': {
        title: '薪资管理',
        endpoint: 'payrolls',
        fields: [
            { name: 'id', label: 'ID', type: 'number', readonly: true, showInTable: true },
            { name: 'employee_id', label: '员工ID', type: 'number', required: true, showInTable: true },
            { name: 'amount', label: '金额', type: 'number', required: true, showInTable: true },
            { name: 'pay_date', label: '支付日期', type: 'date', showInTable: true, searchable: true, dateRange: true },
            { name: 'deductions', label: '扣除', type: 'number', showInTable: true },
            { name: 'bonuses', label: '奖金', type: 'number', showInTable: true }
        ]
    },
    'inventory': {
        title: '库存管理',
        endpoint: 'inventory',
        fields: [
            { name: 'id', label: 'ID', type: 'number', readonly: true, showInTable: true },
            { name: 'inventory_id', label: '库存编号', type: 'text', readonly: true, showInTable: true },
            { name: 'name', label: '名称', type: 'text', required: true, showInTable: true, searchable: true },
            { name: 'category', label: '类别', type: 'text', showInTable: true, searchable: true },
            { name: 'current_stock', label: '当前库存', type: 'number', showInTable: true },
            { name: 'unit_cost', label: '单价', type: 'number', showInTable: true },
            { name: 'status', label: '状态', type: 'select', options: ['available', 'out_of_stock', 'discontinued'], showInTable: true }
        ]
    },
    'gift-types': {
        title: '礼品类型',
        endpoint: 'gift-types',
        fields: [
            { name: 'id', label: 'ID', type: 'number', readonly: true, showInTable: true },
            { name: 'name', label: '名称', type: 'text', required: true, showInTable: true, searchable: true },
            { name: 'category', label: '类别', type: 'text', showInTable: true },
            { name: 'unit_cost', label: '单价', type: 'number', showInTable: true },
            { name: 'requires_inventory', label: '需要库存', type: 'checkbox', showInTable: true }
        ]
    },
    'gifts': {
        title: '礼品管理',
        endpoint: 'gifts',
        fields: [
            { name: 'id', label: 'ID', type: 'number', readonly: true, showInTable: true },
            { name: 'gift_id', label: '礼品编号', type: 'text', readonly: true, showInTable: true },
            { name: 'gift_type_id', label: '礼品类型ID', type: 'number', required: true, showInTable: true },
            { name: 'quantity', label: '数量', type: 'number', showInTable: true },
            { name: 'total_value', label: '总价值', type: 'number', showInTable: true },
            { name: 'distribution_status', label: '配送状态', type: 'select', options: ['pending', 'shipped', 'delivered'], showInTable: true }
        ]
    },
    'inventory-transactions': {
        title: '库存交易',
        endpoint: 'inventory-transactions',
        fields: [
            { name: 'id', label: 'ID', type: 'number', readonly: true, showInTable: true },
            { name: 'transaction_id', label: '交易编号', type: 'text', readonly: true, showInTable: true },
            { name: 'inventory_id', label: '库存ID', type: 'number', required: true, showInTable: true },
            { name: 'transaction_type', label: '交易类型', type: 'text', required: true, showInTable: true, searchable: true },
            { name: 'quantity_change', label: '数量变化', type: 'number', required: true, showInTable: true },
            { name: 'transaction_date', label: '交易日期', type: 'date', showInTable: true, searchable: true, dateRange: true }
        ]
    },
    'deliveries': {
        title: '配送管理',
        endpoint: 'deliveries',
        fields: [
            { name: 'id', label: 'ID', type: 'number', readonly: true, showInTable: true },
            { name: 'delivery_id', label: '配送编号', type: 'text', readonly: true, showInTable: true },
            { name: 'inventory_id', label: '库存ID', type: 'number', required: true, showInTable: true },
            { name: 'quantity', label: '数量', type: 'number', required: true, showInTable: true },
            { name: 'recipient_name', label: '收件人', type: 'text', showInTable: true, searchable: true },
            { name: 'delivery_date', label: '配送日期', type: 'date', showInTable: true, searchable: true, dateRange: true },
            { name: 'status', label: '状态', type: 'select', options: ['pending', 'in_transit', 'delivered'], showInTable: true }
        ]
    }
};

// 全局状态
let currentEntity = 'projects';
let currentData = [];
let editingItem = null;

// 初始化
document.addEventListener('DOMContentLoaded', function() {
    initNavigation();
    loadEntity('projects');
});

// 导航初始化
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

    // 按钮事件
    document.getElementById('btn-add').addEventListener('click', () => openModal());
    document.getElementById('btn-search').addEventListener('click', () => searchData());
    document.getElementById('btn-reset').addEventListener('click', () => resetSearch());
    document.getElementById('btn-save').addEventListener('click', () => saveData());
    document.getElementById('btn-cancel').addEventListener('click', () => closeModal());
    document.getElementById('modal-close').addEventListener('click', () => closeModal());
}

// 加载实体数据
async function loadEntity(entity) {
    currentEntity = entity;
    const config = ENTITY_CONFIG[entity];
    
    document.getElementById('page-title').textContent = config.title;
    
    // 生成搜索控件
    generateSearchControls(config);
    
    // 生成表格
    generateTable(config);
    
    // 加载数据
    await fetchData();
}

// 生成搜索控件
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
            select.innerHTML = '<option value="">全部</option>';
            field.options.forEach(opt => {
                select.innerHTML += `<option value="${opt}">${opt}</option>`;
            });
            formGroup.appendChild(select);
        } else if (field.dateRange) {
            const startInput = document.createElement('input');
            startInput.type = 'date';
            startInput.id = `search-${field.name}-start`;
            startInput.placeholder = '开始日期';
            formGroup.appendChild(startInput);
            
            const endInput = document.createElement('input');
            endInput.type = 'date';
            endInput.id = `search-${field.name}-end`;
            endInput.placeholder = '结束日期';
            formGroup.appendChild(endInput);
        } else {
            const input = document.createElement('input');
            input.type = 'text';
            input.id = `search-${field.name}`;
            input.placeholder = `搜索${field.label}`;
            formGroup.appendChild(input);
        }
        
        searchControls.appendChild(formGroup);
    });
}

// 生成表格
function generateTable(config) {
    const tableHead = document.getElementById('table-head');
    const headerRow = document.createElement('tr');
    
    config.fields.filter(f => f.showInTable).forEach(field => {
        const th = document.createElement('th');
        th.textContent = field.label;
        headerRow.appendChild(th);
    });
    
    const actionTh = document.createElement('th');
    actionTh.textContent = '操作';
    headerRow.appendChild(actionTh);
    
    tableHead.innerHTML = '';
    tableHead.appendChild(headerRow);
}

// 获取数据
async function fetchData() {
    try {
        const config = ENTITY_CONFIG[currentEntity];
        const response = await fetch(`${API_BASE_URL}/${config.endpoint}`);
        const result = await response.json();
        
        currentData = result.data || [];
        renderTable();
        
        document.getElementById('data-count').textContent = `总计: ${currentData.length} 条`;
    } catch (error) {
        showToast('加载数据失败: ' + error.message, 'error');
    }
}

// 渲染表格
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
                <button class="btn btn-small btn-secondary" onclick="editItem(${item.id})">编辑</button>
                <button class="btn btn-small btn-danger" onclick="deleteItem(${item.id})">删除</button>
            </div>
        `;
        row.appendChild(actionTd);
        
        tableBody.appendChild(row);
    });
}

// 搜索数据
function searchData() {
    // 简单的前端过滤，实际应该后端实现
    fetchData();
}

// 重置搜索
function resetSearch() {
    const searchInputs = document.querySelectorAll('#search-controls input, #search-controls select');
    searchInputs.forEach(input => input.value = '');
    fetchData();
}

// 打开模态框
function openModal(item = null) {
    editingItem = item;
    const config = ENTITY_CONFIG[currentEntity];
    const modal = document.getElementById('edit-modal');
    const form = document.getElementById('edit-form');
    
    document.getElementById('modal-title').textContent = item ? '编辑' + config.title : '新增' + config.title;
    
    form.innerHTML = '';
    config.fields.forEach(field => {
        if (field.readonly && !item) return;
        
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
            input.innerHTML = '<option value="">请选择</option>';
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

// 关闭模态框
function closeModal() {
    document.getElementById('edit-modal').classList.remove('show');
    editingItem = null;
}

// 保存数据
async function saveData() {
    const config = ENTITY_CONFIG[currentEntity];
    const form = document.getElementById('edit-form');
    const formData = new FormData(form);
    
    const data = {};
    config.fields.forEach(field => {
        const input = document.getElementById(`edit-${field.name}`);
        if (input) {
            if (field.type === 'checkbox') {
                data[field.name] = input.checked;
            } else if (field.type === 'number') {
                data[field.name] = input.value ? parseFloat(input.value) : null;
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
        
        const response = await fetch(url, {
            method: method,
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(data)
        });
        
        const result = await response.json();
        
        if (response.ok) {
            showToast(result.message || '保存成功', 'success');
            closeModal();
            fetchData();
        } else {
            showToast(result.error || '保存失败', 'error');
        }
    } catch (error) {
        showToast('保存失败: ' + error.message, 'error');
    }
}

// 编辑项
function editItem(id) {
    const item = currentData.find(i => i.id === id);
    if (item) {
        openModal(item);
    }
}

// 删除项
async function deleteItem(id) {
    if (!confirm('确定要删除这条记录吗？')) return;
    
    try {
        const config = ENTITY_CONFIG[currentEntity];
        const response = await fetch(`${API_BASE_URL}/${config.endpoint}/${id}`, {
            method: 'DELETE'
        });
        
        const result = await response.json();
        
        if (response.ok) {
            showToast(result.message || '删除成功', 'success');
            fetchData();
        } else {
            showToast(result.error || '删除失败', 'error');
        }
    } catch (error) {
        showToast('删除失败: ' + error.message, 'error');
    }
}

// 显示提示
function showToast(message, type = 'success') {
    const toast = document.getElementById('toast');
    const toastMessage = document.getElementById('toast-message');
    
    toastMessage.textContent = message;
    toast.className = `toast ${type} show`;
    
    setTimeout(() => {
        toast.classList.remove('show');
    }, 3000);
}
