PRAGMA foreign_keys = on;

CREATE TABLE IF NOT EXISTS "draft" (
	"draft_id" INTEGER,
	"name" TEXT NOT NULL, 
	"content" TEXT,
	"created_at" DATE DEFAULT (datetime('now','localtime')),
	PRIMARY KEY("draft_id" AUTOINCREMENT)
);

-- create index for text search that groups by name
CREATE INDEX draft_index2 on "draft" (name, draft_id);

CREATE TABLE IF NOT EXISTS "comment" (
	"comment_id" INTEGER,
	"draft_id" INTEGER, -- FK
	"parent_id" INTEGER DEFAULT null, -- FK (self)
	"content" TEXT,
	"created_by" TEXT NOT NULL, -- FK
	"created_at" DATE DEFAULT (datetime('now','localtime')),
	PRIMARY KEY("comment_id" AUTOINCREMENT),
	FOREIGN KEY("draft_id")
		REFERENCES "draft" ("draft_id")
			ON DELETE CASCADE,
	FOREIGN KEY("created_by")
		REFERENCES "user" ("user_id")
			ON DELETE CASCADE,
	FOREIGN KEY("parent_id") -- parent of a comment
		REFERENCES "comment" ("comment_id")
			ON DELETE CASCADE
);

--create index for retrieving all comments by draft_id
CREATE INDEX comment_index2 on "comment" (draft_id);

CREATE TABLE IF NOT EXISTS "reaction" (
	"reaction_id" INTEGER,
	"comment_id" INTEGER, -- FK
	"reaction" TEXT,
	"created_by" TEXT NOT NULL, -- FK
	"created_at" DATE DEFAULT (datetime('now','localtime')),
	PRIMARY KEY("reaction_id" AUTOINCREMENT),
	FOREIGN KEY("comment_id")
		REFERENCES "comment" ("comment_id")
			ON DELETE CASCADE,
	FOREIGN KEY("created_by")
		REFERENCES "user" ("user_id")
			ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS "user" (
	"user_id" INTEGER,
	"name" TEXT,
	"role" TEXT,
	"token" TEXT, --would want this encrypted, or at least hashed
	PRIMARY KEY("user_id" AUTOINCREMENT)
);
