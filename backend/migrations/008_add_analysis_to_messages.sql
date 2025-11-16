-- Add analysis column to messages table for conversational insights
ALTER TABLE messages ADD COLUMN IF NOT EXISTS analysis TEXT;

COMMENT ON COLUMN messages.analysis IS 'Conversational analysis and insights about the query results';

