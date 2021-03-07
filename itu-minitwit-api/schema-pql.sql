drop table if exists users;
create table users (
  user_id serial primary key,
  username varchar(256) not null,
  email varchar(256) not null,
  pw_hash varchar(256) not null
);

drop table if exists followers;
create table followers (
  who_id integer,
  whom_id integer
);

drop table if exists messages;
create table messages (
  message_id serial primary key,
  author_id integer not null,
  text varchar(256) not null,
  pub_date integer,
  flagged integer
);

drop table if exists latest;
create table latest (
  latest integer not null default 0
);