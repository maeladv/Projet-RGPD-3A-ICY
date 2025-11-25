CREATE TABLE IF NOT EXISTS answers (
    id SERIAL PRIMARY KEY,

    nom VARCHAR(100) NOT NULL,
    prenom VARCHAR(100) NOT NULL,

    date_naissance DATE NOT NULL,
    ville_naissance VARCHAR(100) NOT NULL,

    niveau_diplome VARCHAR(100) NOT NULL,

    mail VARCHAR(255) NOT NULL,

    adresse VARCHAR(255) NOT NULL,
    complement_adresse VARCHAR(255),
    code_postal VARCHAR(20) NOT NULL,
    ville VARCHAR(100) NOT NULL,
    pays VARCHAR(100) NOT NULL,

    num_secu_sociale VARCHAR(50) NOT NULL,
    num_telephone VARCHAR(50) NOT NULL
);