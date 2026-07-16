CREATE TABLE users (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  email varchar(255) UNIQUE NOT NULL,
  password_hash varchar(255) NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);