CREATE TABLE "Users" (
  "id" int UNIQUE PRIMARY KEY,
  "name" varchar NOT NULL,
  "email" varchar NOT NULL,
  "hashed_pw" varchar NOT NULL,
  "image" varchar,
  "status" varchar,
  "created_at" timestamptz,
);

CREATE TABLE "Message" (
  "id" int UNIQUE PRIMARY KEY,
  "user_id" int,
  "content" varchar,
  "created_at" timestamptz
);

CREATE TABLE "Conversation" (
  "id" int UNIQUE PRIMARY KEY,
  "unread" int DEFAULT 0,
  "last" int,
  "messages" int
);

CREATE TABLE "user_conversation" (
  "user_id" int,
  "conv_id" int
);

ALTER TABLE "Message" ADD FOREIGN KEY ("user_id") REFERENCES "Users" ("message");

ALTER TABLE "Users" ADD FOREIGN KEY ("id") REFERENCES "user_conversation" ("user_id");

ALTER TABLE "Conversation" ADD FOREIGN KEY ("id") REFERENCES "user_conversation" ("conv_id");

ALTER TABLE "Message" ADD PRIMARY KEY ("user_id","conv_id");
ALTER TABLE "Message" ADD FOREIGN KEY ("id") REFERENCES "Conversation" ("last");

ALTER TABLE "Message" ADD FOREIGN KEY ("id") REFERENCES "Conversation" ("messages");
