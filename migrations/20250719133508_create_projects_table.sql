-- +goose Up
-- +goose StatementBegin
create table projects(
    id bigint unsigned primary key auto_increment,
    name varchar(60) not null,
    due datetime null,
    created_at datetime not null default current_timestamp,
    updated_at datetime not null default current_timestamp on update current_timestamp,
    deleted_at datetime default null
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table projects;
-- +goose StatementEnd
