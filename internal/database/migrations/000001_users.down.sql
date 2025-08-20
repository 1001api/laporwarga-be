-- Drop trigger
DROP TRIGGER IF EXISTS set_users_timestamp ON users;

-- Drop trigger function
DROP FUNCTION IF EXISTS trigger_set_timestamp();

-- Drop index
DROP INDEX IF EXISTS idx_users_password_reset_token;

-- Drop users table
DROP TABLE IF EXISTS users;
