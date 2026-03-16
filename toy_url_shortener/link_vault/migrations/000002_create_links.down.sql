-- DOWN: undo this migration

DROP INDEX IF EXISTS idx_links_user_id;
DROP INDEX IF EXISTS idx_links_short_code;
DROP TABLE IF EXISTS links;
