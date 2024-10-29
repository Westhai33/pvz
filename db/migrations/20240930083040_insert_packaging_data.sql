-- +goose Up
-- +goose StatementBegin

INSERT INTO packaging (type, cost, max_weight) VALUES
                                                   ('Коробка', 20.0, 30.0),
                                                   ('Пленка', 1.0, 0.0),
                                                   ('Пакет', 5.0, 10.0);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DELETE FROM packaging WHERE type IN ('Коробка', 'Пленка', 'Пакет');

-- +goose StatementEnd
