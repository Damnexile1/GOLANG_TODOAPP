-- Добавляем колонки без NOT NULL constraint сначала
ALTER TABLE todoapp.tasks
ADD COLUMN status_key integer,
ADD COLUMN deadline timestamptz;

-- Устанавливаем статус на основе существующих данных
UPDATE todoapp.tasks
SET status_key = CASE
    WHEN completed = true THEN 2  -- completed
    WHEN completed = false THEN 1 -- created
END;

-- Теперь добавляем NOT NULL constraint и CHECK constraint
ALTER TABLE todoapp.tasks
ALTER COLUMN status_key SET NOT NULL,
ADD CONSTRAINT tasks_status_key_check CHECK (status_key IN (1, 2, 3));

COMMENT ON COLUMN todoapp.tasks.status_key IS '1=created, 2=completed, 3=failed';
