const API_BASE = '/api/v1';

const api = {
    async request(method, path, body = null) {
        const opts = {
            method,
            headers: { 'Content-Type': 'application/json' },
        };
        if (body) opts.body = JSON.stringify(body);

        const res = await fetch(`${API_BASE}${path}`, opts);
        const data = await res.json();

        if (!data.success) {
            throw new Error(data.error?.message || 'Something went wrong');
        }
        return data.data;
    },

    getTodos()          { return this.request('GET', '/todos'); },
    getCompleted()      { return this.request('GET', '/todos/completed'); },
    getPending()        { return this.request('GET', '/todos/pending'); },
    getTodo(id)         { return this.request('GET', `/todos/${id}`); },
    search(q)           { return this.request('GET', `/todos/search?q=${encodeURIComponent(q)}`); },

    createTodo(title, description) {
        return this.request('POST', '/todos', { title, description });
    },

    updateTodo(id, title, description, completed) {
        return this.request('PUT', `/todos/${id}`, { title, description, completed });
    },

    deleteTodo(id)       { return this.request('DELETE', `/todos/${id}`); },
    markCompleted(id)    { return this.request('PATCH', `/todos/${id}/complete`); },
    markPending(id)      { return this.request('PATCH', `/todos/${id}/pending`); },
};


// =============================================
// STATE
// =============================================
let allTodos = [];
let currentFilter = 'all';
let searchDebounceTimer = null;


// =============================================
// DOM REFERENCES
// =============================================
const $ = (sel) => document.querySelector(sel);
const $$ = (sel) => document.querySelectorAll(sel);

const todoList       = $('#todo-list');
const emptyState     = $('#empty-state');
const loadingState   = $('#loading-state');
const addForm        = $('#add-todo-form');
const titleInput     = $('#todo-title');
const descInput      = $('#todo-description');
const searchInput    = $('#search-input');
const editModal      = $('#edit-modal');
const editForm       = $('#edit-form');
const editId         = $('#edit-id');
const editTitle      = $('#edit-title');
const editDesc       = $('#edit-description');
const editCompleted  = $('#edit-completed');
const modalCloseBtn  = $('#modal-close');
const modalCancelBtn = $('#modal-cancel');
const statTotal      = $('#stat-total .stat-number');
const statPending    = $('#stat-pending .stat-number');
const statCompleted  = $('#stat-completed .stat-number');
const toastContainer = $('#toast-container');


// =============================================
// TERMINAL EFFECTS
// =============================================
function initTerminalEffects() {
    // Nothing heavy — the CSS handles scanlines and cursor blink
}


// =============================================
// RENDERING
// =============================================
function renderTodos(todos) {
    loadingState.classList.add('hidden');

    if (!todos || todos.length === 0) {
        todoList.innerHTML = '';
        emptyState.classList.remove('hidden');
        return;
    }

    emptyState.classList.add('hidden');
    todoList.innerHTML = todos.map((todo, i) => createTodoHTML(todo, i)).join('');
}

function createTodoHTML(todo, index) {
    const isCompleted = todo.completed?.Bool || todo.completed === true;
    const desc = todo.description?.String || todo.description || '';
    const date = todo.created_at?.Time || todo.created_at || '';
    const formattedDate = date ? formatDate(date) : '';
    const checkbox = isCompleted ? '[x]' : '[ ]';

    return `
        <div class="todo-item ${isCompleted ? 'completed' : ''}" data-id="${todo.id}" style="animation-delay: ${index * 0.04}s">
            <div class="todo-checkbox ${isCompleted ? 'checked' : ''}" onclick="toggleTodo(${todo.id}, ${isCompleted})">
                <span class="checkbox-text">${checkbox}</span>
            </div>
            <div class="todo-content">
                <div class="todo-title"><span class="line-prefix">&gt;</span> ${escapeHtml(todo.title)}</div>
                ${desc ? `<div class="todo-description"><span class="comment-prefix">#</span> ${escapeHtml(desc)}</div>` : ''}
                <div class="todo-meta">
                    ${formattedDate ? `<span class="todo-date">${formattedDate}</span>` : ''}
                    <span class="todo-badge ${isCompleted ? 'done' : 'pending'}">${isCompleted ? 'DONE' : 'PENDING'}</span>
                </div>
            </div>
            <div class="todo-actions">
                <button class="action-btn" onclick="openEditModal(${todo.id})" title="Edit">[edit]</button>
                <button class="action-btn delete" onclick="deleteTodo(${todo.id})" title="Delete">[rm]</button>
            </div>
        </div>
    `;
}

function updateStats() {
    const total = allTodos.length;
    const completed = allTodos.filter(t => t.completed?.Bool || t.completed === true).length;
    const pending = total - completed;

    animateNumber(statTotal, total);
    animateNumber(statPending, pending);
    animateNumber(statCompleted, completed);
}

function animateNumber(el, target) {
    const current = parseInt(el.textContent) || 0;
    if (current === target) return;

    const diff = target - current;
    const steps = Math.min(Math.abs(diff), 15);
    const duration = 300;
    const stepTime = duration / steps;
    let step = 0;

    const interval = setInterval(() => {
        step++;
        const progress = step / steps;
        const eased = 1 - Math.pow(1 - progress, 3);
        el.textContent = Math.round(current + diff * eased);
        if (step >= steps) {
            el.textContent = target;
            clearInterval(interval);
        }
    }, stepTime);
}


// =============================================
// DATA FETCHING
// =============================================
async function loadTodos() {
    try {
        loadingState.classList.remove('hidden');
        emptyState.classList.add('hidden');
        todoList.innerHTML = '';

        allTodos = await api.getTodos() || [];
        updateStats();
        applyFilter();
    } catch (err) {
        loadingState.classList.add('hidden');
        showToast('Failed to load todos', 'error');
        console.error(err);
    }
}

function applyFilter() {
    let filtered = allTodos;

    if (currentFilter === 'completed') {
        filtered = allTodos.filter(t => t.completed?.Bool || t.completed === true);
    } else if (currentFilter === 'pending') {
        filtered = allTodos.filter(t => !(t.completed?.Bool || t.completed === true));
    }

    renderTodos(filtered);
}


// =============================================
// ACTIONS
// =============================================
async function toggleTodo(id, isCurrentlyCompleted) {
    try {
        if (isCurrentlyCompleted) {
            await api.markPending(id);
            showToast('Marked as pending', 'success');
        } else {
            await api.markCompleted(id);
            showToast('Marked as completed', 'success');
        }

        // Animate the checkbox
        const item = document.querySelector(`.todo-item[data-id="${id}"]`);
        if (item) {
            const checkbox = item.querySelector('.todo-checkbox');
            checkbox.style.animation = 'checkBounce 0.3s ease';
            setTimeout(() => checkbox.style.animation = '', 300);
        }

        await loadTodos();
    } catch (err) {
        showToast(err.message, 'error');
    }
}

async function deleteTodo(id) {
    const item = document.querySelector(`.todo-item[data-id="${id}"]`);
    if (item) {
        item.classList.add('removing');
        await new Promise(r => setTimeout(r, 350));
    }

    try {
        await api.deleteTodo(id);
        showToast('Todo deleted', 'success');
        await loadTodos();
    } catch (err) {
        showToast(err.message, 'error');
        if (item) item.classList.remove('removing');
    }
}

async function openEditModal(id) {
    const todo = allTodos.find(t => t.id === id);
    if (!todo) return;

    editId.value = id;
    editTitle.value = todo.title;
    editDesc.value = todo.description?.String || todo.description || '';
    editCompleted.value = todo.completed?.Bool || todo.completed === true;
    editModal.classList.remove('hidden');
    editTitle.focus();
}

function closeEditModal() {
    editModal.classList.add('hidden');
}


// =============================================
// EVENT LISTENERS
// =============================================

// Add todo
addForm.addEventListener('submit', async (e) => {
    e.preventDefault();
    const title = titleInput.value.trim();
    const desc = descInput.value.trim();

    if (!title) return;

    try {
        const addBtn = addForm.querySelector('.btn-add');
        addBtn.disabled = true;
        addBtn.style.opacity = '0.5';

        await api.createTodo(title, desc);
        titleInput.value = '';
        descInput.value = '';
        showToast('Todo created!', 'success');
        await loadTodos();

        addBtn.disabled = false;
        addBtn.style.opacity = '1';
    } catch (err) {
        showToast(err.message, 'error');
        const addBtn = addForm.querySelector('.btn-add');
        addBtn.disabled = false;
        addBtn.style.opacity = '1';
    }
});

// Edit form submit
editForm.addEventListener('submit', async (e) => {
    e.preventDefault();
    const id = parseInt(editId.value);
    const title = editTitle.value.trim();
    const desc = editDesc.value.trim();
    const completed = editCompleted.value === 'true';

    if (!title) return;

    try {
        await api.updateTodo(id, title, desc, completed);
        closeEditModal();
        showToast('Todo updated!', 'success');
        await loadTodos();
    } catch (err) {
        showToast(err.message, 'error');
    }
});

// Close modal
modalCloseBtn.addEventListener('click', closeEditModal);
modalCancelBtn.addEventListener('click', closeEditModal);
editModal.addEventListener('click', (e) => {
    if (e.target === editModal) closeEditModal();
});
document.addEventListener('keydown', (e) => {
    if (e.key === 'Escape' && !editModal.classList.contains('hidden')) {
        closeEditModal();
    }
});

// Filter tabs
$$('.filter-btn').forEach(btn => {
    btn.addEventListener('click', () => {
        $$('.filter-btn').forEach(b => b.classList.remove('active'));
        btn.classList.add('active');
        currentFilter = btn.dataset.filter;
        applyFilter();
    });
});

// Search with debounce
searchInput.addEventListener('input', () => {
    clearTimeout(searchDebounceTimer);
    searchDebounceTimer = setTimeout(async () => {
        const query = searchInput.value.trim();
        if (!query) {
            await loadTodos();
            return;
        }
        try {
            const results = await api.search(query) || [];
            allTodos = results;
            updateStats();
            applyFilter();
        } catch (err) {
            showToast(err.message, 'error');
        }
    }, 350);
});

// Auto-resize description textarea
descInput.addEventListener('input', () => {
    descInput.style.height = 'auto';
    descInput.style.height = descInput.scrollHeight + 'px';
});


// =============================================
// TOAST NOTIFICATIONS
// =============================================
function showToast(message, type = 'success') {
    const iconSvg = type === 'success'
        ? '<svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3" stroke-linecap="round" stroke-linejoin="round"><polyline points="20 6 9 17 4 12"/></svg>'
        : '<svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3" stroke-linecap="round" stroke-linejoin="round"><line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/></svg>';

    const toast = document.createElement('div');
    toast.className = `toast ${type}`;
    toast.innerHTML = `
        <div class="toast-icon">${iconSvg}</div>
        <span>${escapeHtml(message)}</span>
    `;
    toastContainer.appendChild(toast);

    setTimeout(() => {
        toast.classList.add('removing');
        setTimeout(() => toast.remove(), 300);
    }, 3000);
}


// =============================================
// UTILITIES
// =============================================
function escapeHtml(str) {
    const div = document.createElement('div');
    div.textContent = str;
    return div.innerHTML;
}

function formatDate(dateStr) {
    try {
        const d = new Date(dateStr);
        const now = new Date();
        const diffMs = now - d;
        const diffMins = Math.floor(diffMs / 60000);
        const diffHours = Math.floor(diffMs / 3600000);
        const diffDays = Math.floor(diffMs / 86400000);

        if (diffMins < 1) return 'Just now';
        if (diffMins < 60) return `${diffMins}m ago`;
        if (diffHours < 24) return `${diffHours}h ago`;
        if (diffDays < 7) return `${diffDays}d ago`;

        return d.toLocaleDateString('en-US', {
            month: 'short',
            day: 'numeric',
            year: d.getFullYear() !== now.getFullYear() ? 'numeric' : undefined,
        });
    } catch {
        return '';
    }
}


// =============================================
// INIT
// =============================================
document.addEventListener('DOMContentLoaded', () => {
    initTerminalEffects();
    loadTodos();
});
