-- +goose Up
-- +goose StatementBegin
CREATE TABLE tracks(
    id INT IDENTITY(1,1) PRIMARY KEY,
    user_id BIGINT,
    action_id INT,
    action_time DATETIME,
    FOREIGN KEY (user_id) REFERENCES users (id),
    FOREIGN KEY (action_id) REFERENCES actions (id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE tracks;
-- +goose StatementEnd
