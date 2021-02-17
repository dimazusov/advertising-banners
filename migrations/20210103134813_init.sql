-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE banner (
    id serial PRIMARY KEY,
    description VARCHAR (100) not null
);
CREATE TABLE banner_place (
    banner_id int NOT NULL,
    place_id int NOT NULL
);
CREATE TABLE event (
   id serial PRIMARY KEY,
   type varchar NOT NULL,
   place_id int NOT NULL,
   banner_id int NOT NULL,
   soc_group_id int NOT NULL,
   time bigint NOT NULL
);
CREATE TABLE place (
   id serial PRIMARY KEY,
   description VARCHAR (100) not null
);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
drop table banner;
drop table banner_place;
drop table event;
drop table place;