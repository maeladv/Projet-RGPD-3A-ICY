document.getElementById('loginForm').addEventListener('submit', async function(e) {
    e.preventDefault();
    const username = document.getElementById('username').value;
    const password = document.getElementById('password').value;
    const errorDiv = document.getElementById('loginError');
    errorDiv.style.display = 'none';

    try {
        const response = await fetch('/api/login', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ username, password }),
            credentials: 'include'
        });
        if (!response.ok) {
            const err = await response.text();
            errorDiv.textContent = err;
            errorDiv.style.display = 'block';
            return;
        }
        
        const data = await response.json();
        if (data.User.Role === 'admin') {
            window.location.href = '/admin';
        } else {
            window.location.href = '/rh';
        }
    } catch (err) {
        errorDiv.textContent = 'Erreur de connexion';
        errorDiv.style.display = 'block';
    }
});
