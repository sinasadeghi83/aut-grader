-- +goose Up
-- +goose StatementBegin
create table tests(
    id bigint unsigned primary key auto_increment,
    name varchar(60) not null,
    scenario_id bigint unsigned not null,
    url varchar(2048) not null,
    method varchar(10) not null,
    req_body TEXT null,
    status_code int unsigned not null,
    res_body TEXT null,
    depends_on_id bigint unsigned null,

    created_at datetime not null default current_timestamp,
    updated_at datetime not null default current_timestamp on update current_timestamp,
    deleted_at datetime default null,

    foreign key (scenario_id) references scenarios(id) on delete cascade,
    foreign key (depends_on_id) references tests(id) on delete cascade,
    
    index idx_tests_depends_on (depends_on_id),
    index idx_tests_scenario_name (scenario_id, name)
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TRIGGER before_test_insert
BEFORE INSERT ON tests
FOR EACH ROW
BEGIN
    IF NEW.depends_on_id IS NOT NULL THEN
        SET @parent_scenario_id = (SELECT scenario_id FROM tests WHERE id = NEW.depends_on_id);
        IF @parent_scenario_id != NEW.scenario_id THEN
            SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'depends_on test must be in the same scenario';
        END IF;
    END IF;
END;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TRIGGER before_test_update
BEFORE UPDATE ON tests
FOR EACH ROW
BEGIN
    IF NEW.depends_on_id IS NOT NULL THEN
        SET @parent_scenario_id = (SELECT scenario_id FROM tests WHERE id = NEW.depends_on_id);
        IF @parent_scenario_id != NEW.scenario_id THEN
            SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'depends_on test must be in the same scenario';
        END IF;
    END IF;
END;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop trigger before_test_update;
-- +goose StatementEnd
-- +goose StatementBegin
drop trigger before_test_insert;
-- +goose StatementEnd
-- +goose StatementBegin
drop table tests;
-- +goose StatementEnd
