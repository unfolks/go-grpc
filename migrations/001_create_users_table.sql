-- Create users table
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    role TEXT NOT NULL DEFAULT 'user',
    attributes JSONB DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Seed some initial users for testing
-- Note: In a real app, passwords should be properly hashed. 
-- For this setup, we'll assume the hash is 'admin123' and 'user123' raw for demonstration,
-- but the implementation will use bcrypt or similar.
INSERT INTO users (username, password_hash, role, attributes)
VALUES 
    ('admin', 'admin123', 'admin', '{"full_name": "System Administrator"}'),
    ('user1', 'user123', 'user', '{"owner_id": "1", "department": "sales"}')
ON CONFLICT (username) DO NOTHING;
