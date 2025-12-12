document.addEventListener('DOMContentLoaded', () => {
    fetch('/api/logout', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        }
    })
    .then(response => {
        if (response.ok) {
            
            window.location.href = '/login';
        } else {
            console.error('Logout failed');
            window.location.href = '/login';
        }
    })
    .catch(error => {
        console.error('Error:', error);
        window.location.href = '/login';
    });
});
