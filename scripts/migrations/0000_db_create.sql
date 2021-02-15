drop database workery_v2_db;
create database workery_v2_db;
\c workery_v2_db;
CREATE USER golang WITH PASSWORD '123password';
GRANT ALL PRIVILEGES ON DATABASE workery_v2_db to golang;
ALTER USER golang CREATEDB;
ALTER ROLE golang SUPERUSER;
