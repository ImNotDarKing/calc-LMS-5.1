const apiUrl = '/api/v1';
let token = localStorage.getItem('token');

function showSection() {
    if (token) {
        document.getElementById('auth-section').classList.add('hidden');
        document.getElementById('app-section').classList.remove('hidden');
        loadExpressions();
    } else {
        document.getElementById('auth-section').classList.remove('hidden');
        document.getElementById('app-section').classList.add('hidden');
    }
}

async function register() {
    const login = document.getElementById('reg-login').value;
    const password = document.getElementById('reg-password').value;
    const res = await fetch(`${apiUrl}/register`, {
        method: 'POST', headers: {'Content-Type':'application/json'},
        body: JSON.stringify({login, password})
    });
    if (res.ok) {
        const data = await res.json();
        token = data.token;
        localStorage.setItem('token', token);
        showSection();
    } else alert('Register failed');
}

async function loginUser() {
    const login = document.getElementById('login-login').value;
    const password = document.getElementById('login-password').value;
    const res = await fetch(`${apiUrl}/login`, {
        method: 'POST', headers: {'Content-Type':'application/json'},
        body: JSON.stringify({login, password})
    });
    if (res.ok) {
        const data = await res.json();
        token = data.token;
        localStorage.setItem('token', token);
        showSection();
    } else alert('Login failed');
}

async function logout() {
    token = null;
    localStorage.removeItem('token');
    showSection();
}

async function submitExpression() {
    const expression = document.getElementById('expr-input').value;
    const res = await fetch(`${apiUrl}/calculate`, {
        method: 'POST', headers: {'Content-Type':'application/json', 'Authorization': 'Bearer ' + token},
        body: JSON.stringify({expression})
    });
    if (res.ok) {
        document.getElementById('expr-input').value = '';
        loadExpressions();
    } else alert('Submit failed');
}

async function loadExpressions() {
    const res = await fetch(`${apiUrl}/expressions`, {
        headers: {'Authorization': 'Bearer ' + token}
    });
    const listEl = document.getElementById('expr-list');
    listEl.innerHTML = '';
    if (res.ok) {
        const data = await res.json();
        data.expressions.forEach(e => {
            const li = document.createElement('li');
            li.textContent = `#${e.id}: ${e.expression} = ${e.result}`;
            listEl.appendChild(li);
        });
    }
}

// Event listeners
document.getElementById('btn-register').onclick = register;
document.getElementById('btn-login').onclick = loginUser;
document.getElementById('btn-logout').onclick = logout;
document.getElementById('btn-submit').onclick = submitExpression;

showSection();