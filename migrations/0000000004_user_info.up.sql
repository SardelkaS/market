alter table "user" add column "name" varchar(256) not null default '';
alter table "user" add column phone_number varchar(32);
alter table "user" add column birth_date varchar(32);
alter table "user" add column contact_data varchar(1024);
alter table "user" add column email varchar(1024);