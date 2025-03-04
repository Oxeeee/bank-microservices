package repo

import (
	"database/sql"
	"errors"
	"log/slog"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/Oxeeee/bank-microservices/billing/internal/models/domain"
	"github.com/Oxeeee/bank-microservices/billing/internal/models/requests"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type BillingRepository interface {
	GetUserByID(uuid uuid.UUID) (*domain.User, error)
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

func (r *billingRepo) GetUserByID(uuid uuid.UUID) (*domain.User, error) {
	const op = "repo.GetUserByID"
	log := r.log.With(
		slog.String("op", op),
	)
	query, args, _ := sqb.
		Select("*").
		From("users").
		Where(sq.Eq{"id": uuid.String()}).
		ToSql()

	var user domain.User
	err := r.db.Get(&user, query, args...)
	if err != nil {
		log.Error("error while get user", "error", err)
		return nil, err
	}

	return &user, err
}
func (r *billingRepo) ProcessPayment(req *requests.BillPayment) (uuid.UUID, error) {
	const op = "repo.processPayment"
	log := r.log.With(slog.String("op", op))
	tx, err := r.db.DB.Begin()
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

	if err := tx.QueryRow(query, args...).Scan(&user.ID, &user.Balance); err != nil {
		tx.Rollback()
		if errors.Is(err, sql.ErrNoRows) {
			log.Debug("user not found", "userID", req.UserID)
			return uuid.Nil, errors.New("user not found")
		}
		log.Error("cannot select user by ID", "error", err, "userID", req.UserID)
		return uuid.Nil, err
	}

	if user.Balance < req.Amount {
		log.Debug("user balance have insufficient money", "user_balance", user.Balance, "amount", req.Amount)
		tx.Rollback()
		return uuid.Nil, errors.New("insufficient balance")
	}

	query, args, _ = sqb.Update("users").Set("balance", user.Balance-req.Amount).Where(sq.Eq{"id": user.ID}).ToSql()

	_, err = tx.Exec(query, args...)
	if err != nil {
		log.Error("cannot update user balance", "error", err, "userID", req.UserID)
		tx.Rollback()
		return uuid.Nil, err
	}

	var paymentID uuid.UUID
	query, args, _ = sqb.Insert("bill_payments").
		Columns("user_id, provider, amount, currency, details, created_at, updated_at").
		Values(user.ID, req.Provider, req.Amount, req.Currency, req.Details, time.Now(), time.Now()).
		Suffix("RETURNING id").
		ToSql()

	err = tx.QueryRow(query, args...).Scan(&paymentID)
	if err != nil {
		log.Error("cannot insert transaction to bill_payments", "error", err)
		tx.Rollback()
		return uuid.Nil, err
	}

	err = tx.Commit()
	if err != nil {
		log.Error("cannot commit transaction", "error", err)
		return uuid.Nil, err
	}

	log.Info("Payment processed successfully for user", "userID", user.ID, "paymentID", paymentID)
	return paymentID, nil
}
