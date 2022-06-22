-- SQL dump generated using DBML (dbml-lang.org)
-- Database: PostgreSQL
-- Generated at: 2022-06-22T08:09:30.995Z

CREATE TABLE "Users" (
  "id" bigserial UNIQUE PRIMARY KEY,
  "name" varchar NOT NULL,
  "email" varchar NOT NULL,
  "hashed_pw" varchar NOT NULL,
  "image" varchar,
  "status" varchar,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "Message" (
  "id" bigserial UNIQUE PRIMARY KEY,
  "from" varchar NOT NULL,
  "content" varchar NOT NULL,
  "created_at" timestamptz DEFAULT (now()),
  "conv_id" bigint
);

CREATE TABLE "Conversation" (
  "id" bigserial UNIQUE PRIMARY KEY,
  "name" varchar
);

CREATE TABLE "user_conversation" (
  "id" bigserial PRIMARY KEY,
  "user_id" bigint,
  "conv_id" bigint
);

CREATE TABLE "sessions" (
  "id" uuid PRIMARY KEY,
  "email" varchar NOT NULL,
  "user_id" bigint NOT NULL,
  "refresh_token" varchar NOT NULL,
  "user_agent" varchar NOT NULL,
  "client_ip" varchar NOT NULL,
  "is_blocked" bool NOT NULL DEFAULT (false),
  "expires_at" timestamptz NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE INDEX ON "user_conversation" ("user_id", "conv_id");

ALTER TABLE "Message" ADD FOREIGN KEY ("conv_id") REFERENCES "Conversation" ("id");

ALTER TABLE "user_conversation" ADD FOREIGN KEY ("user_id") REFERENCES "Users" ("id");

ALTER TABLE "user_conversation" ADD FOREIGN KEY ("conv_id") REFERENCES "Conversation" ("id");

ALTER TABLE "sessions" ADD FOREIGN KEY ("email") REFERENCES "Users" ("email");

ALTER TABLE "sessions" ADD FOREIGN KEY ("user_id") REFERENCES "Users" ("id");
