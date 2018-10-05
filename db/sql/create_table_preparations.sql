CREATE TABLE Preparations
(
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT,
    activeIngredient TEXT
);