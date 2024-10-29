-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS orders (
                                      order_id BIGSERIAL PRIMARY KEY,
                                      user_id INT NOT NULL,
                                      acceptance_date TIMESTAMP NOT NULL,
                                      expiration_date TIMESTAMP NOT NULL,
                                      weight FLOAT NOT NULL,
                                      base_cost FLOAT NOT NULL,
                                      packaging_cost FLOAT NOT NULL,
                                      total_cost FLOAT NOT NULL,
                                      packaging_id INT,
                                      status_id INT,
                                      issue_date TIMESTAMP,
                                      with_film BOOLEAN NOT NULL
);

-- Создание индекса на поле user_id для улучшения производительности запросов
CREATE INDEX idx_orders_user_id ON orders (user_id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP INDEX IF EXISTS idx_orders_user_id;  -- Удаление индекса при откате
DROP TABLE IF EXISTS orders CASCADE;

-- +goose StatementEnd
