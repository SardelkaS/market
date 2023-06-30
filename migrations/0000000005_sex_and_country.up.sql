create table if not exists sex (
    id bigserial primary key,
    "name" varchar(32) not null
);

create table if not exists country (
    id bigserial primary key,
    "name" varchar not null
);

alter table product add column if not exists sex_id bigint references sex(id);
alter table product add column if not exists country_id bigint references country(id);