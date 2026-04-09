# Техническое задание: Система оценки навыков и управления задачами

## 1. Общее описание проекта

### 1.1 Цель проекта
Создание системы управления задачами с встроенной системой оценки навыков сотрудников/пользователей на основе их активности и результатов выполнения задач. Система должна предоставлять визуализацию прогресса в виде радар-чарта (как в Dota 2) по различным категориям навыков.

### 1.2 Ключевые особенности
- Многоуровневая система ролей (пользователь, руководитель, администратор)
- Планирование задач на различные периоды (день, неделя, месяц)
- Привязка задач к категориям навыков
- Система оценки и начисления баллов
- Обязательное подтверждение выполнения задач (документы, фото)
- Визуализация прогресса и статистики
- Панель управления для руководителей

### 1.3 Целевая аудитория
- Команды разработчиков
- Образовательные учреждения
- Корпоративные отделы
- Фрилансеры и самозанятые

---

## 2. Функциональные требования

### 2.1 Система аутентификации и авторизации

#### 2.1.1 Регистрация и вход
- Регистрация по email + password
- Вход по email + password
- JWT токены для авторизации
- Refresh token механизм
- Восстановление пароля (опционально для MVP)

#### 2.1.2 Роли пользователей
1. **User (Пользователь)**
   - Может создавать личные задачи
   - Может выполнять назначенные задачи
   - Видит свою статистику и прогресс
   - Доступ только к своим данным

2. **Manager (Руководитель)**
   - Все права User
   - Может назначать задачи подчиненным
   - Видит статистику своей команды
   - Может создавать планы на период
   - Может утверждать/отклонять выполненные задачи

3. **Admin (Администратор)**
   - Все права Manager
   - Управление пользователями
   - Управление категориями навыков
   - Настройка системы оценки
   - Доступ ко всей статистике

#### 2.1.3 Иерархия пользователей
- User может быть привязан к Manager
- Manager может иметь несколько подчиненных Users
- Admin видит всех пользователей системы

---

### 2.2 Категории навыков (Skills)

#### 2.2.1 Структура категории
```
Skill {
  id: int
  name: string (например: "Backend Development", "Communication", "Problem Solving")
  description: string
  icon: string (опционально)
  max_points: int (максимальное количество баллов, например 100)
  created_at: timestamp
}
```

#### 2.2.2 Примеры категорий
- **Technical Skills**: Backend, Frontend, DevOps, Testing, Database
- **Soft Skills**: Communication, Leadership, Time Management, Teamwork
- **Domain Skills**: Business Analysis, Product Management, Design
- **Learning**: Self-education, Mentoring, Documentation

#### 2.2.3 Управление категориями
- Admin может создавать/редактировать/удалять категории
- Категории глобальные для всей системы
- Каждая задача может быть привязана к 1-3 категориям

---

### 2.3 Планы на период (Period Plans)

#### 2.3.1 Структура плана
```
PeriodPlan {
  id: int
  user_id: int (кому назначен план)
  created_by_user_id: int (кто создал план)
  title: string
  description: string
  period_type: enum (daily, weekly, monthly, custom)
  start_date: date
  end_date: date
  status: enum (draft, active, completed, cancelled)
  created_at: timestamp
  updated_at: timestamp
}
```

#### 2.3.2 Типы периодов
- **Daily**: план на день
- **Weekly**: план на неделю (понедельник-воскресенье)
- **Monthly**: план на месяц (1-е число - последнее число)
- **Custom**: произвольный период

#### 2.3.3 Логика работы с планами
- План создается Manager или самим User
- В начале периода план переходит в статус "active"
- К плану привязываются задачи
- По окончании периода план переходит в "completed"
- Можно отменить план (cancelled)

---

### 2.4 Задачи (Tasks) - расширенная версия

#### 2.4.1 Обновленная структура Task
```
Task {
  // Существующие поля
  id: int
  version: int
  title: string
  description: string
  completed: bool
  status_key: int (1=created, 2=completed, 3=failed)
  deadline: timestamp
  created_at: timestamp
  completed_at: timestamp
  author_user_id: int
  
  // Новые поля
  assigned_to_user_id: int (кому назначена задача)
  period_plan_id: int (FK к PeriodPlan, nullable)
  priority: enum (low, medium, high, critical)
  estimated_hours: float (оценка времени на выполнение)
  actual_hours: float (фактическое время)
  points: int (баллы за выполнение)
  requires_proof: bool (требуется подтверждение)
  approval_status: enum (pending, approved, rejected, not_required)
  approved_by_user_id: int (кто утвердил)
  approved_at: timestamp
  rejection_reason: string
}
```

#### 2.4.2 Связь задач с навыками
```
TaskSkill {
  id: int
  task_id: int (FK)
  skill_id: int (FK)
  weight: float (вес навыка для этой задачи, 0.0-1.0)
  points_earned: int (заработанные баллы по этому навыку)
}
```

Логика:
- Задача может быть привязана к 1-3 навыкам
- Сумма весов всех навыков = 1.0
- При выполнении задачи баллы распределяются пропорционально весам

Пример:
```
Task: "Разработать REST API для модуля пользователей" (100 баллов)
- Backend Development (weight: 0.6) → 60 баллов
- Database Design (weight: 0.3) → 30 баллов
- Documentation (weight: 0.1) → 10 баллов
```

#### 2.4.3 Жизненный цикл задачи
1. **Created** - задача создана
2. **In Progress** - пользователь начал работу (опционально)
3. **Pending Approval** - задача выполнена, ожидает проверки (если requires_proof=true)
4. **Approved** - задача утверждена, баллы начислены
5. **Rejected** - задача отклонена, нужно переделать
6. **Completed** - задача завершена (если не требует проверки)
7. **Failed** - задача провалена (deadline истек)

---

### 2.5 Подтверждения выполнения (Task Proofs)

#### 2.5.1 Структура подтверждения
```
TaskProof {
  id: int
  task_id: int (FK)
  user_id: int (кто приложил)
  proof_type: enum (document, image, link, text)
  file_path: string (путь к файлу, если document/image)
  file_name: string
  file_size: int
  url: string (если link)
  text_content: string (если text)
  comment: string (комментарий пользователя)
  created_at: timestamp
}
```

#### 2.5.2 Типы подтверждений
- **Document**: PDF, DOCX, TXT и т.д.
- **Image**: JPG, PNG, GIF
- **Link**: ссылка на внешний ресурс (GitHub, Figma, Google Docs)
- **Text**: текстовое описание результата

#### 2.5.3 Хранение файлов
- Локальное хранилище: `/uploads/proofs/{user_id}/{task_id}/{filename}`
- Ограничение размера: 10MB на файл
- Максимум 5 файлов на задачу
- Опционально: интеграция с S3/MinIO для продакшена

---

### 2.6 Система оценки и баллов

#### 2.6.1 Начисление баллов
```
UserSkillScore {
  id: int
  user_id: int (FK)
  skill_id: int (FK)
  total_points: int (общее количество баллов)
  level: int (уровень навыка, рассчитывается автоматически)
  last_updated: timestamp
}
```

#### 2.6.2 Расчет уровня навыка
```
Level 1: 0-100 баллов
Level 2: 101-250 баллов
Level 3: 251-500 баллов
Level 4: 501-1000 баллов
Level 5: 1001+ баллов
```

#### 2.6.3 Логика начисления
1. Задача выполнена и утверждена (или не требует утверждения)
2. Баллы распределяются по навыкам согласно весам
3. Обновляется UserSkillScore для каждого навыка
4. Пересчитывается уровень навыка
5. Создается запись в истории (UserSkillHistory)

#### 2.6.4 История изменений баллов
```
UserSkillHistory {
  id: int
  user_id: int (FK)
  skill_id: int (FK)
  task_id: int (FK, nullable)
  points_change: int (может быть отрицательным)
  reason: string
  created_at: timestamp
}
```

---

### 2.7 Статистика и аналитика

#### 2.7.1 Личная статистика пользователя
- Радар-чарт по всем навыкам (как в Dota 2)
- Прогресс по текущему плану (% выполнения)
- Количество выполненных задач за период
- Общее количество баллов
- Топ-3 навыка
- График активности (задачи по дням)
- Среднее время выполнения задач

#### 2.7.2 Статистика для Manager
- Сводка по команде (таблица пользователей)
- Средние баллы команды по навыкам
- Топ исполнителей
- Задачи, ожидающие утверждения
- Просроченные задачи
- Сравнение пользователей (радар-чарты рядом)

#### 2.7.3 Статистика для Admin
- Общая статистика системы
- Количество пользователей, задач, планов
- Самые популярные навыки
- Средняя активность пользователей
- Графики роста системы

---

### 2.8 API Endpoints

#### 2.8.1 Authentication
```
POST   /api/v1/auth/register
POST   /api/v1/auth/login
POST   /api/v1/auth/refresh
POST   /api/v1/auth/logout
POST   /api/v1/auth/forgot-password (опционально)
```

#### 2.8.2 Users
```
GET    /api/v1/users/me
PATCH  /api/v1/users/me
GET    /api/v1/users/{id}
GET    /api/v1/users (Admin/Manager)
PATCH  /api/v1/users/{id}/role (Admin)
DELETE /api/v1/users/{id} (Admin)
GET    /api/v1/users/{id}/subordinates (Manager)
```

#### 2.8.3 Skills
```
GET    /api/v1/skills
GET    /api/v1/skills/{id}
POST   /api/v1/skills (Admin)
PATCH  /api/v1/skills/{id} (Admin)
DELETE /api/v1/skills/{id} (Admin)
```

#### 2.8.4 Period Plans
```
GET    /api/v1/plans
GET    /api/v1/plans/{id}
POST   /api/v1/plans
PATCH  /api/v1/plans/{id}
DELETE /api/v1/plans/{id}
GET    /api/v1/plans/{id}/tasks
GET    /api/v1/plans/active (текущий активный план)
```

#### 2.8.5 Tasks (расширенные)
```
GET    /api/v1/tasks
GET    /api/v1/tasks/{id}
POST   /api/v1/tasks
PATCH  /api/v1/tasks/{id}
DELETE /api/v1/tasks/{id}
POST   /api/v1/tasks/{id}/start (начать работу)
POST   /api/v1/tasks/{id}/complete (завершить)
POST   /api/v1/tasks/{id}/approve (Manager)
POST   /api/v1/tasks/{id}/reject (Manager)
GET    /api/v1/tasks/pending-approval (Manager)
```

#### 2.8.6 Task Proofs
```
POST   /api/v1/tasks/{id}/proofs (upload)
GET    /api/v1/tasks/{id}/proofs
GET    /api/v1/proofs/{id}
DELETE /api/v1/proofs/{id}
GET    /api/v1/proofs/{id}/download
```

#### 2.8.7 Statistics
```
GET    /api/v1/stats/me
GET    /api/v1/stats/users/{id}
GET    /api/v1/stats/team (Manager)
GET    /api/v1/stats/system (Admin)
GET    /api/v1/stats/skills-radar
GET    /api/v1/stats/activity-chart
```

---

## 3. Нефункциональные требования

### 3.1 Производительность
- Время ответа API: < 200ms для 95% запросов
- Поддержка до 1000 одновременных пользователей
- Размер базы данных: до 100GB

### 3.2 Безопасность
- Все пароли хешируются (bcrypt)
- JWT токены с коротким временем жизни (15 мин access, 7 дней refresh)
- HTTPS обязателен в продакшене
- Rate limiting на API endpoints
- Валидация всех входных данных
- SQL injection защита (prepared statements)
- XSS защита
- CORS настройка

### 3.3 Масштабируемость
- Stateless API (можно горизонтально масштабировать)
- Кеширование частых запросов (Redis, опционально)
- Пагинация для всех списков
- Индексы в БД для частых запросов

### 3.4 Надежность
- Graceful shutdown
- Логирование всех ошибок
- Backup БД (ежедневно)
- Транзакции для критичных операций

---

## 4. Архитектура системы

### 4.1 Слои приложения
```
┌─────────────────────────────────────┐
│         Transport Layer             │
│  (HTTP Handlers, Middleware, DTOs)  │
└─────────────────────────────────────┘
              ↓
┌─────────────────────────────────────┐
│         Service Layer               │
│  (Business Logic, Validation)       │
└─────────────────────────────────────┘
              ↓
┌─────────────────────────────────────┐
│         Repository Layer            │
│  (Database Queries, Models)         │
└─────────────────────────────────────┘
              ↓
┌─────────────────────────────────────┐
│         Domain Layer                │
│  (Entities, Value Objects)          │
└─────────────────────────────────────┘
```

### 4.2 Структура проекта
```
GOLANG_TODOAPP/
├── cmd/
│   └── todoapp/
│       └── main.go
├── internal/
│   ├── core/
│   │   ├── config/
│   │   ├── domain/
│   │   │   ├── user.go
│   │   │   ├── task.go
│   │   │   ├── skill.go
│   │   │   ├── period_plan.go
│   │   │   ├── task_proof.go
│   │   │   └── user_skill_score.go
│   │   ├── errors/
│   │   ├── logger/
│   │   ├── repository/
│   │   └── transport/
│   │       └── http/
│   │           ├── middleware/
│   │           │   ├── auth.go
│   │           │   ├── role.go
│   │           │   └── rate_limit.go
│   │           └── server/
│   └── features/
│       ├── auth/
│       │   ├── repository/
│       │   ├── service/
│       │   └── transport/
│       ├── users/
│       ├── tasks/
│       ├── skills/
│       ├── plans/
│       ├── proofs/
│       └── statistics/
├── migrations/
├── uploads/
│   └── proofs/
├── documentation/
│   └── technical_specification.md
├── docker-compose.yaml
└── Makefile
```

### 4.3 База данных

#### Основные таблицы:
1. **users** - пользователи
2. **tasks** - задачи (расширенная)
3. **skills** - категории навыков
4. **period_plans** - планы на период
5. **task_skills** - связь задач и навыков
6. **task_proofs** - подтверждения выполнения
7. **user_skill_scores** - баллы пользователей по навыкам
8. **user_skill_history** - история изменений баллов

#### Связи:
- users 1:N tasks (author)
- users 1:N tasks (assigned_to)
- users 1:N period_plans
- period_plans 1:N tasks
- tasks N:M skills (через task_skills)
- tasks 1:N task_proofs
- users N:M skills (через user_skill_scores)

---

## 5. План разработки (Roadmap)

### Phase 1: Фундамент (2-3 недели)
- [ ] Система аутентификации (JWT)
- [ ] Роли и права доступа
- [ ] Middleware для авторизации
- [ ] Обновление структуры User
- [ ] Миграции БД

### Phase 2: Навыки и планы (2 недели)
- [ ] CRUD для Skills
- [ ] CRUD для Period Plans
- [ ] Связь планов с задачами
- [ ] Обновление Tasks под новые требования

### Phase 3: Система оценки (2 недели)
- [ ] Связь задач с навыками (task_skills)
- [ ] Логика начисления баллов
- [ ] UserSkillScore и история
- [ ] Расчет уровней

### Phase 4: Подтверждения (1-2 недели)
- [ ] Upload файлов
- [ ] CRUD для Task Proofs
- [ ] Система утверждения задач
- [ ] Хранилище файлов

### Phase 5: Статистика (2 недели)
- [ ] Личная статистика
- [ ] Статистика команды
- [ ] Радар-чарт данные
- [ ] Графики активности

### Phase 6: Полировка (1 неделя)
- [ ] Оптимизация запросов
- [ ] Индексы БД
- [ ] Документация API (Swagger)
- [ ] Тестирование

---

## 6. Технологический стек

### Backend
- **Язык**: Go 1.22+
- **Web Framework**: net/http (stdlib)
- **Database**: PostgreSQL 15+
- **Migrations**: golang-migrate
- **Validation**: go-playground/validator
- **Logging**: zap
- **Config**: envconfig
- **JWT**: golang-jwt/jwt
- **File Upload**: multipart/form-data
- **Documentation**: Swagger (swaggo)

### Infrastructure
- **Containerization**: Docker
- **Orchestration**: Docker Compose
- **Database**: PostgreSQL в контейнере
- **Storage**: Локальная FS (или MinIO для продакшена)

### Development Tools
- **IDE**: GoLand / VS Code
- **API Testing**: Postman / Insomnia
- **DB Client**: DBeaver / pgAdmin
- **Version Control**: Git

---

## 7. Примеры использования

### 7.1 Сценарий 1: Создание плана на месяц
1. Manager логинится в систему
2. Создает Period Plan на апрель 2026 для User
3. Добавляет 20 задач в план
4. Каждую задачу привязывает к навыкам с весами
5. Устанавливает requires_proof=true для важных задач
6. User видит план и начинает выполнять задачи

### 7.2 Сценарий 2: Выполнение задачи с подтверждением
1. User выбирает задачу из плана
2. Нажимает "Start" (опционально)
3. Выполняет задачу
4. Загружает скриншот результата + ссылку на GitHub
5. Нажимает "Complete"
6. Manager получает уведомление
7. Manager проверяет результат и нажимает "Approve"
8. Баллы начисляются User по навыкам
9. Статистика обновляется

### 7.3 Сценарий 3: Просмотр статистики
1. User открывает страницу статистики
2. Видит радар-чарт по 5 навыкам
3. Видит прогресс по текущему плану (15/20 задач)
4. Видит график активности за месяц
5. Видит топ-3 навыка: Backend (Level 3), Database (Level 2), Testing (Level 2)

---

## 8. Риски и ограничения

### 8.1 Технические риски
- **Сложность расчета баллов**: нужна четкая формула
- **Производительность статистики**: много агрегаций
- **Хранение файлов**: может быстро расти

### 8.2 Бизнес-риски
- **Геймификация может демотивировать**: если баллы несправедливы
- **Субъективность оценки**: Manager может быть необъективным
- **Сложность для новых пользователей**: много сущностей

### 8.3 Ограничения MVP
- Нет real-time уведомлений (WebSocket)
- Нет мобильного приложения
- Нет интеграций (Slack, Telegram)
- Нет экспорта отчетов (PDF, Excel)
- Нет системы достижений (badges)

---

## 9. Метрики успеха

### 9.1 Технические метрики
- Uptime > 99%
- API response time < 200ms
- Zero critical bugs в продакшене
- Test coverage > 70%

### 9.2 Бизнес-метрики
- Количество активных пользователей
- Среднее количество задач на пользователя
- Процент выполненных планов
- Средний рейтинг системы от пользователей

---

## 10. Дальнейшее развитие (Post-MVP)

### 10.1 Фичи для будущих версий
- [ ] Real-time уведомления (WebSocket)
- [ ] Система достижений и бейджей
- [ ] Таблица лидеров (leaderboard)
- [ ] Экспорт отчетов (PDF, Excel)
- [ ] Интеграции (Slack, Telegram, Jira)
- [ ] Мобильное приложение
- [ ] Календарь задач
- [ ] Recurring tasks (повторяющиеся задачи)
- [ ] Комментарии к задачам
- [ ] Mentions (@username)
- [ ] Темная тема UI

### 10.2 Оптимизации
- [ ] Кеширование (Redis)
- [ ] Full-text search (Elasticsearch)
- [ ] CDN для статики
- [ ] S3 для файлов
- [ ] Микросервисная архитектура

---

## 11. Заключение

Данное ТЗ описывает амбициозную систему управления задачами с геймификацией и оценкой навыков. Проект сохраняет существующий функционал TODO-приложения и значительно расширяет его возможности.

Ключевые преимущества системы:
- Мотивация через геймификацию
- Прозрачная оценка навыков
- Контроль выполнения задач
- Визуализация прогресса
- Гибкая система ролей

Рекомендуется разрабатывать систему итеративно, начиная с MVP (Phase 1-3), затем добавляя дополнительные фичи на основе обратной связи пользователей.

---

**Версия документа**: 1.0  
**Дата создания**: 2026-04-07  
**Автор**: OpenCode AI Assistant  
**Статус**: Draft
