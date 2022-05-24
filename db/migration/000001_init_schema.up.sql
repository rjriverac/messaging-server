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
  "created_at" timestamptz not null DEFAULT (now()),
  "conv_id" bigint not null
);

CREATE TABLE "Conversation" (
  "id" bigserial UNIQUE PRIMARY KEY,
  "name" varchar
);

CREATE TABLE "user_conversation" (
  "id" bigserial PRIMARY KEY,
  "user_id" bigint not null,
  "conv_id" bigint not null
);

CREATE INDEX ON "user_conversation" ("user_id", "conv_id");

ALTER TABLE "Message" ADD FOREIGN KEY ("conv_id") REFERENCES "Conversation" ("id");

ALTER TABLE "user_conversation" ADD FOREIGN KEY ("user_id") REFERENCES "Users" ("id");

ALTER TABLE "user_conversation" ADD FOREIGN KEY ("conv_id") REFERENCES "Conversation" ("id");
