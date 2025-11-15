-- Enable Row Level Security (RLS) on tables
ALTER TABLE transactions ENABLE ROW LEVEL SECURITY;
ALTER TABLE conversations ENABLE ROW LEVEL SECURITY;
ALTER TABLE messages ENABLE ROW LEVEL SECURITY;
ALTER TABLE audit_logs ENABLE ROW LEVEL SECURITY;

-- RLS Policies for transactions table
-- Policy: Users can read transactions based on their role permissions
CREATE POLICY transactions_read_policy ON transactions
    FOR SELECT
    USING (
        -- Admin and Manager can see all transactions
        EXISTS (
            SELECT 1 FROM users u
            JOIN roles r ON u.role_id = r.id
            WHERE u.id = current_setting('app.current_user_id')::INTEGER
            AND r.name IN ('admin', 'manager')
        )
        OR
        -- Analyst can see all transactions
        EXISTS (
            SELECT 1 FROM users u
            JOIN roles r ON u.role_id = r.id
            WHERE u.id = current_setting('app.current_user_id')::INTEGER
            AND r.name = 'analyst'
        )
        OR
        -- Viewer can see transactions (limited columns enforced at application level)
        EXISTS (
            SELECT 1 FROM users u
            JOIN roles r ON u.role_id = r.id
            WHERE u.id = current_setting('app.current_user_id')::INTEGER
            AND r.name = 'viewer'
        )
    );

-- RLS Policies for conversations table
-- Policy: Users can only see their own conversations, except admins who can see all
CREATE POLICY conversations_read_policy ON conversations
    FOR SELECT
    USING (
        user_id = current_setting('app.current_user_id')::INTEGER
        OR
        EXISTS (
            SELECT 1 FROM users u
            JOIN roles r ON u.role_id = r.id
            WHERE u.id = current_setting('app.current_user_id')::INTEGER
            AND r.name = 'admin'
        )
    );

CREATE POLICY conversations_insert_policy ON conversations
    FOR INSERT
    WITH CHECK (
        user_id = current_setting('app.current_user_id')::INTEGER
    );

CREATE POLICY conversations_update_policy ON conversations
    FOR UPDATE
    USING (
        user_id = current_setting('app.current_user_id')::INTEGER
        OR
        EXISTS (
            SELECT 1 FROM users u
            JOIN roles r ON u.role_id = r.id
            WHERE u.id = current_setting('app.current_user_id')::INTEGER
            AND r.name = 'admin'
        )
    );

CREATE POLICY conversations_delete_policy ON conversations
    FOR DELETE
    USING (
        user_id = current_setting('app.current_user_id')::INTEGER
        OR
        EXISTS (
            SELECT 1 FROM users u
            JOIN roles r ON u.role_id = r.id
            WHERE u.id = current_setting('app.current_user_id')::INTEGER
            AND r.name = 'admin'
        )
    );

-- RLS Policies for messages table
-- Policy: Users can only see messages from their own conversations
CREATE POLICY messages_read_policy ON messages
    FOR SELECT
    USING (
        EXISTS (
            SELECT 1 FROM conversations c
            WHERE c.id = messages.conversation_id
            AND (
                c.user_id = current_setting('app.current_user_id')::INTEGER
                OR
                EXISTS (
                    SELECT 1 FROM users u
                    JOIN roles r ON u.role_id = r.id
                    WHERE u.id = current_setting('app.current_user_id')::INTEGER
                    AND r.name = 'admin'
                )
            )
        )
    );

CREATE POLICY messages_insert_policy ON messages
    FOR INSERT
    WITH CHECK (
        EXISTS (
            SELECT 1 FROM conversations c
            WHERE c.id = messages.conversation_id
            AND c.user_id = current_setting('app.current_user_id')::INTEGER
        )
    );

-- RLS Policies for audit_logs table
-- Policy: Users can see their own audit logs, admins can see all
CREATE POLICY audit_logs_read_policy ON audit_logs
    FOR SELECT
    USING (
        user_id = current_setting('app.current_user_id')::INTEGER
        OR
        EXISTS (
            SELECT 1 FROM users u
            JOIN roles r ON u.role_id = r.id
            WHERE u.id = current_setting('app.current_user_id')::INTEGER
            AND r.name = 'admin'
        )
    );

COMMENT ON POLICY transactions_read_policy ON transactions IS 'RLS policy controlling read access to transactions based on user role';
COMMENT ON POLICY conversations_read_policy ON conversations IS 'RLS policy allowing users to read their own conversations';
COMMENT ON POLICY messages_read_policy ON messages IS 'RLS policy allowing users to read messages from their conversations';

