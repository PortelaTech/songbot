create role ptuser password 'ptpass';
create database songs;
ALTER ROLE ptuser WITH LOGIN;
grant all on database songs to ptuser;
