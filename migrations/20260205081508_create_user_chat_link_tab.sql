-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
CREATE TYPE user_role_enum AS ENUM ('Участник','Администратор','Гость');
CREATE TABLE user_chat_links (
    user_id INTEGER NOT NULL,
    chat_id INTEGER NOT NULL REFERENCES chats(id) ON DELETE CASCADE,
    user_blocked BOOLEAN NOT NULL,
    user_role user_role_enum NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    PRIMARY KEY(user_id, chat_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE user_chat_link;
DROP TYPE user_role_enum;
-- +goose StatementEnd
