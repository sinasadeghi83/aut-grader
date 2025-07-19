-- +goose Up
-- +goose StatementBegin
create table sections(
    id bigint unsigned primary key auto_increment,
    name varchar(60) not null,
    depends_on_id bigint unsigned null,
    project_id bigint unsigned not null,
    created_at datetime not null default current_timestamp,
    updated_at datetime not null default current_timestamp on update current_timestamp,
    deleted_at datetime default null,

    foreign key (depends_on_id) references sections(id) on delete cascade,
    foreign key (project_id) references projects(id) on delete cascade,

    index idx_sections_project_name (project_id, name),
    index idx_sections_depends_on (depends_on_id),
    index idx_sections_deleted (deleted_at)
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TRIGGER before_section_insert
BEFORE INSERT ON sections
FOR EACH ROW
BEGIN
    IF NEW.depends_on_id IS NOT NULL THEN
        SET @parent_project_id = (SELECT project_id FROM sections WHERE id = NEW.depends_on_id);
        IF @parent_project_id != NEW.project_id THEN
            SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'depends_on section must be in the same project';
        END IF;
    END IF;
END;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TRIGGER before_section_update
BEFORE UPDATE ON sections
FOR EACH ROW
BEGIN
    IF NEW.depends_on_id IS NOT NULL THEN
        SET @parent_project_id = (SELECT project_id FROM sections WHERE id = NEW.depends_on_id);
        IF @parent_project_id != NEW.project_id THEN
            SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'depends_on section must be in the same project';
        END IF;
    END IF;
END;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop trigger before_section_update;
-- +goose StatementEnd
-- +goose StatementBegin
drop trigger before_section_insert;
-- +goose StatementEnd
-- +goose StatementBegin
drop table sections;
-- +goose StatementEnd
