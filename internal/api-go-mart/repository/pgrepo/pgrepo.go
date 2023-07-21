package pgrepo

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"

	"github.com/alexlzrv/go-mart/internal/api-go-mart/entities"
	postgres "github.com/alexlzrv/go-mart/sql"
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

func (repo *PostgresRepo) GetUserOrders(userID int64) ([]byte, error) {
	query := `SELECT order_num, status, accrual, uploaded_at
			  FROM orders 
			  WHERE user_id = $1
			  ORDER BY uploaded_at ASC`

	rows, err := repo.db.Query(query, userID)
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
		return nil, errors.New("getUserOrders, no data")
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

	var (
		accrual sql.NullFloat64
		userID  int64
	)

	query := `INSERT INTO orders(user_id, order_num, status, uploaded_at)
				VALUES($1, $2, $3, now())`

	oldOrder := `SELECT user_id, order_num, status, accrual, uploaded_at
			 	 FROM orders 
			 	 WHERE order_num = $1`

	err = repo.db.
		QueryRowContext(ctx, oldOrder, order.Number).
		Scan(&userID, &order.Number, &order.Status, &accrual, &order.UploadedAt)

	switch err {
	case nil:
		if order.UserID == userID {
			repo.log.Infof("loadOrder, order %s already added, userID %d", order.Number, userID)
			return entities.ErrOrderAlreadyAdded
		}
		repo.log.Infof("loadOrder, order %s already added, older userID %d userID %d", order.Number, userID, order.UserID)
		return entities.ErrOrderAddedByOther

	case sql.ErrNoRows:
		repo.log.Infof("loadOrder, order %s added, userID %d", order.Number, order.UserID)
		_, err = tx.ExecContext(ctx, query, order.UserID, order.Number, order.Status)
		if err != nil {
			repo.log.Errorf("loadOrder, error loading order %d - %s", order.UserID, err)
			return err
		}

	default:
		return err
	}

	return tx.Commit()
}

func (repo *PostgresRepo) GetBalanceInfo(login string) ([]byte, error) {
	var result []byte
	return result, nil
}

func (repo *PostgresRepo) Withdraw(login string, orderID string, sum float64) error {
	return nil
}

func (repo *PostgresRepo) GetWithdrawals(login string) ([]byte, error) {
	var result []byte
	return result, nil
}
