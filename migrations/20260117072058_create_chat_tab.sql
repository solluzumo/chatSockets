-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
CREATE TYPE chat_status_enum AS ENUM ('Приватный','Публичный','Канал');

CREATE TABLE chats (
    id SERIAL PRIMARY KEY,
    title VARCHAR(200) UNIQUE NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    chat_status chat_status_enum NOT NULL
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE chats;
DROP TYPE chat_status_enum;

-- +goose StatementEnd
