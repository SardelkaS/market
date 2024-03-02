create table if not exists public."user"(
    id bigserial primary key,
    "login" text unique not null,
    "password" text not null,
    "user_role" text not null default 'user',
    internal_id text not null,
    ban bool not null default false,
    "timezone" text not null
);