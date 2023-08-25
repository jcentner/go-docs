-- test data
PRAGMA foreign_keys = on;

INSERT INTO "draft" (name, content) VALUES ("doc 0", "some text for draft 0");
INSERT INTO "draft" (name, content) VALUES ("doc 0", "some updated text for draft 1 on document 0");
INSERT INTO "draft" (name, content) VALUES ("doc 1", "some text for document 1");

INSERT INTO "user" (name, role, token) VALUES ("admin user", "admin", "thisisatesttoken");
INSERT INTO "user" (name, role, token) VALUES ("regular user", "user", "regularusertoken");

INSERT INTO "comment" (draft_id, content, created_by) VALUES (1, "a comment on the first draft", 1);
INSERT INTO "comment" (draft_id, parent_id, content, created_by) VALUES (1, 1, "a comment on the comment", 1);
