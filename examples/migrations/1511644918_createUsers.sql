-- Migrations 1511644918_createUsers.sql
-- +migrataur up
create table Users (
  id int primary key,
  email varchar(250),
  registered_at datetime
);
-- -migrataur up


-- +migrataur down
drop table Users;
-- -migrataur down
