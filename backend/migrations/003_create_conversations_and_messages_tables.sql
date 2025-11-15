-- Create conversations table
CREATE TABLE IF NOT EXISTS conversations (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title VARCHAR(255),
    parent_branch_id INTEGER REFERENCES conversations(id) ON DELETE SET NULL,
    branch_point_message_id INTEGER,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create messages table
CREATE TABLE IF NOT EXISTS messages (
    id SERIAL PRIMARY KEY,
    conversation_id INTEGER NOT NULL REFERENCES conversations(id) ON DELETE CASCADE,
    user_message TEXT NOT NULL,
    sql_query TEXT,
    result_data JSONB,
    result_format VARCHAR(20) CHECK (result_format IN ('text', 'table', 'chart', 'error')),
    error_message TEXT,
    execution_time_ms INTEGER,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_conversations_user_id ON conversations(user_id);
CREATE INDEX IF NOT EXISTS idx_conversations_parent_branch_id ON conversations(parent_branch_id);
CREATE INDEX IF NOT EXISTS idx_conversations_created_at ON conversations(created_at);
CREATE INDEX IF NOT EXISTS idx_messages_conversation_id ON messages(conversation_id);
CREATE INDEX IF NOT EXISTS idx_messages_created_at ON messages(created_at);

-- Create full-text search index for conversations
CREATE INDEX IF NOT EXISTS idx_conversations_title_search ON conversations USING gin(to_tsvector('english', COALESCE(title, '')));

-- Create full-text search index for messages
CREATE INDEX IF NOT EXISTS idx_messages_user_message_search ON messages USING gin(to_tsvector('english', user_message));

COMMENT ON TABLE conversations IS 'Chat conversation sessions with support for branching';
COMMENT ON TABLE messages IS 'Individual messages within conversations, storing queries and results';
COMMENT ON COLUMN conversations.parent_branch_id IS 'Reference to parent conversation when this is a branch';
COMMENT ON COLUMN conversations.branch_point_message_id IS 'Message ID where the branch was created';
COMMENT ON COLUMN messages.result_format IS 'Format of the result: text, table, chart, or error';

