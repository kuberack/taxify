create table taxify.users (
  user_id int unsigned auto_increment,
  constraint pk_user_id primary key (user_id),
  phone_number varchar(64) not null unique,
  verify_sid varchar(255)
);
