BEGIN;

-- Create the 'chats' table
CREATE TABLE IF NOT EXISTS chats (
                                     id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                                     user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                                     title TEXT NOT NULL,
                                     summary TEXT,
                                     created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Add an index on user_id for faster lookups
CREATE INDEX IF NOT EXISTS idx_chats_user_id ON chats(user_id);

-- Create the 'messages' table
CREATE TABLE IF NOT EXISTS messages (
                                        id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                                        chat_id UUID NOT NULL REFERENCES chats(id) ON DELETE CASCADE,
                                        role TEXT NOT NULL, -- "user" or "assistant"
                                        content TEXT NOT NULL,
                                        metadata JSONB,
                                        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Add an index on chat_id for faster retrieval of messages for a chat
CREATE INDEX IF NOT EXISTS idx_messages_chat_id ON messages(chat_id);

COMMIT;
