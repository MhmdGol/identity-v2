-- +goose Up
-- +goose StatementBegin
CREATE TABLE users(
    id BIGINT PRIMARY KEY,
    uun VARCHAR(255),
    username VARCHAR(255) UNIQUE,
    hashed_password VARCHAR(255),
    email VARCHAR(255) UNIQUE,
    created_at DATETIME,
    totp_secret VARCHAR(255),
    role_id INT,
    status_id INT,
    FOREIGN KEY (role_id) REFERENCES roles (id),
    FOREIGN KEY (status_id) REFERENCES statuses (id),
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
-- +goose StatementEnd
