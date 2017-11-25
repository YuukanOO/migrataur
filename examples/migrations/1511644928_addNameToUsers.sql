-- Migrations 1511644928_addNameToUsers.sql
-- +migrataur up
alter table Users add Name varchar(50);
-- -migrataur up


-- +migrataur down
alter table Users drop column Name;
-- -migrataur down
