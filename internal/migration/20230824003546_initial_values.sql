-- +goose Up
-- +goose StatementBegin
INSERT INTO roles VALUES('admin');
INSERT INTO roles VALUES('staff');
INSERT INTO roles VALUES('user');
INSERT INTO statuses VALUES('active');
INSERT INTO statuses VALUES('suspend');
INSERT INTO actions VALUES('login');
INSERT INTO actions VALUES('logout');
INSERT INTO actions VALUES('pass_recovery');
INSERT INTO actions VALUES('update_pass');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM roles;
DELETE FROM statuses;
DELETE FROM actions;
-- +goose StatementEnd
