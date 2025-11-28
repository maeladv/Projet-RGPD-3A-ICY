// ######    CONFIGURATION DU PROJET    #####

const port = 2021;

// module express pour gérer les routes
const express = require("express");
const app = express();
// extension pour parser le json
app.use(express.json());
app.use(express.urlencoded({ extended: true }));

// path
const path = require("path");



// #####  ROUTES ET REQUETES  #####

// req et res correspondent à la requête et la réponse HTTP

//  Page d'affichage des scores
app.get('/', (req, res) => {
    res.sendFile(path.join(__dirname, '/static/pages/formulaire.html'));
});

// css files
app.get('/css/scores.css', (req, res) => {
    res.sendFile(path.join(__dirname, '/static/scores.css'))
})

// Polices d'écriture : tous les fichiers du dossier font sont publics
app.get('/fonts/:file', (req, res) => {
    const fileName = req.params.file;
    const filePath = path.join(__dirname, 'static/fonts', fileName)
    if (fs.existsSync(filePath)) {
        res.sendFile(filePath);
    } else {
        res.status(404).send('Fichier non trouvé');
    }
})


// images
app.get('/images/:file', (req, res) => {
    const fileName = req.params.file;
    const filePath = path.join(__dirname, 'static/img', fileName)
    if (fs.existsSync(filePath)) {
        res.sendFile(filePath);
    } else {
        res.status(404).send('Image non trouvée');
    }
})



// Démarrer le serveur
app.listen(port, () => {
    console.log(` Serveur Lancé sur http://localhost:${port}`);
});