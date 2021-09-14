create database casher;
\c casher

drop table operations;
drop table users;

create table users (
    id serial primary key,
    login varchar(256) unique not null,
    password varchar(256) not null,
    name varchar(256) not null,
    birth timestamp with time zone not null,
    created_at timestamp with time zone default now() not null
);
create index if not exists login_queue_idx on users (login);

create table operations (
    id serial primary key,
    user_id bigint references users (id) not null,
    subject varchar(256) not null,
    amount bigint not null,
    type int not null,
    message text,
    created_at timestamp with time zone default now() not null
);