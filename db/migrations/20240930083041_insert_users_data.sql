-- +goose Up
-- +goose StatementBegin

INSERT INTO users (username) VALUES
                                 ('Иван'),
                                 ('Мария'),
                                 ('Алексей');

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DELETE FROM users WHERE username IN ('Иван', 'Мария', 'Алексей');

-- +goose StatementEnd
