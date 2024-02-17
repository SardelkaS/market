create table if not exists subcategory (
    id bigserial primary key,
    "name" varchar not null,
    category_id bigint not null references category(id)
);

drop table if exists product_category;

alter table product add column if not exists subcategory_id bigint references subcategory(id);