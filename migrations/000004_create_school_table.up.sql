CREATE TABLE schools (
  id uuid DEFAULT uuid_generate_v4() PRIMARY KEY,
  name VARCHAR(255) UNIQUE 
)