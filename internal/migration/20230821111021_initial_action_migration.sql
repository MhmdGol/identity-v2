-- +goose Up
-- +goose StatementBegin
CREATE TABLE actions(
    id INT IDENTITY(1,1) PRIMARY KEY,
    name VARCHAR(255) UNIQUE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE actions;
-- +goose StatementEnd
