-- +goose Up
-- +goose StatementBegin
ALTER TABLE positions ADD COLUMN id SERIAL PRIMARY KEY;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE positions DROP COLUMN id;
-- +goose StatementEnd
