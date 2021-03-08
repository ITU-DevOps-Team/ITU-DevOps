ceate sequence user_id_seq;
drop table if exists users;
create table users (
  user_id primary key not null default nextval('user_id_seq'),
  username varchar(256) not null,
  email varchar(256) not null,
  pw_hash varchar(256) not null
);

drop table if exists followers;
create table followers (
  who_id integer,
  whom_id integer
);

create sequence message_id_seq;
drop table if exists messages;
create table messages (
  message_id primary key not null default nextval('message_id_seq'),
  author_id integer not null,
  text varchar(256) not null,
  pub_date integer,
  flagged integer
);

drop table if exists latests;
create table latests (
  latest integer not null default 0
);