BEGIN TRANSACTION;

CREATE TABLE IF NOT EXISTS users(
                        id BIGINT GENERATED ALWAYS AS IDENTITY,
                        login TEXT NOT NULL UNIQUE,
                        password TEXT NOT NULL,
                        PRIMARY KEY(id)
);

CREATE TABLE IF NOT EXISTS orders(
                        id BIGINT GENERATED ALWAYS AS IDENTITY,
                        user_id BIGINT NOT NULL,
                        order_num TEXT NOT NULL UNIQUE,
                        status TEXT,
                        accrual DOUBLE PRECISION,
                        uploaded_at TIMESTAMP WITH TIME ZONE NOT NULL,
                        PRIMARY KEY(id),
                        FOREIGN KEY(user_id) REFERENCES users(id)
);

CREATE TABLE IF NOT EXISTS balance(
                        id BIGINT GENERATED ALWAYS AS IDENTITY,
                        user_id BIGINT,
                        balance DOUBLE PRECISION,
                        FOREIGN KEY(user_id) REFERENCES users(id)
);

CREATE TABLE IF NOT EXISTS withdraw(
                        user_id BIGINT NOT NULL PRIMARY KEY REFERENCES users(id),
                        amount DOUBLE PRECISION,
                        processed_at TIMESTAMP WITH TIME ZONE NOT NULL
);

COMMIT;