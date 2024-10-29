-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS return_reasons (
                                              reason_id BIGSERIAL PRIMARY KEY,
                                              reason TEXT NOT NULL
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS return_reasons CASCADE;

-- +goose StatementEnd
