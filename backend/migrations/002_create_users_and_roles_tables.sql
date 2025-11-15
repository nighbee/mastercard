-- Create roles table
CREATE TABLE IF NOT EXISTS roles (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) UNIQUE NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create permissions table
CREATE TABLE IF NOT EXISTS permissions (
    id SERIAL PRIMARY KEY,
    resource VARCHAR(100) NOT NULL,
    action VARCHAR(50) NOT NULL,
    conditions JSONB,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(resource, action)
);

-- Create role_permissions junction table
CREATE TABLE IF NOT EXISTS role_permissions (
    role_id INTEGER NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    permission_id INTEGER NOT NULL REFERENCES permissions(id) ON DELETE CASCADE,
    PRIMARY KEY (role_id, permission_id)
);

-- Create users table
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    full_name VARCHAR(255) NOT NULL,
    role_id INTEGER REFERENCES roles(id) ON DELETE SET NULL,
    is_active BOOLEAN DEFAULT TRUE,
    last_login TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_role_id ON users(role_id);
CREATE INDEX IF NOT EXISTS idx_users_is_active ON users(is_active);
CREATE INDEX IF NOT EXISTS idx_role_permissions_role_id ON role_permissions(role_id);
CREATE INDEX IF NOT EXISTS idx_role_permissions_permission_id ON role_permissions(permission_id);

-- Insert default roles
INSERT INTO roles (name, description) VALUES
    ('admin', 'Administrator with full system access'),
    ('manager', 'Manager with access to all data and reports'),
    ('analyst', 'Analyst with read access to transaction data'),
    ('viewer', 'Viewer with limited read-only access')
ON CONFLICT (name) DO NOTHING;

-- Insert default permissions
INSERT INTO permissions (resource, action, conditions) VALUES
    ('transactions', 'read', '{}'),
    ('transactions', 'read_limited', '{"columns": ["date", "merch_name", "trx_amount_usd", "location_city"]}'),
    ('conversations', 'create', '{}'),
    ('conversations', 'read_own', '{}'),
    ('conversations', 'read_all', '{}'),
    ('conversations', 'update_own', '{}'),
    ('conversations', 'delete_own', '{}'),
    ('conversations', 'delete_all', '{}'),
    ('users', 'read', '{}'),
    ('users', 'create', '{}'),
    ('users', 'update', '{}'),
    ('users', 'delete', '{}'),
    ('audit_logs', 'read', '{}'),
    ('audit_logs', 'export', '{}')
ON CONFLICT (resource, action) DO NOTHING;

-- Assign permissions to roles
-- Admin: all permissions
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r, permissions p
WHERE r.name = 'admin'
ON CONFLICT DO NOTHING;

-- Manager: read transactions, manage own conversations, read audit logs
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r, permissions p
WHERE r.name = 'manager'
  AND (p.resource = 'transactions' AND p.action = 'read')
  OR (p.resource = 'conversations' AND p.action IN ('create', 'read_own', 'update_own', 'delete_own'))
  OR (p.resource = 'audit_logs' AND p.action = 'read')
ON CONFLICT DO NOTHING;

-- Analyst: read transactions, manage own conversations
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r, permissions p
WHERE r.name = 'analyst'
  AND (p.resource = 'transactions' AND p.action = 'read')
  OR (p.resource = 'conversations' AND p.action IN ('create', 'read_own', 'update_own', 'delete_own'))
ON CONFLICT DO NOTHING;

-- Viewer: limited read access
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r, permissions p
WHERE r.name = 'viewer'
  AND (p.resource = 'transactions' AND p.action = 'read_limited')
  OR (p.resource = 'conversations' AND p.action IN ('create', 'read_own'))
ON CONFLICT DO NOTHING;

COMMENT ON TABLE roles IS 'User roles for RBAC';
COMMENT ON TABLE permissions IS 'Permissions defining what actions can be performed on resources';
COMMENT ON TABLE role_permissions IS 'Many-to-many relationship between roles and permissions';
COMMENT ON TABLE users IS 'User accounts with authentication and role assignment';

