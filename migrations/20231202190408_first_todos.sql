-- +goose Up
-- +goose StatementBegin
CREATE TABLE todos (
    id SERIAL PRIMARY KEY,
    title VARCHAR NOT NULL,
    description TEXT,
    date DATE NOT NULL,
    status VARCHAR NOT NULL,
    created_at timestamp DEFAULT NOW(),
    updated_at timestamp
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE todos;
-- +goose StatementEnd
