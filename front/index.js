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

app.use('/static', express.static(path.join(__dirname, 'static')));
app.use('/script', express.static(path.join(__dirname, 'script')));

// req et res correspondent à la requête et la réponse HTTP

app.get('/', (req, res) => {
    res.sendFile(path.join(__dirname, '/static/pages/formulaire.html'));
});

// css files
app.get('/css/style.css', (req, res) => {
    res.sendFile(path.join(__dirname, '/static/css/style.css'))
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
    const filePath = path.join(__dirname, 'static/images', fileName)
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