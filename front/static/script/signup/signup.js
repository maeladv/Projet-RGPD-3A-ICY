document.getElementById('signupForm').addEventListener('submit', async function(e) {
    e.preventDefault();

    const username = document.getElementById('username').value.trim();
    const password = document.getElementById('password').value;
    const role = document.getElementById('role').value;

    const response = await fetch('/api/user/add', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ username, password, role })
    });

    const messageDiv = document.getElementById('signupMessage');
    if (response.ok) {
        messageDiv.style.color = 'green';
        messageDiv.textContent = "Utilisateur créé avec succès !";
        document.getElementById('signupForm').reset();
        window.location.href = '/login';
    } else {
        const data = await response.json();
        messageDiv.style.color = '#d9534f';
        messageDiv.textContent = data.error || "Erreur lors de la création.";
    }
});