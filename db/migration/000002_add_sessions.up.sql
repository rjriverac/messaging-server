CREATE TABLE "sessions" (
  "id" uuid PRIMARY KEY,
  "email" varchar NOT NULL,
  "user_id" BIGINT NOT NULL,
  "refresh_token" varchar NOT NULL,
  "user_agent" VARCHAR not null,
  "client_ip" VARCHAR not null,
  "is_blocked" boolean NOT NULL DEFAULT false,
  "expires_at" timestamptz NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

ALTER TABLE "sessions" ADD FOREIGN KEY ("user_id") REFERENCES "Users" ("id");
ALTER TABLE "sessions" ADD FOREIGN KEY ("email") REFERENCES "Users" ("email");