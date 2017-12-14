create sequence loginuser_id_seq;
create table login_user (
  id int primary key default nextval('loginuser_id_seq'::regclass),
  username varchar(255) NOT NULL UNIQUE,
  password varchar(255) NOT NULL,
  name varchar(255) NOT NULL,
  roles varchar(255) NOT NULL,
  is_enabled boolean NOT NULL DEFAULT FALSE,
  created timestamp with time zone NOT NULL,
  last_login timestamp with time zone NOT NULL,
  is_sync boolean NOT NULL DEFAULT FALSE
 );

create sequence category_id_seq;
create table category (
  id int primary key default nextval('category_id_seq'::regclass),
  name varchar(255) NOT NULL UNIQUE,
  is_sync boolean NOT NULL DEFAULT FALSE
 );

create sequence country_id_seq;
create table country (
  id int primary key default nextval('country_id_seq'::regclass),
  name varchar(255) NOT NULL UNIQUE,
  is_sync boolean NOT NULL DEFAULT FALSE
 );

create sequence customer_id_seq;
create table customer (
  id int primary key default nextval('customer_id_seq'::regclass),
  name varchar(255) NOT NULL UNIQUE,
  address varchar(255),
  is_sync boolean NOT NULL DEFAULT FALSE
 );

create sequence supplier_id_seq;
create table supplier (
  id int primary key default nextval('supplier_id_seq'::regclass),
  name varchar(255) NOT NULL UNIQUE,
  address varchar(255),
  is_sync boolean NOT NULL DEFAULT FALSE
 );

create sequence unitofmeasurement_id_seq;
create table unitofmeasurement (
  id int primary key default nextval('unitofmeasurement_id_seq'::regclass),
  abbr varchar(255) NOT NULL UNIQUE,
  name varchar(255) NOT NULL,
  is_sync boolean NOT NULL DEFAULT FALSE
 );

 create sequence item_id_seq;
 create table item (
   id int primary key default nextval('item_id_seq'::regclass),
   code varchar(255) NOT NULL UNIQUE,
   name varchar(255) NOT NULL,
   type varchar(255) NOT NULL,
   brand_name varchar(255),
   packing_size varchar(255),
   case_pack float NOT NULL CHECK (case_pack > 0),
   caseuom int NOT NULL REFERENCES unitofmeasurement ON DELETE RESTRICT,
   pieceuom int NOT NULL REFERENCES unitofmeasurement ON DELETE RESTRICT,
   category int NOT NULL REFERENCES category ON DELETE RESTRICT,
   created timestamp with time zone NOT NULL,
   last_modified timestamp with time zone NOT NULL
 );