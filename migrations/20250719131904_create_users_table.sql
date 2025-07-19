-- +goose Up
-- +goose StatementBegin
create table users(
    id bigint unsigned primary key auto_increment,
    username varchar(25) not null,
    password  varchar(25) not null,
    created_at datetime not null,
    updated_at datetime not null,
    deleted_at datetime default NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table users;
-- +goose StatementEnd

