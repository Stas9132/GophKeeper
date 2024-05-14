create table if not exists users(
                                    id serial primary key ,
                                    user_id varchar(255) unique not null,
    hash varchar(255));

create table if not exists meta(
    id serial primary key ,
    user_id varchar(255),
    obj_name varchar(255),
    obj_id varchar(255),
    obj_type numeric);