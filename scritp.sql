-- drop database if exists db_marketplace;

CREATE DATABASE db_marketplace CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
use db_marketplace;

create table clients(
    id int unsigned auto_increment primary key,
    name varchar(255) not null,
    email varchar(255) not null,
    phone varchar(100) not null,
    `created_at` datetime default now()
);

insert into clients (name, email, phone) values ('Higor Diego', 'higordiegoti@gmail.com', '88997613741');

create table destinations (
    id int unsigned auto_increment primary key,
    street varchar(255) not null,
    district varchar(255) not null,
    city varchar(255) not null,
    country varchar(255) not null,
    state varchar(255) not null,
    number varchar(255) not null,
    lat varchar(255),
    lng varchar(255),
    zipCode varchar(100) not null,
    id_client int unsigned not null,
	FOREIGN KEY (id_client) REFERENCES clients(id),
    `created_at` datetime default now()
);

insert into destinations (street, district, city, country, state, number, lat, lng, zipCode, id_client) values ('Padre josé alves', 'Salesianos', 'Juazeiro do Norte', 'Brasil', 'Ceará', '790', '-7.205440', '-39.324280', '63050222', 1);

create table products (
    id int unsigned auto_increment primary key,
    name varchar(255) not null,
    price decimal(10, 2) not null,
    nfe text not null,
    id_client int unsigned not null,
	FOREIGN KEY (id_client) REFERENCES clients(id),
    `created_at` datetime default now()
);

insert into products (`name`, price, nfe, id_client) values ('Iphone 12', 4900.00, '20320131203213021301321', 1);
