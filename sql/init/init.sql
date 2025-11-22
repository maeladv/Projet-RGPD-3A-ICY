CREATE TABLE answers (
    id SERIAL PRIMARY KEY,

    nom VARCHAR(50) NOT NULL,
    prenom VARCHAR(50) NOT NULL,

    date_naissance DATE NOT NULL,
    ville_naissance VARCHAR(50) NOT NULL,

    niveau_diplome VARCHAR(100) NOT NULL,

    mail VARCHAR(100) NOT NULL,

    adresse VARCHAR(100) NOT NULL,
    complement_adresse VARCHAR(100),
    code_postal VARCHAR(100) NOT NULL,
    ville VARCHAR(100) NOT NULL,
    pays VARCHAR(100) NOT NULL,

    num_secu_sociale VARCHAR(15) NOT NULL,
    num_telephone VARCHAR(12) NOT NULL
);
