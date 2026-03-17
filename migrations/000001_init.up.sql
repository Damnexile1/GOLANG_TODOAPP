CREATE SCHEMA IF NOT EXISTS todoapp;

create table todoapp.users (
    id           serial                primary key,
    version      bigint       not null default 1,
    full_name    varchar(100) not null check (char_length(full_name) between 3 and 100),
    phone_number varchar(15)           check (
        phone_number ~ '^\+[0-9]+$'
        and
        char_length(phone_number) between 10 and 15
    )
);

create table todoapp.tasks (
    id              serial                      primary key,
    version         bigint          not null    default 1,
    title           varchar(100)    not null    check ( char_length(title) between 1 and 100),
    description     varchar(1000),              check ( char_length(description) between 1 and 100),
    completed       boolean         not null,
    created_at      timestamptz     not null,
    completed_at    timestamptz,

    check (
        (completed=false and completed_at is null)
        or
        (completed=true and completed_at is not null and completed_at >= created_at)
    ),

    author_user_id  integer         not null    references todoapp.users(id)
);
