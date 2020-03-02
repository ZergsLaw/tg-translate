-- +goose Up
-- SQL in this section is executed when the migration is applied.
create type lang_edit_state as enum ('from', 'to');

create table users
(
    id                 serial,
    telegram_id        integer                 not null,
    current_lang_from  text                    not null,
    current_lang_to    text                    not null,
    current_lang_state lang_edit_state,
    created_at         timestamp default now() not null,
    last_action_time   timestamp default now() not null,

    unique (telegram_id),
    primary key (id)
);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
drop table users;
drop type lang_edit_state;