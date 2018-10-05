CREATE TABLE suppliers_preparations
(
    id SERIAL PRIMARY KEY,
    preparation_id  int REFERENCES preparations (id) ON UPDATE CASCADE ON DELETE CASCADE,
    supplier_id  int REFERENCES suppliers (id) ON UPDATE CASCADE ON DELETE CASCADE,
    price REAL
);