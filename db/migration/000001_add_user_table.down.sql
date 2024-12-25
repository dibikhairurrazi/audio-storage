BEGIN;

DROP INDEX IF EXISTS index_user_unique_email;
DROP TABLE IF EXISTS users;

COMMIT;
