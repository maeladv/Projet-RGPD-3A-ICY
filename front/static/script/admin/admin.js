const API_URL = '/api';

let allForms = [];
let filteredForms = [];

// Vérification de l'authentification au chargement
document.addEventListener('DOMContentLoaded', () => {
    // On ne vérifie plus le cookie ici car il est HttpOnly
    // C'est loadForms qui gérera la redirection si 401
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

// checkAuth supprimé car inutile avec cookie HttpOnly
// La vérification se fait via les appels API

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
                window.location.href = '/login.html';
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
                <div class="actions-cell">
                    <button class="viewBtn" data-id="${form.ID || form.id}">Éditer</button>
                    <button class="deleteBtn" data-id="${form.ID || form.id}">Suppr.</button>
                </div>
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

    // Ajoute les événements sur les boutons "Supprimer"
    document.querySelectorAll('.deleteBtn').forEach(btn => {
        btn.addEventListener('click', (e) => {
            const id = e.target.dataset.id;
            if(confirm('Êtes-vous sûr de vouloir supprimer ce formulaire ?')) {
                deleteForm(id);
            }
        });
    });
}

function formatDate(dateStr) {
    if (!dateStr) return '';
    const date = new Date(dateStr);
    return date.toLocaleDateString('fr-FR');
}

function formatDateForInput(dateStr) {
    if (!dateStr) return '';
    const date = new Date(dateStr);
    return date.toISOString().split('T')[0];
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
            <form id="editForm" onsubmit="return false;">
                <input type="hidden" name="ID" value="${form.ID || form.id}">
                <div class="detail-grid">
                    <div class="form-group">
                        <label>Nom</label>
                        <input type="text" name="Nom" value="${form.Nom || ''}">
                    </div>
                    <div class="form-group">
                        <label>Prénom</label>
                        <input type="text" name="Prenom" value="${form.Prenom || ''}">
                    </div>
                    <div class="form-group">
                        <label>Date de naissance</label>
                        <input type="date" name="DateNaissance" value="${formatDateForInput(form.DateNaissance)}">
                    </div>
                    <div class="form-group">
                        <label>Ville de naissance</label>
                        <input type="text" name="VilleNaissance" value="${form.VilleNaissance || ''}">
                    </div>
                    <div class="form-group">
                        <label>Niveau diplôme</label>
                        <input type="text" name="NiveauDiplome" value="${form.NiveauDiplome || ''}">
                    </div>
                    <div class="form-group">
                        <label>Email</label>
                        <input type="email" name="Mail" value="${form.Mail || ''}">
                    </div>
                    <div class="form-group">
                        <label>Téléphone</label>
                        <input type="tel" name="Telephone" value="${form.Telephone || ''}">
                    </div>
                    <div class="form-group">
                        <label>Adresse</label>
                        <input type="text" name="Adresse" value="${form.Adresse || ''}">
                    </div>
                    <div class="form-group">
                        <label>Complément</label>
                        <input type="text" name="Complement" value="${form.Complement || ''}">
                    </div>
                    <div class="form-group">
                        <label>Code postal</label>
                        <input type="text" name="CodePostal" value="${form.CodePostal || ''}">
                    </div>
                    <div class="form-group">
                        <label>Ville</label>
                        <input type="text" name="Ville" value="${form.Ville || ''}">
                    </div>
                    <div class="form-group">
                        <label>Pays</label>
                        <input type="text" name="Pays" value="${form.Pays || ''}">
                    </div>
                    <div class="form-group">
                        <label>N° Sécurité sociale</label>
                        <input type="text" name="NumSecu" value="${form.NumSecu || ''}">
                    </div>
                </div>
                <div class="modal-actions" style="margin-top: 20px; text-align: right;">
                    <button type="button" class="button primary" onclick="updateForm()">Sauvegarder les modifications</button>
                </div>
            </form>
        `;

        detailSection.style.display = 'block';
    } catch (err) {
        alert('Erreur: ' + err.message);
    }
}

async function updateForm() {
    const formElement = document.getElementById('editForm');
    const formData = new FormData(formElement);
    const data = {};
    
    formData.forEach((value, key) => {
        if (key === 'ID') {
            data[key] = parseInt(value);
        } else if (key === 'DateNaissance') {
            data[key] = new Date(value).toISOString();
        } else {
            data[key] = value;
        }
    });

    try {
        const response = await fetch(`${API_URL}/form/update`, {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(data),
            credentials: 'include'
        });

        if (!response.ok) {
            throw new Error('Erreur lors de la mise à jour');
        }

        alert('Formulaire mis à jour avec succès');
        closeDetail();
        loadForms(); // Recharger la liste
    } catch (err) {
        alert('Erreur: ' + err.message);
    }
}

async function deleteForm(id) {
    try {
        const response = await fetch(`${API_URL}/form/delete?id=${id}`, {
            method: 'DELETE',
            credentials: 'include'
        });

        if (!response.ok) {
            throw new Error('Erreur lors de la suppression');
        }

        alert('Formulaire supprimé avec succès');
        loadForms(); // Recharger la liste
    } catch (err) {
        alert('Erreur: ' + err.message);
    }
}

function closeDetail() {
    document.getElementById('formDetail').style.display = 'none';
}

function logout() {
    document.cookie = 'jwt=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;';
    window.location.href = '/login.html';
}
