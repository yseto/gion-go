CREATE TABLE "users" (
  "id" serial NOT NULL,
  "name" character varying(32) DEFAULT NULL,
  "digest" character varying(200),
  "numentry" bigint DEFAULT 0 NOT NULL,
  "nopinlist" boolean DEFAULT false NOT NULL,
  "numsubstr" bigint DEFAULT 0 NOT NULL,
  "autoseen" boolean DEFAULT false NOT NULL,
  "last_login" timestamp(0) DEFAULT CURRENT_TIMESTAMP NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "name" UNIQUE ("name")
);

CREATE TABLE "category" (
  "id" serial NOT NULL,
  "user_id" bigint NOT NULL,
  "name" character varying(60) NOT NULL,
  PRIMARY KEY ("id")
);
CREATE INDEX "user_id" on "category" ("user_id");

CREATE TABLE "feed" (
  "id" serial NOT NULL,
  "url" text NOT NULL,
  "siteurl" text NOT NULL,
  "title" character varying(200) NOT NULL,
  "http_status" character varying(3) NOT NULL,
  "pubdate" timestamp(0) DEFAULT CURRENT_TIMESTAMP NOT NULL,
  "term" character varying(1) DEFAULT '1' NOT NULL,
  "cache" text NOT NULL,
  "next_serial" bigint DEFAULT 0 NOT NULL,
  PRIMARY KEY ("id")
);

CREATE TABLE "subscription" (
  "id" serial NOT NULL,
  "category_id" bigint NOT NULL,
  "feed_id" bigint NOT NULL,
  "user_id" bigint NOT NULL,
  PRIMARY KEY ("id")
);
CREATE INDEX "category_id" on "subscription" ("category_id");
CREATE INDEX "feed_id" on "subscription" ("feed_id");
CREATE INDEX "subscription_idx_1" on "subscription" ("user_id");

CREATE TABLE "entry" (
  "serial" bigint NOT NULL,
  "pubdate" timestamp(0) DEFAULT CURRENT_TIMESTAMP NOT NULL,
  "update_at" timestamp(0) DEFAULT CURRENT_TIMESTAMP NOT NULL,
  "readflag" smallint NOT NULL,
  "subscription_id" bigint NOT NULL,
  "feed_id" bigint NOT NULL,
  "user_id" bigint NOT NULL,
  CONSTRAINT "serial_2" UNIQUE ("serial", "feed_id", "user_id")
);
CREATE INDEX "subscription_id" on "entry" ("subscription_id");
CREATE INDEX "entry_idx_1" on "entry" ("user_id");

CREATE TABLE "story" (
  "feed_id" bigint NOT NULL,
  "serial" bigint NOT NULL,
  "title" character varying(80) NOT NULL,
  "description" text NOT NULL,
  "url" text NOT NULL,
  PRIMARY KEY ("feed_id", "serial")
);

ALTER TABLE "category" ADD CONSTRAINT "category_ibfk_1" FOREIGN KEY ("user_id")
  REFERENCES "users" ("id") ON DELETE CASCADE ON UPDATE CASCADE DEFERRABLE;

ALTER TABLE "subscription" ADD CONSTRAINT "subscription_ibfk_1" FOREIGN KEY ("category_id")
  REFERENCES "category" ("id") ON DELETE CASCADE ON UPDATE CASCADE DEFERRABLE;

ALTER TABLE "subscription" ADD CONSTRAINT "subscription_ibfk_2" FOREIGN KEY ("feed_id")
  REFERENCES "feed" ("id") ON DELETE CASCADE ON UPDATE CASCADE DEFERRABLE;

ALTER TABLE "subscription" ADD CONSTRAINT "subscription_ibfk_3" FOREIGN KEY ("user_id")
  REFERENCES "users" ("id") ON DELETE CASCADE ON UPDATE CASCADE DEFERRABLE;

ALTER TABLE "entry" ADD CONSTRAINT "entry_ibfk_1" FOREIGN KEY ("subscription_id")
  REFERENCES "subscription" ("id") ON DELETE CASCADE ON UPDATE CASCADE DEFERRABLE;

ALTER TABLE "entry" ADD CONSTRAINT "entry_ibfk_2" FOREIGN KEY ("user_id")
  REFERENCES "users" ("id") ON DELETE CASCADE ON UPDATE CASCADE DEFERRABLE;

