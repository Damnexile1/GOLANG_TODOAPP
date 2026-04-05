ALTER TABLE todoapp.tasks
DROP COLUMN IF EXISTS status_key,
DROP COLUMN IF EXISTS deadline;
