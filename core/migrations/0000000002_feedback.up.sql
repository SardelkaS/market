alter table public."feedback" add column if not exists internal_id text not null default '';
alter table public."feedback" add column if not exists create_date timestamp not null default now();
alter table public."feedback" add column if not exists update_date timestamp;
alter table public."feedback" add column if not exists is_removed bool not null default false;

create table if not exists feedback_like(
    id bigserial primary key,
    feedback_id bigint not null references feedback(id),
    user_id bigint not null
);