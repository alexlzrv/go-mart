package pgrepo

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"

	"github.com/alexlzrv/go-mart/internal/api-go-mart/entities"
	postgres "github.com/alexlzrv/go-mart/sql"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type PostgresRepo struct {
	db  *postgres.Postgres
	log *zap.SugaredLogger
}

func NewRepository(db *postgres.Postgres, log *zap.SugaredLogger) *PostgresRepo {
	return &PostgresRepo{db: db, log: log}
}

func (repo *PostgresRepo) Register(ctx context.Context, user *entities.User) error {
	tx, err := repo.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer func(tx *sql.Tx) {
		err = tx.Rollback()
		if err != nil {
			return
		}
	}(tx)

	cryptPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		repo.log.Errorf("this password is not allowed: %s", err)
		return err
	}

	user.CryptPassword = cryptPassword

	var id int64

	query := `INSERT INTO users (login, password) VALUES($1, $2)
                      RETURNING id`

	err = tx.QueryRowContext(ctx, query, user.Login, user.CryptPassword).Scan(&id)
	if err != nil {
		repo.log.Errorf("register, error with scan row %s", err)
		return err
	}

	query = `INSERT INTO balance(user_id, balance) 
				VALUES($1, 0)`

	_, err = tx.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	user.ID = id

	return tx.Commit()
}

func (repo *PostgresRepo) Login(ctx context.Context, user *entities.User) error {
	query := `SELECT id, password
				   FROM users
				   WHERE login = $1`

	err := repo.db.QueryRowContext(ctx, query, user.Login).Scan(&user.ID, &user.CryptPassword)
	if !errors.Is(err, nil) && !errors.Is(err, sql.ErrNoRows) {
		repo.log.Errorf("login, error with scan row %s", err)
		return err
	}

	return nil
}

func (repo *PostgresRepo) GetUserOrders(ctx context.Context, userID int64) ([]byte, error) {
	query := `SELECT order_num, status, accrual, uploaded_at
			  FROM orders 
			  WHERE user_id = $1
			  ORDER BY uploaded_at ASC`

	rows, err := repo.db.QueryContext(ctx, query, userID)
	if err != nil {
		repo.log.Errorf("getUserOrders, error with get row in query %s", err)
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err = rows.Close()
		if err != nil {
			return
		}
	}(rows)

	orders := make([]entities.Order, 0)

	for rows.Next() {
		var (
			accrual sql.NullFloat64
			order   = entities.Order{UserID: userID}
		)

		err = rows.Scan(&order.Number, &order.Status, &accrual, &order.UploadedAt)
		if err != nil {
			repo.log.Errorf("getUserOrders, error with scan rows %s", err)
			return nil, err
		}
		order.Accrual = accrual.Float64

		orders = append(orders, order)
	}

	if err = rows.Err(); err != nil {
		repo.log.Errorf("getUserOrders, rows error %s", err)
		return nil, err
	}

	if len(orders) == 0 {
		return nil, entities.ErrNoData
	}

	result, err := json.Marshal(orders)
	if err != nil {
		repo.log.Errorf("getUserOrders, error with marshal %s", err)
		return nil, err
	}

	return result, nil
}

func (repo *PostgresRepo) LoadOrder(ctx context.Context, order *entities.Order) error {
	tx, err := repo.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer func(tx *sql.Tx) {
		err = tx.Rollback()
		if err != nil {
			return
		}
	}(tx)

	query := `INSERT INTO orders(user_id, order_num, status, uploaded_at)
				VALUES($1, $2, $3, now())`

	_, err = tx.ExecContext(ctx, query, order.UserID, order.Number, order.Status)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code != pgerrcode.UniqueViolation {
			return err
		}

		existOrder, err := repo.CheckOrder(ctx, order.Number)
		if err != nil {
			repo.log.Errorf("checkOrder, error %s", err)
			return err
		}

		if order.UserID == existOrder.UserID {
			return entities.ErrOrderAlreadyAdded
		}

		return entities.ErrOrderAddedByOther
	}

	return tx.Commit()
}

func (repo *PostgresRepo) CheckOrder(ctx context.Context, number string) (*entities.Order, error) {
	query := `SELECT user_id, order_num, status, accrual, uploaded_at
			 	 FROM orders 
			 	 WHERE order_num = $1`

	var (
		accrual sql.NullFloat64
		order   = &entities.Order{}
	)

	err := repo.db.QueryRowContext(ctx, query, number).Scan(
		&order.UserID, &order.Number, &order.Status, &accrual, &order.UploadedAt,
	)
	if err != nil {
		return nil, err
	}

	order.Accrual = accrual.Float64

	return order, nil
}

func (repo *PostgresRepo) UpdateOrder(ctx context.Context, order *entities.Order) error {
	tx, err := repo.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer func(tx *sql.Tx) {
		err = tx.Rollback()
		if err != nil {
			return
		}
	}(tx)

	query := `UPDATE orders
		SET status = $1, accrual = $2
		WHERE order_num = $3`

	_, err = tx.ExecContext(ctx, query, order.Status, order.Accrual, order.Number)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (repo *PostgresRepo) GetAllOrder(ctx context.Context) ([]entities.Order, error) {
	query := `SELECT user_id, order_num, status, uploaded_at
		FROM orders 
		WHERE status = 'NEW' OR status = 'PROCESSING'
		ORDER BY uploaded_at ASC`

	rows, err := repo.db.QueryContext(ctx, query)
	if err != nil {
		repo.log.Errorf("getAllOrder, error with query %s", err)
		return nil, err
	}

	defer func(rows *sql.Rows) {
		err = rows.Close()
		if err != nil {
			return
		}
	}(rows)

	orders := make([]entities.Order, 0)

	for rows.Next() {
		order := entities.Order{}

		err = rows.Scan(&order.UserID, &order.Number, &order.Status, &order.UploadedAt)
		if err != nil {
			repo.log.Errorf("getAllOrder, error with scan row %s", err)
			return nil, err
		}

		orders = append(orders, order)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}

func (repo *PostgresRepo) GetBalanceInfo(ctx context.Context, userID int64) ([]byte, error) {
	query := `SELECT b.balance, COALESCE(w.amount, 0)
				FROM balance b
				LEFT JOIN (SELECT user_id, SUM(amount) AS amount
				    		FROM withdraw
				    		GROUP BY user_id) w ON b.user_id = w.user_id
				WHERE b.user_id = $1`

	balance := entities.Balance{UserID: userID}

	err := repo.db.QueryRowContext(ctx, query, userID).Scan(&balance.Current, &balance.Withdrawn)
	if err != nil {
		repo.log.Errorf("getBalanceInfo, error with %s", err)
		return nil, err
	}

	result, err := json.Marshal(balance)
	if err != nil {
		repo.log.Errorf("getBalanceInfo, error with marshal %s", err)
		return nil, err
	}

	return result, nil
}

func (repo *PostgresRepo) Withdraw(ctx context.Context, userID int64) ([]byte, error) {
	query := `SELECT order_num, amount, processed_at
				FROM withdraw
				WHERE user_id = $1
				ORDER BY processed_at ASC`

	rows, err := repo.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}

	defer func(rows *sql.Rows) {
		err = rows.Close()
		if err != nil {
			return
		}
	}(rows)

	withdrawals := make([]entities.BalanceChange, 0)

	for rows.Next() {
		withdrawal := entities.BalanceChange{
			UserID: userID,
		}

		err = rows.Scan(&withdrawal.Order, &withdrawal.Amount, &withdrawal.ProcessedAt)
		if err != nil {
			return nil, err
		}

		withdrawals = append(withdrawals, withdrawal)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	result, err := json.Marshal(withdrawals)
	if err != nil {
		repo.log.Errorf("withdraw, error with marshal %s", err)
		return nil, err
	}

	return result, nil
}

func (repo *PostgresRepo) GetWithdrawals(ctx context.Context, change *entities.BalanceChange) error {
	queryBalance := `UPDATE balance
				SET balance = balance - $1
				WHERE user_id = $2`

	queryWithdraw := `INSERT INTO withdraw(user_id, order_num, amount, processed_at)
						VALUES ($1, $2, $3, now())`

	tx, err := repo.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer func(tx *sql.Tx) {
		err = tx.Rollback()
		if err != nil {
			return
		}
	}(tx)

	_, err = tx.ExecContext(ctx, queryBalance, change.Amount, change.UserID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.CheckViolation {
			return entities.ErrNegativeBalance
		}

		repo.log.Errorf("getWithdrawals, error %s", err)
		return err
	}

	_, err = tx.ExecContext(ctx, queryWithdraw, change.UserID, change.Order, change.Amount)
	if err != nil {
		repo.log.Errorf("getWithdrawals, error exec %s", err)
		return err
	}

	return tx.Commit()
}
