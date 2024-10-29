-- +goose Up
-- +goose StatementBegin

INSERT INTO return_reasons (reason) VALUES
                                        ('Истек срок хранения'),
                                        ('Вернул покупатель');

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DELETE FROM return_reasons WHERE reason IN ('Истек срок хранения', 'Вернул покупатель');

-- +goose StatementEnd
