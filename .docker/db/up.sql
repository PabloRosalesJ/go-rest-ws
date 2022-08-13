drop table if exists users;

create table users (
  id varchar(255) primary key,
  email varchar(255) not null,
  password varchar(255) not null,
  created_at timestamp not null default now()
);