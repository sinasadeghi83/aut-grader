-- +goose Up
-- +goose StatementBegin
create table users(
    id bigint unsigned primary key auto_increment,
    username varchar(25) not null,
    password  varchar(25) not null,
    phone     varchar(12) not null unique, 
    credits   int unsigned default 0,
    created_at datetime not null,
    updated_at datetime not null,
    deleted_at datetime default NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table users;
-- +goose StatementEnd

