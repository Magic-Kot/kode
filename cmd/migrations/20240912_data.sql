-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users
(
    id            SERIAL         PRIMARY KEY,
    username      VARCHAR(30)    NOT NULL        UNIQUE,
    password      VARCHAR        NOT NULL,
    name          VARCHAR        DEFAULT(''),
    surname       VARCHAR        DEFAULT(''),
    email         VARCHAR        DEFAULT('')     UNIQUE,
    age           INTEGER        DEFAULT(0),
    avatar        VARCHAR        DEFAULT('')
);

CREATE TABLE IF NOT EXISTS notes
(
    id        SERIAL         PRIMARY KEY,
    texts     VARCHAR        NOT NULL
);

CREATE TABLE IF NOT EXISTS users_notes
(
    id          SERIAL     PRIMARY KEY,
    user_id     INTEGER references users (id) on delete cascade    NOT NULL,
    note_id     INTEGER references notes (id) on delete cascade    NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users_notes;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS notes;
-- +goose StatementEnd