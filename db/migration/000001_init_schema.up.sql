CREATE TABLE "Users" (
  "id" bigserial UNIQUE PRIMARY KEY,
  "name" varchar NOT NULL,
  "email" varchar NOT NULL,
  "hashed_pw" varchar NOT NULL,
  "image" varchar,
  "status" varchar,
  "created_at" timestamptz DEFAULT (now())
);

CREATE TABLE "Message" (
  "id" bigserial UNIQUE PRIMARY KEY,
  "user_id" bigint NOT NULL,
  "content" varchar,
  "created_at" timestamptz DEFAULT (now())
);

CREATE TABLE "Conversation" (
  "id" bigserial UNIQUE PRIMARY KEY,
  "unread" int DEFAULT 0,
  "last" bigint,
  "messages" bigint
);

CREATE TABLE "user_conversation" (
  "user_id" bigint,
  "conv_id" bigint,
  PRIMARY KEY ("user_id", "conv_id")
);

ALTER TABLE "Message" ADD FOREIGN KEY ("user_id") REFERENCES "Users" ("id");

ALTER TABLE "user_conversation" ADD FOREIGN KEY ("user_id") REFERENCES "Users" ("id");

ALTER TABLE "user_conversation" ADD FOREIGN KEY ("conv_id") REFERENCES "Conversation" ("id");

ALTER TABLE "Conversation" ADD FOREIGN KEY ("last") REFERENCES "Message" ("id");

ALTER TABLE "Conversation" ADD FOREIGN KEY ("messages") REFERENCES "Message" ("id");
