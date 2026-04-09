-- Добавляем поля для аутентификации и ролей
ALTER TABLE todoapp.users
ADD COLUMN email varchar(255),
ADD COLUMN password_hash varchar(255),
ADD COLUMN role integer,
ADD COLUMN manager_id integer;

-- Устанавливаем значения по умолчанию для существующих пользователей
UPDATE todoapp.users
SET 
    email = 'user' || id || '@example.com',
    password_hash = '',
    role = 1, -- UserRoleUser
    manager_id = NULL
WHERE email IS NULL;

-- Добавляем ограничения
ALTER TABLE todoapp.users
ALTER COLUMN email SET NOT NULL,
ALTER COLUMN password_hash SET NOT NULL,
ALTER COLUMN role SET NOT NULL,
ADD CONSTRAINT users_email_unique UNIQUE (email),
ADD CONSTRAINT users_email_check CHECK (email ~* '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$'),
ADD CONSTRAINT users_role_check CHECK (role IN (1, 2, 3)),
ADD CONSTRAINT users_manager_fk FOREIGN KEY (manager_id) REFERENCES todoapp.users(id) ON DELETE SET NULL;

-- Создаем индексы для производительности
CREATE INDEX idx_users_email ON todoapp.users(email);
CREATE INDEX idx_users_role ON todoapp.users(role);
CREATE INDEX idx_users_manager_id ON todoapp.users(manager_id);

COMMENT ON COLUMN todoapp.users.role IS '1=user, 2=manager, 3=admin';
COMMENT ON COLUMN todoapp.users.manager_id IS 'Reference to manager user (for hierarchy)';
