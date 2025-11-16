-- Update RBAC roles and permissions according to new requirements

-- Update role names and descriptions
UPDATE roles SET 
    name = 'analyzer',
    description = 'Core User: Can perform the main function of the application'
WHERE name = 'analyst';

UPDATE roles SET 
    description = 'Team Lead: Has all Analyzer abilities, plus can manage users and view their activity'
WHERE name = 'manager';

UPDATE roles SET 
    description = 'Super User: Has all Manager abilities, plus full control over system configuration and all users'
WHERE name = 'admin';

-- Delete viewer role if it exists
DELETE FROM role_permissions WHERE role_id IN (SELECT id FROM roles WHERE name = 'viewer');
DELETE FROM users WHERE role_id IN (SELECT id FROM roles WHERE name = 'viewer');
DELETE FROM roles WHERE name = 'viewer';

-- Clear existing role-permission mappings to rebuild
DELETE FROM role_permissions;

-- Ensure all required permissions exist
INSERT INTO permissions (resource, action, conditions) VALUES
    -- Query permissions (main application function)
    ('queries', 'execute', '{}'),
    
    -- Conversation permissions
    ('conversations', 'create', '{}'),
    ('conversations', 'read_own', '{}'),
    ('conversations', 'read_all', '{}'),
    ('conversations', 'update_own', '{}'),
    ('conversations', 'delete_own', '{}'),
    ('conversations', 'delete_all', '{}'),
    
    -- Transaction data permissions
    ('transactions', 'read', '{}'),
    
    -- User management permissions (for Manager and Admin)
    ('users', 'read', '{}'),
    ('users', 'create', '{}'),
    ('users', 'update', '{}'),
    ('users', 'delete', '{}'),
    
    -- Audit log permissions
    ('audit_logs', 'read', '{}'),
    ('audit_logs', 'read_all', '{}'),
    
    -- System configuration permissions (Admin only)
    ('system', 'configure', '{}')
ON CONFLICT (resource, action) DO NOTHING;

-- Assign permissions to Analyzer role
-- Analyzer: Core User - queries and own conversations
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r, permissions p
WHERE r.name = 'analyzer'
  AND (
    (p.resource = 'queries' AND p.action = 'execute') OR
    (p.resource = 'transactions' AND p.action = 'read') OR
    (p.resource = 'conversations' AND p.action IN ('create', 'read_own', 'update_own', 'delete_own'))
  )
ON CONFLICT DO NOTHING;

-- Assign permissions to Manager role
-- Manager: All Analyzer abilities + manage users + view activity
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r, permissions p
WHERE r.name = 'manager'
  AND (
    -- All analyzer permissions
    (p.resource = 'queries' AND p.action = 'execute') OR
    (p.resource = 'transactions' AND p.action = 'read') OR
    (p.resource = 'conversations' AND p.action IN ('create', 'read_own', 'read_all', 'update_own', 'delete_own')) OR
    -- Manager additional permissions
    (p.resource = 'users' AND p.action IN ('read', 'create', 'update')) OR
    (p.resource = 'audit_logs' AND p.action = 'read')
  )
ON CONFLICT DO NOTHING;

-- Assign permissions to Admin role
-- Admin: All Manager abilities + full system control
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r, permissions p
WHERE r.name = 'admin'
  AND (
    -- All manager permissions (inherited)
    (p.resource = 'queries' AND p.action = 'execute') OR
    (p.resource = 'transactions' AND p.action = 'read') OR
    (p.resource = 'conversations' AND p.action IN ('create', 'read_own', 'read_all', 'update_own', 'delete_own', 'delete_all')) OR
    (p.resource = 'users' AND p.action IN ('read', 'create', 'update', 'delete')) OR
    (p.resource = 'audit_logs' AND p.action IN ('read', 'read_all')) OR
    -- Admin exclusive permissions
    (p.resource = 'system' AND p.action = 'configure')
  )
ON CONFLICT DO NOTHING;

COMMENT ON TABLE roles IS 'User roles: analyzer (core user), manager (team lead), admin (super user)';

