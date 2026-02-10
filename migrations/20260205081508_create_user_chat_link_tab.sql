-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
CREATE TABLE user_chat_link (
    user_id INTEGER NOT NULL,
    chat_id INTEGER NOT NULL REFERENCES chats(id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE user_chat_link;
-- +goose StatementEnd
