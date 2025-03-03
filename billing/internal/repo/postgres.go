package repo

import (
	"log/slog"

	sq "github.com/Masterminds/squirrel"
	"github.com/Oxeeee/bank-microservices/billing/internal/models/domain"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type BillingRepository interface {
	GetUserByID(uuid uuid.UUID) (*domain.User, error)
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
	log := slog.With(
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
