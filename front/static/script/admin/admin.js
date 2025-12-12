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
        if (!Array.isArray(allForms)) allForms = [];
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

function showNotification(message, type = "info") {
    const notif = document.getElementById('notification');
    notif.textContent = message;
    notif.style.display = 'block';
    notif.style.background = type === "error" ? "#f8d7da" : "#d1e7dd";
    notif.style.color = type === "error" ? "#842029" : "#0f5132";
    notif.style.border = "1px solid " + (type === "error" ? "#f5c2c7" : "#badbcc");
    notif.style.padding = "10px";
    notif.style.margin = "10px 0";
    notif.style.borderRadius = "5px";
    setTimeout(() => { notif.style.display = 'none'; }, 4000);
}

function showConfirmDialog(message, onConfirm) {
    const dialog = document.getElementById('confirmDialog');
    const msg = document.getElementById('confirmDialogMessage');
    const yesBtn = document.getElementById('confirmDialogYes');
    const noBtn = document.getElementById('confirmDialogNo');
    msg.textContent = message;
    dialog.style.display = 'flex';

    // Nettoie les anciens listeners
    yesBtn.onclick = null;
    noBtn.onclick = null;

    yesBtn.onclick = () => {
        dialog.style.display = 'none';
        onConfirm(true);
    };
    noBtn.onclick = () => {
        dialog.style.display = 'none';
        onConfirm(false);
    };
}

function displayForms(forms) {
    if(!Array.isArray(forms)) {
        forms = [];
    }
    const tbody = document.getElementById('formsBody');
    tbody.innerHTML = '';

    if (forms.length === 0) {
        tbody.innerHTML = '<tr><td colspan="8">Aucun formulaire trouvé</td></tr>';
        return;
    }

    forms.forEach(form => {
        const row = document.createElement('tr');
        row.innerHTML = `
            <td>${form.id}</td>
            <td>${form.nom || ''}</td>
            <td>${form.prenom || ''}</td>
            <td>${form.mail || ''}</td>
            <td>${form.num_telephone || ''}</td>
            <td>${formatDate(form.date_naissance)}</td>
            <td>${form.niveau_diplome || ''}</td>
            <td>
                <div class="actions-cell">
                    <button class="viewBtn" data-id="${form.id}">Éditer</button>
                    <button class="deleteBtn" data-id="${form.id}">Suppr.</button>
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
            showConfirmDialog('Êtes-vous sûr de vouloir supprimer ce formulaire ?', (confirmed) => {
                if (confirmed) deleteForm(id);
            });
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
            (form.nom && form.nom.toLowerCase().includes(globalTerm)) ||
            (form.prenom && form.prenom.toLowerCase().includes(globalTerm)) ||
            (form.mail && form.mail.toLowerCase().includes(globalTerm)) ||
            (form.num_telephone && form.num_telephone.toLowerCase().includes(globalTerm));

        // Filtres spécifiques
        const matchName = !searchName || (form.nom && form.nom.toLowerCase().includes(searchName));
        const matchSurname = !searchSurname || (form.prenom && form.prenom.toLowerCase().includes(searchSurname));
        const matchEmail = !searchEmail || (form.mail && form.mail.toLowerCase().includes(searchEmail));
        const matchPhone = !searchPhone || (form.num_telephone && form.num_telephone.toLowerCase().includes(searchPhone));
        
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
                <input type="hidden" name="id" value="${form.id}">
                <div class="detail-grid">
                    <div class="form-group">
            });
                        <label>Nom</label>
                        <input type="text" name="nom" value="${form.nom || ''}">
                    </div>
                    <div class="form-group">
                        <label>Prénom</label>
                        <input type="text" name="prenom" value="${form.prenom || ''}">
                    </div>
                    <div class="form-group">
                        <label>Date de naissance</label>
                        <input type="date" name="date_naissance" value="${formatDateForInput(form.date_naissance)}">
                    </div>
                    <div class="form-group">
                        <label>Ville de naissance</label>
                        <input type="text" name="ville_naissance" value="${form.ville_naissance || ''}">
                    </div>
                    <div class="form-group">
                        <label>Niveau diplôme</label>
                        <input type="text" name="niveau_diplome" value="${form.niveau_diplome || ''}">
                    </div>
                    <div class="form-group">
                        <label>Email</label>
                        <input type="email" name="mail" value="${form.mail || ''}">
                    </div>
                    <div class="form-group">
                        <label>Téléphone</label>
                        <input type="tel" name="num_telephone" value="${form.num_telephone || ''}">
                    </div>
                    <div class="form-group">
                        <label>Adresse</label>
                        <input type="text" name="adresse" value="${form.adresse || ''}">
                    </div>
                    <div class="form-group">
                        <label>Complément</label>
                        <input type="text" name="complement" value="${form.complement || ''}">
                    </div>
                    <div class="form-group">
                        <label>Code postal</label>
                        <input type="text" name="code_postal" value="${form.code_postal || ''}">
                    </div>
                    <div class="form-group">
                        <label>Ville</label>
                        <input type="text" name="ville" value="${form.ville || ''}">
                    </div>
                    <div class="form-group">
                        <label>Pays</label>
                        <input type="text" name="pays" value="${form.pays || ''}">
                    </div>
                    <div class="form-group">
                        <label>N° Sécurité sociale</label>
                        <input type="text" name="num_secu_sociale" value="${form.num_secu_sociale || ''}">
                    </div>
                </div>
                <div class="modal-actions" style="margin-top: 20px; text-align: right;">
                    <button type="button" class="button primary" onclick="updateForm()">Sauvegarder les modifications</button>
                </div>
            </form>
        `;

        detailSection.style.display = 'block';
    } catch (err) {
        showNotification('Erreur: ' + err.message, "error");
    }
}

async function updateForm() {
    const formElement = document.getElementById('editForm');
    const formData = new FormData(formElement);
    const data = {
        id: parseInt(formData.get('id')),
        nom: formData.get('nom'),
        prenom: formData.get('prenom'),
        date_naissance: new Date(formData.get('date_naissance')).toISOString(),
        ville_naissance: formData.get('ville_naissance'),
        niveau_diplome: formData.get('niveau_diplome'),
        mail: formData.get('mail'),
        num_telephone: formData.get('num_telephone'),
        adresse: formData.get('adresse'),
        complement: formData.get('complement'),
        code_postal: formData.get('code_postal'),
        ville: formData.get('ville'),
        pays: formData.get('pays'),
        num_secu_sociale: formData.get('num_secu_sociale')
    };

    try {
        const response = await fetch(`${API_URL}/form/modify`, {
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

        showNotification('Formulaire mis à jour avec succès', "info");
        closeDetail();
        loadForms();
    } catch (err) {
        showNotification('Erreur: ' + err.message, "error");
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

        showNotification('Formulaire supprimé avec succès', "info");
        loadForms(); // Recharger la liste
    } catch (err) {
        showNotification('Erreur: ' + err.message, "error");
    }
}

function closeDetail() {
    document.getElementById('formDetail').style.display = 'none';
}
