ALTER TABLE "Conversation"
  ALTER COLUMN "Messages" type bigint[] USING array["Messages"]::BIGINT[],
  ALTER COLUMN SET DEFAULT "{}";
