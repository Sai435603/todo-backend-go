const API_BASE = '/api/v1';

// =============================================
// AUTH STORE
// =============================================
const authStore = {
    getToken() { return localStorage.getItem('auth_token'); },
    setToken(t) { localStorage.setItem('auth_token', t); },
    clearToken() { localStorage.removeItem('auth_token'); },

    getUser() { const u = localStorage.getItem('auth_user'); return u ? JSON.parse(u) : null; },
    setUser(u) { localStorage.setItem('auth_user', JSON.stringify(u)); },
    clearUser() { localStorage.removeItem('auth_user'); },

    isAuthenticated() { return !!this.getToken(); },

    clear() { this.clearToken(); this.clearUser(); },
};

// =============================================
// SCREEN MANAGEMENT
// =============================================
function showScreen(screen) {
    var login = document.getElementById('login-screen');
    var app = document.getElementById('app-main');
    if (screen === 'login') {
        login.classList.remove('hidden');
        app.classList.add('hidden');
    } else {
        login.classList.add('hidden');
        app.classList.remove('hidden');
    }
}

function populateUserProfile() {
    var user = authStore.getUser();
    if (!user) return;
    var avatar = document.getElementById('user-avatar');
    var name = document.getElementById('user-name');
    if (user.avatar_url) {
        avatar.src = user.avatar_url;
        avatar.style.display = 'block';
    } else {
        avatar.style.display = 'none';
    }
    name.textContent = user.name || user.email;
}

// =============================================
// OAUTH LOGIN (REDIRECT FLOW)
// =============================================
function startLogin() {
    // Full-page redirect to Google OAuth — no popup needed.
    // After auth, server redirects back to /#token=xxx&user=json
    window.location.href = API_BASE + '/auth/google';
}

// Parse auth credentials from URL fragment after OAuth redirect.
// The backend callback redirects to /#token=xxx&user={json}
function handleOAuthRedirect() {
    var hash = window.location.hash;
    if (!hash || hash.indexOf('token=') === -1) return false;

    // Parse fragment params (everything after #)
    var fragment = hash.substring(1); // remove #
    var params = {};
    var pairs = fragment.split('&');
    for (var i = 0; i < pairs.length; i++) {
        var kv = pairs[i].split('=');
        if (kv.length === 2) {
            params[kv[0]] = decodeURIComponent(kv[1]);
        }
    }

    if (!params.token) return false;

    // Store token
    authStore.setToken(params.token);

    // Store user info if present
    if (params.user) {
        try {
            var user = JSON.parse(params.user);
            authStore.setUser(user);
        } catch (e) {
            // User data malformed, token is still valid
        }
    }

    // Clean URL — remove fragment so token isn't visible
    if (window.history && window.history.replaceState) {
        window.history.replaceState(null, '', window.location.pathname);
    } else {
        window.location.hash = '';
    }

    return true;
}

function logout() {
    authStore.clear();
    showScreen('login');
}

// =============================================
// API CLIENT
// =============================================
var api = {
    request: function (method, path, body) {
        var opts = {
            method: method,
            headers: { 'Content-Type': 'application/json' },
        };

        var token = authStore.getToken();
        if (token) {
            opts.headers['Authorization'] = 'Bearer ' + token;
        }

        if (body) opts.body = JSON.stringify(body);

        return fetch(API_BASE + path, opts).then(function (res) {
            if (res.status === 401) {
                authStore.clear();
                showScreen('login');
                return Promise.reject(new Error('Session expired, please login again'));
            }
            return res.json();
        }).then(function (data) {
            if (!data.success) {
                return Promise.reject(new Error(data.error && data.error.message ? data.error.message : 'Something went wrong'));
            }
            return data.data;
        });
    },

    getTodos: function () { return this.request('GET', '/todos'); },
    getCompleted: function () { return this.request('GET', '/todos/completed'); },
    getPending: function () { return this.request('GET', '/todos/pending'); },
    getTodo: function (id) { return this.request('GET', '/todos/' + id); },
    search: function (q) { return this.request('GET', '/todos/search?q=' + encodeURIComponent(q)); },

    createTodo: function (title, description) {
        return this.request('POST', '/todos', { title: title, description: description });
    },

    updateTodo: function (id, title, description, completed) {
        return this.request('PUT', '/todos/' + id, { title: title, description: description, completed: completed });
    },

    deleteTodo: function (id) { return this.request('DELETE', '/todos/' + id); },
    markCompleted: function (id) { return this.request('PATCH', '/todos/' + id + '/complete'); },
    markPending: function (id) { return this.request('PATCH', '/todos/' + id + '/pending'); },
};

// =============================================
// STATE
// =============================================
var allTodos = [];
var currentFilter = 'all';
var searchDebounceTimer = null;

// =============================================
// DOM HELPERS
// =============================================
function $(sel) { return document.querySelector(sel); }
function $$(sel) { return document.querySelectorAll(sel); }

// =============================================
// RENDERING
// =============================================
function renderTodos(todos) {
    var loadingState = $('#loading-state');
    var emptyState = $('#empty-state');
    var todoList = $('#todo-list');

    loadingState.classList.add('hidden');

    if (!todos || todos.length === 0) {
        todoList.innerHTML = '';
        emptyState.classList.remove('hidden');
        return;
    }

    emptyState.classList.add('hidden');
    var html = '';
    for (var i = 0; i < todos.length; i++) {
        html += createTodoHTML(todos[i], i);
    }
    todoList.innerHTML = html;
}

function createTodoHTML(todo, index) {
    var isCompleted = (todo.completed && todo.completed.Bool) || todo.completed === true;
    var desc = (todo.description && todo.description.String) || todo.description || '';
    var date = (todo.created_at && todo.created_at.Time) || todo.created_at || '';
    var formattedDate = date ? formatDate(date) : '';
    var checkbox = isCompleted ? '[x]' : '[ ]';

    return '<div class="todo-item ' + (isCompleted ? 'completed' : '') + '" data-id="' + todo.id + '" style="animation-delay: ' + (index * 0.04) + 's">'
        + '<div class="todo-checkbox ' + (isCompleted ? 'checked' : '') + '" onclick="toggleTodo(' + todo.id + ', ' + isCompleted + ')">'
        + '<span class="checkbox-text">' + checkbox + '</span>'
        + '</div>'
        + '<div class="todo-content">'
        + '<div class="todo-title"><span class="line-prefix">&gt;</span> ' + escapeHtml(todo.title) + '</div>'
        + (desc ? '<div class="todo-description"><span class="comment-prefix">#</span> ' + escapeHtml(desc) + '</div>' : '')
        + '<div class="todo-meta">'
        + (formattedDate ? '<span class="todo-date">' + formattedDate + '</span>' : '')
        + '<span class="todo-badge ' + (isCompleted ? 'done' : 'pending') + '">' + (isCompleted ? 'DONE' : 'PENDING') + '</span>'
        + '</div></div>'
        + '<div class="todo-actions">'
        + '<button class="action-btn" onclick="openEditModal(' + todo.id + ')" title="Edit">[edit]</button>'
        + '<button class="action-btn delete" onclick="deleteTodo(' + todo.id + ')" title="Delete">[rm]</button>'
        + '</div></div>';
}

function updateStats() {
    var statTotal = $('#stat-total .stat-number');
    var statPending = $('#stat-pending .stat-number');
    var statCompleted = $('#stat-completed .stat-number');

    var total = allTodos.length;
    var completed = 0;
    for (var i = 0; i < allTodos.length; i++) {
        if ((allTodos[i].completed && allTodos[i].completed.Bool) || allTodos[i].completed === true) {
            completed++;
        }
    }
    var pending = total - completed;

    animateNumber(statTotal, total);
    animateNumber(statPending, pending);
    animateNumber(statCompleted, completed);
}

function animateNumber(el, target) {
    var current = parseInt(el.textContent) || 0;
    if (current === target) return;
    var diff = target - current;
    var steps = Math.min(Math.abs(diff), 15);
    var duration = 300;
    var stepTime = duration / steps;
    var step = 0;
    var interval = setInterval(function () {
        step++;
        var progress = step / steps;
        var eased = 1 - Math.pow(1 - progress, 3);
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
function loadTodos() {
    var loadingState = $('#loading-state');
    var emptyState = $('#empty-state');
    var todoList = $('#todo-list');

    loadingState.classList.remove('hidden');
    emptyState.classList.add('hidden');
    todoList.innerHTML = '';

    api.getTodos().then(function (todos) {
        allTodos = todos || [];
        updateStats();
        applyFilter();
    }).catch(function (err) {
        loadingState.classList.add('hidden');
        showToast(err.message || 'Failed to load todos', 'error');
        console.error(err);
    });
}

function applyFilter() {
    var filtered = allTodos;
    if (currentFilter === 'completed') {
        filtered = allTodos.filter(function (t) {
            return (t.completed && t.completed.Bool) || t.completed === true;
        });
    } else if (currentFilter === 'pending') {
        filtered = allTodos.filter(function (t) {
            return !((t.completed && t.completed.Bool) || t.completed === true);
        });
    }
    renderTodos(filtered);
}

// =============================================
// ACTIONS
// =============================================
function toggleTodo(id, isCurrentlyCompleted) {
    var promise = isCurrentlyCompleted ? api.markPending(id) : api.markCompleted(id);
    var msg = isCurrentlyCompleted ? 'Marked as pending' : 'Marked as completed';
    promise.then(function () {
        showToast(msg, 'success');
        loadTodos();
    }).catch(function (err) {
        showToast(err.message, 'error');
    });
}

function deleteTodo(id) {
    var item = document.querySelector('.todo-item[data-id="' + id + '"]');
    if (item) item.classList.add('removing');

    setTimeout(function () {
        api.deleteTodo(id).then(function () {
            showToast('Todo deleted', 'success');
            loadTodos();
        }).catch(function (err) {
            showToast(err.message, 'error');
            if (item) item.classList.remove('removing');
        });
    }, 350);
}

function openEditModal(id) {
    var todo = null;
    for (var i = 0; i < allTodos.length; i++) {
        if (allTodos[i].id === id) { todo = allTodos[i]; break; }
    }
    if (!todo) return;
    $('#edit-id').value = id;
    $('#edit-title').value = todo.title;
    $('#edit-description').value = (todo.description && todo.description.String) || todo.description || '';
    $('#edit-completed').value = (todo.completed && todo.completed.Bool) || todo.completed === true;
    $('#edit-modal').classList.remove('hidden');
    $('#edit-title').focus();
}

function closeEditModal() {
    $('#edit-modal').classList.add('hidden');
}

// =============================================
// EVENT LISTENERS
// =============================================
function bindEvents() {
    var addForm = $('#add-todo-form');
    var editForm = $('#edit-form');
    var modalCloseBtn = $('#modal-close');
    var modalCancelBtn = $('#modal-cancel');
    var editModal = $('#edit-modal');
    var searchInput = $('#search-input');
    var descInput = $('#todo-description');

    addForm.addEventListener('submit', function (e) {
        e.preventDefault();
        var titleInput = $('#todo-title');
        var title = titleInput.value.trim();
        var desc = descInput.value.trim();
        if (!title) return;

        var addBtn = addForm.querySelector('.btn-add');
        addBtn.disabled = true;
        addBtn.style.opacity = '0.5';

        api.createTodo(title, desc).then(function () {
            titleInput.value = '';
            descInput.value = '';
            showToast('Todo created!', 'success');
            addBtn.disabled = false;
            addBtn.style.opacity = '1';
            loadTodos();
        }).catch(function (err) {
            showToast(err.message, 'error');
            addBtn.disabled = false;
            addBtn.style.opacity = '1';
        });
    });

    editForm.addEventListener('submit', function (e) {
        e.preventDefault();
        var id = parseInt($('#edit-id').value);
        var title = $('#edit-title').value.trim();
        var desc = $('#edit-description').value.trim();
        var completed = $('#edit-completed').value === 'true';
        if (!title) return;

        api.updateTodo(id, title, desc, completed).then(function () {
            closeEditModal();
            showToast('Todo updated!', 'success');
            loadTodos();
        }).catch(function (err) {
            showToast(err.message, 'error');
        });
    });

    modalCloseBtn.addEventListener('click', closeEditModal);
    modalCancelBtn.addEventListener('click', closeEditModal);
    editModal.addEventListener('click', function (e) {
        if (e.target === editModal) closeEditModal();
    });
    document.addEventListener('keydown', function (e) {
        if (e.key === 'Escape' && !editModal.classList.contains('hidden')) {
            closeEditModal();
        }
    });

    var filterBtns = $$('.filter-btn');
    for (var i = 0; i < filterBtns.length; i++) {
        filterBtns[i].addEventListener('click', function () {
            for (var j = 0; j < filterBtns.length; j++) filterBtns[j].classList.remove('active');
            this.classList.add('active');
            currentFilter = this.dataset.filter;
            applyFilter();
        });
    }

    searchInput.addEventListener('input', function () {
        clearTimeout(searchDebounceTimer);
        searchDebounceTimer = setTimeout(function () {
            var query = searchInput.value.trim();
            if (!query) { loadTodos(); return; }
            api.search(query).then(function (results) {
                allTodos = results || [];
                updateStats();
                applyFilter();
            }).catch(function (err) {
                showToast(err.message, 'error');
            });
        }, 350);
    });

    descInput.addEventListener('input', function () {
        descInput.style.height = 'auto';
        descInput.style.height = descInput.scrollHeight + 'px';
    });
}

// =============================================
// TOAST
// =============================================
function showToast(message, type) {
    type = type || 'success';
    var toastContainer = $('#toast-container');
    var toast = document.createElement('div');
    toast.className = 'toast ' + type;
    toast.innerHTML = '<span>' + escapeHtml(message) + '</span>';
    toastContainer.appendChild(toast);
    setTimeout(function () {
        toast.classList.add('removing');
        setTimeout(function () { toast.remove(); }, 300);
    }, 3000);
}

// =============================================
// UTILITIES
// =============================================
function escapeHtml(str) {
    var div = document.createElement('div');
    div.textContent = str;
    return div.innerHTML;
}

function formatDate(dateStr) {
    try {
        var d = new Date(dateStr);
        var now = new Date();
        var diffMs = now - d;
        var diffMins = Math.floor(diffMs / 60000);
        var diffHours = Math.floor(diffMs / 3600000);
        var diffDays = Math.floor(diffMs / 86400000);
        if (diffMins < 1) return 'Just now';
        if (diffMins < 60) return diffMins + 'm ago';
        if (diffHours < 24) return diffHours + 'h ago';
        if (diffDays < 7) return diffDays + 'd ago';
        return d.toLocaleDateString('en-US', { month: 'short', day: 'numeric' });
    } catch (e) { return ''; }
}

// =============================================
// INIT
// =============================================
document.addEventListener('DOMContentLoaded', function () {
    bindEvents();

    // Check if we just returned from OAuth redirect
    var justLoggedIn = handleOAuthRedirect();

    // Check for auth errors in query string
    var urlParams = new URLSearchParams(window.location.search);
    var authError = urlParams.get('auth_error');
    if (authError) {
        showToast('Login failed: ' + authError.replace(/_/g, ' '), 'error');
        // Clean the error from URL
        if (window.history && window.history.replaceState) {
            window.history.replaceState(null, '', window.location.pathname);
        }
    }

    if (justLoggedIn || authStore.isAuthenticated()) {
        showScreen('app');
        populateUserProfile();
        loadTodos();
    } else {
        showScreen('login');
    }
});
