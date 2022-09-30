BEGIN TRANSACTION;

CREATE TABLE "users" (
  "id" INTEGER PRIMARY KEY NOT NULL,
  "name" varchar(32) DEFAULT NULL,
  "digest" varchar(200),
  "numentry" integer(11) NOT NULL DEFAULT 0,
  "nopinlist" integer(11) NOT NULL DEFAULT 0,
  "numsubstr" integer(11) NOT NULL DEFAULT 0,
  "autoseen" boolean NOT NULL DEFAULT 0,
  "last_login" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX "name" ON "users" ("name");

CREATE TABLE "category" (
  "id" INTEGER PRIMARY KEY NOT NULL,
  "user_id" integer(11) NOT NULL,
  "name" varchar(60) NOT NULL,
  FOREIGN KEY ("user_id") REFERENCES "users"("id") ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE INDEX "user_id" ON "category" ("user_id");

CREATE TABLE "feed" (
  "id" INTEGER PRIMARY KEY NOT NULL,
  "url" text NOT NULL,
  "siteurl" text NOT NULL,
  "title" varchar(200) NOT NULL,
  "http_status" varchar(3) NOT NULL,
  "pubdate" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "term" varchar(1) NOT NULL DEFAULT '1',
  "cache" text NOT NULL,
  "next_serial" integer(11) NOT NULL DEFAULT 0
);

CREATE TABLE "subscription" (
  "id" INTEGER PRIMARY KEY NOT NULL,
  "category_id" integer(11) NOT NULL,
  "feed_id" integer(11) NOT NULL,
  "user_id" integer(11) NOT NULL,
  FOREIGN KEY ("category_id") REFERENCES "category"("id") ON DELETE CASCADE ON UPDATE CASCADE,
  FOREIGN KEY ("feed_id") REFERENCES "feed"("id") ON DELETE CASCADE ON UPDATE CASCADE,
  FOREIGN KEY ("user_id") REFERENCES "users"("id") ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE INDEX "category_id" ON "subscription" ("category_id");

CREATE INDEX "feed_id" ON "subscription" ("feed_id");

CREATE INDEX "subscription_idx" ON "subscription" ("user_id");

CREATE TABLE "entry" (
  "serial" integer(11) NOT NULL,
  "pubdate" datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "update_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "readflag" tinyint(4) NOT NULL,
  "subscription_id" integer(11) NOT NULL,
  "feed_id" integer(11) NOT NULL,
  "user_id" integer(11) NOT NULL,
  FOREIGN KEY ("subscription_id") REFERENCES "subscription"("id") ON DELETE CASCADE ON UPDATE CASCADE,
  FOREIGN KEY ("user_id") REFERENCES "users"("id") ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE INDEX "subscription_id" ON "entry" ("subscription_id");

CREATE INDEX "entry_idx" ON "entry" ("user_id");

CREATE UNIQUE INDEX "serial_2" ON "entry" ("serial", "feed_id", "user_id");

CREATE TABLE "story" (
  "feed_id" integer(11) NOT NULL,
  "serial" integer(11) NOT NULL,
  "title" varchar(80) NOT NULL,
  "description" tinytext(255) NOT NULL,
  "url" text NOT NULL,
  PRIMARY KEY ("feed_id", "serial")
);

COMMIT;
