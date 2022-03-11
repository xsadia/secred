CREATE TABLE warehouse_items (
    id uuid DEFAULT uuid_generate_v4() PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE,
    quantity int DEFAULT 0,
    min int DEFAULT 0,
    max int DEFAULT 0 
);