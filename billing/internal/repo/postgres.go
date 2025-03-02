package repo

import (
	"log"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

type BillingRepository interface {
}

type billingRepo struct {
	db *sqlx.DB
}

func NewBillingRepository(db *sqlx.DB) BillingRepository {
	return &billingRepo{
		db: db,
	}
}

type User struct {
	ID   int    `db:"id"`
	Name string `db:"name"`
}

func (r *billingRepo) GetActiveUsers() ([]User, error) {
	query, args, _ := sq.
		Select("id", "name").
		From("users").
		Where(sq.Eq{"status": "active"}).
		ToSql()

	var users []User
	err := r.db.Select(&users, query, args...)
	if err != nil {
		log.Printf("Ошибка при получении пользователей: %v", err)
		return nil, err
	}

	return users, nil
}
