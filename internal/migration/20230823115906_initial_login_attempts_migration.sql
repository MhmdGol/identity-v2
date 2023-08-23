-- +goose Up
-- +goose StatementBegin
CREATE TABLE login_attempts(
    id INT IDENTITY(1,1) PRIMARY KEY,
    user_id BIGINT,
    attempts INT NOT NULL DEFAULT 0,
    last_attempt DATETIME,
    ban_expiry DATETIME,
    FOREIGN KEY (user_id) REFERENCES users(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE login_attempts;
-- +goose StatementEnd
