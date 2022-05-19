ALTER TABLE "Conversation"
ALTER COLUMN last
set not null;
ALTER TABLE "Conversation"
ALTER COLUMN messages
set not null;