-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS statuses (
                                        status_id BIGSERIAL PRIMARY KEY,
                                        status_name TEXT NOT NULL
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS statuses CASCADE;

-- +goose StatementEnd
