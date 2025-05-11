const API_BASE = '/api/v1';
let token = localStorage.getItem('jwt') || '';

const loginForm    = document.getElementById('login-form');
const registerForm = document.getElementById('register-form');
const appDiv       = document.getElementById('app');
const logoutBtn    = document.getElementById('btn-logout');
const btnShowLogin = document.getElementById('btn-show-login');
const btnShowReg   = document.getElementById('btn-show-register');

btnShowLogin.addEventListener('click', () => {
  loginForm.classList.remove('hidden');
  registerForm.classList.add('hidden');
});
btnShowReg.addEventListener('click', () => {
  registerForm.classList.remove('hidden');
  loginForm.classList.add('hidden');
});
logoutBtn.addEventListener('click', () => {
  localStorage.removeItem('jwt');
  token = '';
  location.reload();
});

function authFetch(url, options = {}) {
  options.headers = options.headers || {};
  options.headers['Content-Type'] = 'application/json';
  if (token) {
    options.headers['Authorization'] = 'Bearer ' + token;
  }
  return fetch(url, options);
}

function updateUI() {
  if (token) {
    loginForm.classList.add('hidden');
    registerForm.classList.add('hidden');
    appDiv.classList.remove('hidden');
    logoutBtn.classList.remove('hidden');
    loadResults();
  } else {
    loginForm.classList.remove('hidden');
    registerForm.classList.add('hidden');
    appDiv.classList.add('hidden');
    logoutBtn.classList.add('hidden');
  }
}

loginForm.addEventListener('submit', async e => {
  e.preventDefault();
  document.getElementById('login-error').textContent = '';
  const login    = document.getElementById('login-login').value.trim();
  const password = document.getElementById('login-password').value;
  try {
    const resp = await fetch(`${API_BASE}/login`, {
      method: 'POST',
      headers: {'Content-Type':'application/json'},
      body: JSON.stringify({login, password})
    });
    const data = await resp.json();
    if (resp.ok && data.token) {
      token = data.token;
      localStorage.setItem('jwt', token);
      updateUI();
    } else {
      document.getElementById('login-error').textContent = data.error || 'Ошибка входа';
    }
  } catch (err) {
    document.getElementById('login-error').textContent = 'Сетевая ошибка';
  }
});

registerForm.addEventListener('submit', async e => {
  e.preventDefault();
  document.getElementById('register-error').textContent = '';
  const login    = document.getElementById('reg-login').value.trim();
  const password = document.getElementById('reg-password').value;
  try {
    const resp = await fetch(`${API_BASE}/register`, {
      method: 'POST',
      headers: {'Content-Type':'application/json'},
      body: JSON.stringify({login, password})
    });
    const data = await resp.json();
    if (resp.ok && data.token) {
      token = data.token;
      localStorage.setItem('jwt', token);
      updateUI();
    } else {
      document.getElementById('register-error').textContent = data.error || 'Ошибка регистрации';
    }
  } catch (err) {
    document.getElementById('register-error').textContent = 'Сетевая ошибка';
  }
});

const calcForm = document.getElementById('calc-form');
calcForm.addEventListener('submit', async e => {
  e.preventDefault();
  document.getElementById('calc-error').textContent = '';
  const raw = document.getElementById('expression').value.trim();
  if (!raw) return;
  try {
    const resp = await authFetch(`${API_BASE}/calculate`, {
      method: 'POST',
      body: JSON.stringify({expression: raw})
    });
    if (resp.status === 201) {
      document.getElementById('expression').value = '';
      setTimeout(loadResults, 100);
    } else {
      const err = await resp.json();
      document.getElementById('calc-error').textContent = err.error || 'Ошибка';
    }
  } catch (err) {
    document.getElementById('calc-error').textContent = 'Сетевая ошибка';
  }
});

async function loadResults() {
  try {
    const resp = await authFetch(`${API_BASE}/expressions`);
    if (!resp.ok) return;
    const data = await resp.json();
    const tbody = document.querySelector('#results-table tbody');
    tbody.innerHTML = '';
    data.expressions.forEach(expr => {
      const tr = document.createElement('tr');
      tr.innerHTML = `
        <td>${expr.id}</td>
        <td>${expr.expression}</td>
        <td>${expr.result}</td>
      `;
      tbody.appendChild(tr);
    });
  } catch (err) {
    console.error(err);
  }
}

updateUI();

setInterval(loadResults, 100);