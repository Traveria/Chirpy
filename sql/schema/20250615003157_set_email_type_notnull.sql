-- +goose Up
-- +goose StatementBegin
ALTER TABLE users ALTER COLUMN email SET NOT NULL; 
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users ALTER COLUMN email DROP NOT NULL;
-- +goose StatementEnd
