-- +goose Up
-- +goose StatementBegin
CREATE TABLE sessions(
    id INT IDENTITY(1,1) PRIMARY KEY,
    user_id BIGINT UNIQUE,
    exp DATETIME,
    FOREIGN KEY (user_id) REFERENCES users(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE sessions;
-- +goose StatementEnd
