CREATE TABLE IF NOT EXISTS users (
                                     user_id SERIAL PRIMARY KEY,
                                     username VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );

CREATE TABLE IF NOT EXISTS packaging (
                                         packaging_id SERIAL PRIMARY KEY,
                                         type VARCHAR(255) NOT NULL,
    cost FLOAT NOT NULL,
    max_weight FLOAT NOT NULL
    );

CREATE TABLE IF NOT EXISTS statuses (
                                        status_id SERIAL PRIMARY KEY,
                                        status_name VARCHAR(255) NOT NULL
    );

CREATE TABLE IF NOT EXISTS return_reasons (
                                              reason_id SERIAL PRIMARY KEY,
                                              reason VARCHAR(255) NOT NULL
    );


CREATE TABLE IF NOT EXISTS orders (
                                      order_id SERIAL PRIMARY KEY,
                                      user_id INT REFERENCES users(user_id),
    acceptance_date TIMESTAMP NOT NULL,
    expiration_date TIMESTAMP NOT NULL,
    weight FLOAT NOT NULL,
    base_cost FLOAT NOT NULL,
    packaging_cost FLOAT NOT NULL,
    total_cost FLOAT NOT NULL,
    packaging_id INT REFERENCES packaging(packaging_id),
    status_id INT REFERENCES statuses(status_id),
    issue_date TIMESTAMP,
    with_film BOOLEAN NOT NULL
    );

CREATE TABLE IF NOT EXISTS returns (
                                       return_id SERIAL PRIMARY KEY,
                                       order_id INT REFERENCES orders(order_id),
    user_id INT REFERENCES users(user_id),
    return_date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    reason_id INT REFERENCES return_reasons(reason_id),
    base_cost FLOAT NOT NULL,
    packaging_cost FLOAT NOT NULL,
    packaging_id INT REFERENCES packaging(packaging_id),
    total_cost FLOAT NOT NULL,
    status_id INT REFERENCES statuses(status_id)
    );


-- Заполнение таблицы статусов
INSERT INTO statuses (status_name) VALUES
                                       ('Создан'),
                                       ('Выдан'),
                                       ('Возврат'),
                                       ('Передан курьеру');

-- Заполнение таблицы причин возвратов
INSERT INTO return_reasons (reason) VALUES
                                        ('Истек срок хранения'),
                                        ('Вернул покупатель');

-- Заполнение таблицы упаковки
INSERT INTO packaging (type, cost, max_weight) VALUES
                                                   ('Коробка', 20.0, 30.0),
                                                   ('Пленка', 1.0, 0.0),
                                                   ('Пакет', 5.0, 10.0);


-- Заполнение таблицы пользователей
INSERT INTO users (username) VALUES
                                 ('Иван'),
                                 ('Мария'),
                                 ('Алексей');
