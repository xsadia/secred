CREATE TABLE users (
    id uuid DEFAULT uuid_generate_v4() PRIMARY KEY, 
    email VARCHAR(255) NOT NULL,
    username VARCHAR(80) NOT NULL,
    password VARCHAR(255) NOT NULL,
    refresh_token VARCHAR(255) 
);