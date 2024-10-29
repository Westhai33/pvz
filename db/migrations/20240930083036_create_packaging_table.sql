-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS packaging (
                                         packaging_id BIGSERIAL PRIMARY KEY,
                                         type TEXT NOT NULL,
                                         cost FLOAT NOT NULL,
                                         max_weight FLOAT NOT NULL
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS packaging CASCADE;

-- +goose StatementEnd
