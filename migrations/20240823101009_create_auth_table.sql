-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
  name VARCHAR(100),
  email VARCHAR(255),
  password VARCHAR(255),
  role smallint
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
-- +goose StatementEnd
