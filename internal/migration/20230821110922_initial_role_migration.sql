-- +goose Up
-- +goose StatementBegin
CREATE TABLE roles(
    id INT IDENTITY(1,1) PRIMARY KEY,
    name VARCHAR(255)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE roles;
-- +goose StatementEnd
