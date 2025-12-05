// API Base URL
const API_URL = window.location.hostname === 'localhost' || window.location.hostname === '127.0.0.1'
    ? 'http://localhost:8080/api'
    : `${window.location.origin}/api`;

// Cambio entre formularios
document.getElementById('show-register')?.addEventListener('click', (e) => {
    e.preventDefault();
    document.getElementById('login-form').classList.remove('active');
    document.getElementById('register-form').classList.add('active');
    hideError();
});

document.getElementById('show-login')?.addEventListener('click', (e) => {
    e.preventDefault();
    document.getElementById('register-form').classList.remove('active');
    document.getElementById('login-form').classList.add('active');
    hideError();
});

// Manejo de login
document.getElementById('loginForm')?.addEventListener('submit', async (e) => {
    e.preventDefault();
    hideError();

    const email = document.getElementById('login-email').value;
    const password = document.getElementById('login-password').value;

    try {
        const response = await fetch(`${API_URL}/auth/login`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ email, password })
        });

        if (response.ok) {
            const data = await response.json();
            
            // Guardar token y usuario en localStorage
            localStorage.setItem('auth_token', data.token);
            localStorage.setItem('user', JSON.stringify(data.user));
            
            // Redirigir al dashboard
            window.location.href = '/';
        } else {
            const error = await response.text();
            showError(error || 'Error al iniciar sesión');
        }
    } catch (error) {
        console.error('Error en login:', error);
        showError('Error de conexión. Por favor, intenta de nuevo.');
    }
});

// Manejo de registro
document.getElementById('registerForm')?.addEventListener('submit', async (e) => {
    e.preventDefault();
    hideError();

    const name = document.getElementById('register-name').value;
    const email = document.getElementById('register-email').value;
    const password = document.getElementById('register-password').value;

    if (password.length < 8) {
        showError('La contraseña debe tener al menos 8 caracteres');
        return;
    }

    try {
        const response = await fetch(`${API_URL}/auth/register`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ name, email, password })
        });

        if (response.ok) {
            const data = await response.json();
            
            // Guardar token y usuario en localStorage
            localStorage.setItem('auth_token', data.token);
            localStorage.setItem('user', JSON.stringify(data.user));
            
            // Redirigir al dashboard
            window.location.href = '/';
        } else {
            const error = await response.text();
            showError(error || 'Error al registrarse');
        }
    } catch (error) {
        console.error('Error en registro:', error);
        showError('Error de conexión. Por favor, intenta de nuevo.');
    }
});

// Funciones auxiliares
function showError(message) {
    const errorDiv = document.getElementById('auth-error');
    errorDiv.textContent = message;
    errorDiv.style.display = 'block';
}

function hideError() {
    const errorDiv = document.getElementById('auth-error');
    errorDiv.style.display = 'none';
}

// Verificar si ya hay sesión activa
function checkExistingSession() {
    const token = localStorage.getItem('auth_token');
    if (token) {
        // Redirigir al dashboard si ya hay sesión
        window.location.href = '/';
    }
}

checkExistingSession();
