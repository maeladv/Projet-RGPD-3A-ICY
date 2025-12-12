const API_URL = '/api';

let allForms = [];
let filteredForms = [];

// Vérification de l'authentification au chargement
document.addEventListener('DOMContentLoaded', () => {
    if (!getCookie('jwt')) {
        window.location.href = '/login';
        return;
    }
    loadForms();
    setupEventListeners();
});

function setupEventListeners() {
    document.getElementById('logoutBtn').addEventListener('click', logout);
    
    // Toggle filters
    document.getElementById('toggleFiltersBtn').addEventListener('click', () => {
        const panel = document.getElementById('advancedFilters');
        panel.style.display = panel.style.display === 'none' ? 'block' : 'none';
    });

    // Global search
    document.getElementById('globalSearch').addEventListener('input', (e) => {
        applyFilters(e.target.value);
    });

    document.getElementById('reloadBtn').addEventListener('click', () => {
        resetFilters();
        loadForms();
    });

    document.getElementById('applyFilters').addEventListener('click', () => applyFilters());
    document.getElementById('resetFilters').addEventListener('click', resetFilters);
    document.getElementById('closeDetail').addEventListener('click', closeDetail);
}

function getCookie(name) {
    const value = `; ${document.cookie}`;
    const parts = value.split(`; ${name}=`);
    if (parts.length === 2) return parts.pop().split(';').shift();
}

async function loadForms() {
    const loading = document.getElementById('loading');
    const error = document.getElementById('error');
    const table = document.getElementById('formsTable');

    try {
        loading.style.display = 'block';
        error.style.display = 'none';
        table.style.display = 'none';

        const response = await fetch(`${API_URL}/forms`, {
            credentials: 'include', // Envoie le cookie JWT
            headers: {
                'Content-Type': 'application/json'
            }
        });

        if (!response.ok) {
            if (response.status === 401) {
                window.location.href = '/login';
                return;
            }
            throw new Error('Erreur lors du chargement des formulaires');
        }

        allForms = await response.json();
        filteredForms = [...allForms];
        displayForms(filteredForms);

        loading.style.display = 'none';
        table.style.display = 'table';
    } catch (err) {
        loading.style.display = 'none';
        error.textContent = err.message;
        error.style.display = 'block';
    }
}

function displayForms(forms) {
    const tbody = document.getElementById('formsBody');
    tbody.innerHTML = '';

    if (forms.length === 0) {
        tbody.innerHTML = '<tr><td colspan="8">Aucun formulaire trouvé</td></tr>';
        return;
    }

    forms.forEach(form => {
        const row = document.createElement('tr');
        row.innerHTML = `
            <td>${form.ID || form.id}</td>
            <td>${form.Nom || ''}</td>
            <td>${form.Prenom || ''}</td>
            <td>${form.Mail || ''}</td>
            <td>${form.Telephone || ''}</td>
            <td>${formatDate(form.DateNaissance)}</td>
            <td>${form.NiveauDiplome || ''}</td>
            <td>
                <button class="viewBtn" data-id="${form.ID || form.id}">Voir détails</button>
            </td>
        `;
        tbody.appendChild(row);
    });

    // Ajoute les événements sur les boutons "Voir détails"
    document.querySelectorAll('.viewBtn').forEach(btn => {
        btn.addEventListener('click', (e) => {
            const id = e.target.dataset.id;
            showFormDetail(id);
        });
    });
}

function formatDate(dateStr) {
    if (!dateStr) return '';
    const date = new Date(dateStr);
    return date.toLocaleDateString('fr-FR');
}

function applyFilters(globalSearchTerm = null) {
    // Si globalSearchTerm est fourni, c'est une recherche rapide
    // Sinon on prend les valeurs des champs
    
    let globalTerm = globalSearchTerm;
    if (globalTerm === null) {
        globalTerm = document.getElementById('globalSearch').value.toLowerCase();
    } else {
        globalTerm = globalTerm.toLowerCase();
    }

    const searchName = document.getElementById('searchName').value.toLowerCase();
    const searchSurname = document.getElementById('searchSurname').value.toLowerCase();
    const searchEmail = document.getElementById('searchEmail').value.toLowerCase();
    const searchPhone = document.getElementById('searchPhone').value.toLowerCase();

    filteredForms = allForms.filter(form => {
        // Filtre global (cherche dans tous les champs pertinents)
        const matchGlobal = !globalTerm || 
            (form.Nom && form.Nom.toLowerCase().includes(globalTerm)) ||
            (form.Prenom && form.Prenom.toLowerCase().includes(globalTerm)) ||
            (form.Mail && form.Mail.toLowerCase().includes(globalTerm)) ||
            (form.Telephone && form.Telephone.toLowerCase().includes(globalTerm));

        // Filtres spécifiques
        const matchName = !searchName || (form.Nom && form.Nom.toLowerCase().includes(searchName));
        const matchSurname = !searchSurname || (form.Prenom && form.Prenom.toLowerCase().includes(searchSurname));
        const matchEmail = !searchEmail || (form.Mail && form.Mail.toLowerCase().includes(searchEmail));
        const matchPhone = !searchPhone || (form.Telephone && form.Telephone.toLowerCase().includes(searchPhone));
        
        return matchGlobal && matchName && matchSurname && matchEmail && matchPhone;
    });

    displayForms(filteredForms);
}

function resetFilters() {
    document.getElementById('globalSearch').value = '';
    document.getElementById('searchName').value = '';
    document.getElementById('searchSurname').value = '';
    document.getElementById('searchEmail').value = '';
    document.getElementById('searchPhone').value = '';
    
    filteredForms = [...allForms];
    displayForms(filteredForms);
}

async function showFormDetail(id) {
    const detailSection = document.getElementById('formDetail');
    const detailContent = document.getElementById('detailContent');

    try {
        const response = await fetch(`${API_URL}/form?id=${id}`, {
            credentials: 'include'
        });

        if (!response.ok) {
            throw new Error('Erreur lors du chargement des détails');
        }

        const form = await response.json();
        
        detailContent.innerHTML = `
            <div class="detail-grid">
                <div><strong>ID:</strong> ${form.ID || form.id}</div>
                <div><strong>Nom:</strong> ${form.Nom || ''}</div>
                <div><strong>Prénom:</strong> ${form.Prenom || ''}</div>
                <div><strong>Date de naissance:</strong> ${formatDate(form.DateNaissance)}</div>
                <div><strong>Ville de naissance:</strong> ${form.VilleNaissance || ''}</div>
                <div><strong>Niveau diplôme:</strong> ${form.NiveauDiplome || ''}</div>
                <div><strong>Email:</strong> ${form.Mail || ''}</div>
                <div><strong>Téléphone:</strong> ${form.Telephone || ''}</div>
                <div><strong>Adresse:</strong> ${form.Adresse || ''}</div>
                <div><strong>Complément:</strong> ${form.Complement || ''}</div>
                <div><strong>Code postal:</strong> ${form.CodePostal || ''}</div>
                <div><strong>Ville:</strong> ${form.Ville || ''}</div>
                <div><strong>Pays:</strong> ${form.Pays || ''}</div>
                <div><strong>N° Sécurité sociale:</strong> ${form.NumSecu || ''}</div>
            </div>
        `;

        detailSection.style.display = 'block';
    } catch (err) {
        alert('Erreur: ' + err.message);
    }
}

function closeDetail() {
    document.getElementById('formDetail').style.display = 'none';
}

function logout() {
    document.cookie = 'jwt=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;';
    window.location.href = '/login';
}