-- +goose Up
-- +goose StatementBegin
create table scenarios(
    id bigint unsigned primary key auto_increment,
    name varchar(60) not null,
    depends_on_id bigint unsigned null,
    section_id bigint unsigned not null,
    created_at datetime not null default current_timestamp,
    updated_at datetime not null default current_timestamp on update current_timestamp,
    deleted_at datetime default null,

    foreign key (depends_on_id) references scenarios(id) on delete cascade,
    foreign key (section_id) references sections(id) on delete cascade,

    index idx_scenario_section_name (section_id, name),
    index idx_scenario_depends_on (depends_on_id),
    index idx_scenario_deleted (deleted_at)
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TRIGGER before_scenario_insert
BEFORE INSERT ON scenarios
FOR EACH ROW
BEGIN
    IF NEW.depends_on_id IS NOT NULL THEN
        SET @parent_section_id = (SELECT section_id FROM scenarios WHERE id = NEW.depends_on_id);
        IF @parent_section_id != NEW.section_id THEN
            SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'depends_on scenario must be in the same section';
        END IF;
    END IF;
END;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TRIGGER before_scenario_update
BEFORE UPDATE ON scenarios
FOR EACH ROW
BEGIN
    IF NEW.depends_on_id IS NOT NULL THEN
        SET @parent_section_id = (SELECT section_id FROM scenarios WHERE id = NEW.depends_on_id);
        IF @parent_section_id != NEW.section_id THEN
            SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'depends_on scenario must be in the same section';
        END IF;
    END IF;
END;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop trigger before_scenario_update;
-- +goose StatementEnd
-- +goose StatementBegin
drop trigger before_scenario_insert;
-- +goose StatementEnd
-- +goose StatementBegin
drop table scenarios;
-- +goose StatementEnd
