package repo

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/Oxeeee/bank-microservices/billing/internal/models/domain"
	custerrors "github.com/Oxeeee/bank-microservices/billing/internal/models/errors"
	"github.com/Oxeeee/bank-microservices/billing/internal/models/requests"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type BillingRepository interface {
	GetPaymentByID(uuid uuid.UUID) (*domain.BillPayment, error)
	ProcessPayment(req *requests.BillPayment) (uuid.UUID, error)
}

type billingRepo struct {
	db  *sqlx.DB
	log *slog.Logger
}

var sqb = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

func NewBillingRepository(db *sqlx.DB, log *slog.Logger) BillingRepository {
	return &billingRepo{
		db:  db,
		log: log,
	}
}
func (r *billingRepo) GetPaymentByID(uuid uuid.UUID) (*domain.BillPayment, error) {
	const op = "repo.GetPaymentByID"
	log := r.log.With(slog.String("op", op))

	query, args, _ := sqb.
		Select("id, user_id, provider, amount, currency, details, status, created_at, updated_at").
		From("bill_payments").
		Where(sq.Eq{"id": uuid}).
		ToSql()

	var bill domain.BillPayment
	err := r.db.Get(&bill, query, args...)
	if err != nil {
		log.Error("error while get bill by id", "error", err)
		return nil, err
	}

	return &bill, nil
}

func (r *billingRepo) ProcessPayment(req *requests.BillPayment) (uuid.UUID, error) {
	const op = "repo.processPayment"
	log := r.log.With(slog.String("op", op))

	ctx := context.Background()

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		log.Error("cannot start transaction", "error", err)
		return uuid.Nil, err
	}

	var user domain.User
	query, args, _ := sqb.
		Select("id, balance").
		From("users").
		Where(sq.Eq{"id": req.UserID}).
		ToSql()

	err = tx.QueryRowContext(ctx, query, args...).Scan(&user.ID, &user.Balance)
	if err != nil {
		tx.Rollback()
		if errors.Is(err, sql.ErrNoRows) {
			log.Debug("user not found", "userID", req.UserID)
			return uuid.Nil, custerrors.ErrUserNotFound
		}
		log.Error("cannot select user by ID", "error", err, "userID", req.UserID)
		return uuid.Nil, err
	}

	if user.Balance < req.Amount {
		log.Debug("user balance insufficient", "user_balance", user.Balance, "amount", req.Amount)
		tx.Rollback()
		return uuid.Nil, custerrors.ErrInsufficientBalance
	}

	// Обновляем баланс пользователя
	query, args, _ = sqb.Update("users").
		Set("balance", user.Balance-req.Amount).
		Where(sq.Eq{"id": user.ID}).
		ToSql()

	_, err = tx.ExecContext(ctx, query, args...)
	if err != nil {
		log.Error("cannot update user balance", "error", err, "userID", req.UserID)
		tx.Rollback()
		return uuid.Nil, err
	}

	// Вставляем запись о платеже в таблицу bill_payments
	var paymentID uuid.UUID
	query, args, _ = sqb.Insert("bill_payments").
		Columns("user_id, provider, amount, currency, details, status, created_at, updated_at").
		Values(user.ID, req.Provider, req.Amount, req.Currency, req.Details, "pending", time.Now(), time.Now()).
		Suffix("RETURNING id").
		ToSql()

	err = tx.QueryRowContext(ctx, query, args...).Scan(&paymentID)
	if err != nil {
		log.Error("cannot insert transaction into bill_payments", "error", err)
		tx.Rollback()
		return uuid.Nil, err
	}

	// Фиксируем транзакцию
	err = tx.Commit()
	if err != nil {
		log.Error("cannot commit transaction", "error", err)
		return uuid.Nil, err
	}

	log.Info("Payment processed successfully", "userID", user.ID, "paymentID", paymentID)
	return paymentID, nil
}
