drop table if exists users;

create table users (
  id varchar(32) primary key,
  email varchar(255) not null,
  password varchar(255) not null,
  created_at timestamp not null default now()
);

drop table if exists posts;

create table posts (
  id varchar(32) primary key,
  user_id varchar(32) not null,
  post_content varchar(120) not null,
  created_at timestamp not null default now(),
  foreign key (user_id) references users (id)
);