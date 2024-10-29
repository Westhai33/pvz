-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS returns (
                                       return_id BIGSERIAL PRIMARY KEY,
                                       order_id INT NOT NULL,
                                       user_id INT NOT NULL,
                                       return_date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                                       reason_id INT,
                                       base_cost FLOAT NOT NULL,
                                       packaging_cost FLOAT NOT NULL,
                                       packaging_id INT,
                                       total_cost FLOAT NOT NULL,
                                       status_id INT
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS returns CASCADE;

-- +goose StatementEnd
