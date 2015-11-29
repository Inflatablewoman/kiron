drop type gender cascade;
create type gender as enum ('male', 'female');

drop type status cascade;
create type status as enum (
  'received',
  'confirmed',
  'in verification',
  'verified',
  'waiting for response',
  'rejected',
  'accepted'
);

drop table if exists education_levels cascade;
create table education_levels (
  id serial primary key,
  education_level text unique not null
);

insert into education_levels (education_level) values
('none'),
('elementary'),
('secondary'),
('associate'),
('bachelor');

drop table if exists roles cascade;
create table roles (
  id serial primary key,
  role text unique not null
);

insert into roles (role) values
('admin'),
('sub-admin'),
('trusted helper'),
('limited helper'),
('applicant');

drop table if exists users cascade;
create table users (
  id serial primary key,
  email text not null unique,
  name text not null,
  lastname text not null,
  password text not null,
  -- make sure that the password is hashed with bcrypt. see also
  -- https://en.wikipedia.org/wiki/Bcrypt
  constraint password_in_bcrypt check (
    password like '$2a$%' or
    password like '$2b$%'
  ),
  created_at timestamp not null,
  role_id integer references roles not null
);

drop table if exists auth_tokens cascade;
create table auth_tokens (
  user_id integer references users not null,
  token text not null,
  expires timestamp not null
);

drop table if exists applications cascade;
create table applications (
  id serial primary key,
  birthday date not null,
  phone text,
  nationality text not null,
  country text not null,  -- address
  city text not null,     -- address
  zip text not null,      -- address
  address_extra text,
  first_page_of_survey_data text,
  gender gender not null,
  study_program text,
  user_id integer references users not null unique,
  education_level_id integer references education_levels not null,
  status status not null default 'received',
  blocked_until timestamp,
  created_at timestamp not null,
  edited_at timestamp not null
);

drop table if exists document_types cascade;
create table document_types (
  id serial primary key,
  document_type text unique not null
);

insert into document_types (document_type) values
('1refugee status'),
('unhcr refugee status'),
('asylum application'),
('refugee camp'),
('aufenthaltserlaubnis'),
('aufenthaltsgestattung'),
('duldung'),
('subsidiary protection status'),
('our-certification');

drop table if exists documents cascade;
create table documents (
  id serial primary key,
  application_id integer references applications not null,
  document_type_id integer references document_types not null,
  contents bytea not null -- document itself
);

drop table if exists comments cascade;
create table comments (
  id serial primary key,
  created_at timestamp not null,
  application_id integer references applications not null,
  user_id integer references users not null,
  contents text not null
);

-- some sample records to work with

begin;

insert into users (name, lastname, email, password, created_at, role_id) values (
  'foo', 'bar', 'foo@example.org',
  '$2a$10$FTHN0Dechb/IiQuyeEwxaOCSdBss1KcC5fBKDKsj85adOYTLOPQf6', NOW(),
  (select id from roles where role = 'applicant')
);

insert into applications
(
  birthday, phone, nationality, country,
  city, zip, address_extra, first_page_of_survey_data, gender, education_level_id,
  user_id, status, created_at, edited_at
) values (
  '2000-01-01', '123456789', 'german', 'germany', 'munich', '80331',
  'po box 123', 'first page of the survey data', 'male',
  (select id from education_levels where education_level = 'elementary'),
  (select id from users where email = 'foo@example.org'),
  'received', now(), now()
);

insert into documents (application_id, document_type_id, contents) values (
  (
    select app.id from applications app
    join users on users.id = app.user_id
    where users.email = 'foo@example.org'
  ),
  (select id from document_types where document_type = '1refugee status'),
  '[contents of a pdf file]'
);

commit;
