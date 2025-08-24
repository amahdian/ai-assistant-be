BEGIN;

-- Drop indexes first to avoid dependency issues
DROP INDEX IF EXISTS idx_messages_chat_id;
DROP INDEX IF EXISTS idx_chats_user_id;

-- Drop tables
DROP TABLE IF EXISTS messages;
DROP TABLE IF EXISTS chats;

COMMIT;
