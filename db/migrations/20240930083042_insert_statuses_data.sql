-- +goose Up
-- +goose StatementBegin

INSERT INTO statuses (status_name) VALUES
                                       ('Создан'),
                                       ('Выдан'),
                                       ('Возврат'),
                                       ('Передан курьеру');

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DELETE FROM statuses WHERE status_name IN ('Создан', 'Выдан', 'Возврат', 'Передан курьеру');

-- +goose StatementEnd
