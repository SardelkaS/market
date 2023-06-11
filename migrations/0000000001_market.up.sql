create table if not exists public."user"(
    id bigserial primary key,
    "login" text unique not null,
    "password" text not null,
    "user_role" text not null default 'user',
    internal_id text not null,
    ban bool not null default false,
    "timezone" text not null
);

create table if not exists public."category" (
    id bigserial primary key,
    "name" text not null
);

create table if not exists public."manufacturer" (
    id bigserial primary key,
    "name" text not null
);

create table if not exists public."product" (
    id bigserial primary key,
    internal_id text not null unique,
    "name" text not null,
    price float8 not null check (price > 0),
    count bigint not null default 0 check (count >= 0),
    manufacturer_id bigint not null references manufacturer(id),
    buy_count bigint not null default 0 check (buy_count >= 0),
    description text,
    pictures text[],
    "show" bool not null default false
);

create table if not exists public."feedback" (
    id bigserial primary key,
    user_id bigint not null references "user"(id),
    product_id bigint not null references product(id),
    stars bigint not null default 5 check (stars > 0 and stars < 6),
    "message" text,
    pictures text[]
);

create table if not exists public."like_product" (
    id bigserial primary key,
    product_id bigint not null references product(id),
    user_id bigint not null references "user"(id)
);

create table if not exists public."product_category" (
    id bigserial primary key,
    product_id bigint not null references product(id),
    category_id bigint not null references category(id)
);

create table if not exists public."basket" (
    id bigserial primary key,
    user_id bigint not null unique references "user"(id),
    product_id bigint not null references product(id),
    count bigint not null default 1 check (count >= 0)
);

create table if not exists public."order_status" (
    id bigserial primary key,
    "name" text not null unique
);

create table if not exists public."order" (
    id bigserial primary key,
    internal_id text unique not null,
    user_id bigint not null references "user"(id),
    status_id bigint not null references order_status(id),
    address text,
    contact_data text,
    create_time timestamp not null,
    update_time timestamp,
    complete_time timestamp
);

create table if not exists public."order_products" (
    id bigserial primary key,
    order_id bigint not null references "order"(id),
    product_id bigint not null references product(id),
    count bigint not null default 1 check (count > 0)
);