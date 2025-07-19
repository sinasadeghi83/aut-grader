-- +goose Up
-- +goose StatementBegin
create table theaders(
    id bigint unsigned primary key auto_increment,
    hkey varchar(300) not null,
    hvalue varchar(300) not null,
    test_id bigint unsigned not null,

    created_at datetime not null default current_timestamp,
    updated_at datetime not null default current_timestamp on update current_timestamp,
    deleted_at datetime default null,

    foreign key (test_id) references tests(id) on delete cascade,
    
    index idx_theaders_test_key (test_id, hkey)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table theaders;
-- +goose StatementEnd
