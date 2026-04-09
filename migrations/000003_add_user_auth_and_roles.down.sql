-- Удаляем индексы
DROP INDEX IF EXISTS todoapp.idx_users_manager_id;
DROP INDEX IF EXISTS todoapp.idx_users_role;
DROP INDEX IF EXISTS todoapp.idx_users_email;

-- Удаляем ограничения
ALTER TABLE todoapp.users
DROP CONSTRAINT IF EXISTS users_manager_fk,
DROP CONSTRAINT IF EXISTS users_role_check,
DROP CONSTRAINT IF EXISTS users_email_check,
DROP CONSTRAINT IF EXISTS users_email_unique;

-- Удаляем колонки
ALTER TABLE todoapp.users
DROP COLUMN IF EXISTS manager_id,
DROP COLUMN IF EXISTS role,
DROP COLUMN IF EXISTS password_hash,
DROP COLUMN IF EXISTS email;
