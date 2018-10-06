CREATE TABLE Preparations
(
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT,
    type TEXT NOT NULL,
    activeIngredient TEXT,
    imageURL TEXT
);