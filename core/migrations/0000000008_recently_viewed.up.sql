create table if not exists recently_viewed (
    id bigserial primary key,
    user_id bigint not null references "user"(id),
    product_id bigint not null references product(id),
    view_date timestamp not null default now()
);