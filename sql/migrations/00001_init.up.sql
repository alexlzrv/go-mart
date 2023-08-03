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
                            status TEXT NOT NULL,
                            accrual DOUBLE PRECISION,
                            uploaded_at TIMESTAMP WITH TIME ZONE NOT NULL,
                            PRIMARY KEY(id),
                            FOREIGN KEY(user_id) REFERENCES users(id)
    );

    CREATE TABLE IF NOT EXISTS balance(
                            id BIGINT GENERATED ALWAYS AS IDENTITY,
                            user_id BIGINT,
                            balance DOUBLE PRECISION NOT NULL CHECK (balance >= 0),
                            PRIMARY KEY(id),
                            FOREIGN KEY(user_id) REFERENCES users(id)
    );

    DO $$
    BEGIN
        IF EXISTS (SELECT 1 FROM pg_type WHERE typname = 'balance_operation')
         THEN
            ALTER TYPE "balance_operation" ADD VALUE IF NOT EXISTS 'withdrawal';
            ALTER TYPE "balance_operation" ADD VALUE IF NOT EXISTS 'refill';
         ELSE
            CREATE TYPE "balance_operation" AS ENUM ('withdrawal', 'refill');
         END IF;
    END
    $$;

    CREATE TABLE IF NOT EXISTS withdraw(
                            id BIGINT GENERATED ALWAYS AS IDENTITY,
                            user_id BIGINT NOT NULL,
                            order_num TEXT NOT NULL,
                            amount DOUBLE PRECISION,
                            operation "balance_operation" NOT NULL,
                            processed_at TIMESTAMP WITH TIME ZONE NOT NULL,
                            PRIMARY KEY(id),
                            FOREIGN KEY(user_id) REFERENCES users(id)
    );

COMMIT;