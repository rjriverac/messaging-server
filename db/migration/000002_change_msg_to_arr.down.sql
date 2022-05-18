ALTER TABLE "Conversation"
ALTER COLUMN "Messages" DROP DEFAULT,
  ALTER COLUMN "Messages" type bigint USING "Messages"[1]::bigint;