create table if not exists users(
    id serial primary key ,
    user_id varchar(255) unique not null,
    hash varchar(255));