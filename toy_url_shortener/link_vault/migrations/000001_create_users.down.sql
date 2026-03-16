-- DOWN: undo this migration (rollback)

DROP INDEX IF EXISTS idx_users_email;
DROP TABLE IF EXISTS users;
