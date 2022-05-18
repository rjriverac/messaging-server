ALTER TABLE "Conversation"
  ALTER "Messages" DROP DEFAULT,
  ALTER "Messages" type bigint[] USING ARRAY["Messages"],
  ALTER SET DEFAULT "{}";
