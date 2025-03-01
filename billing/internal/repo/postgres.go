package repo

import "github.com/jmoiron/sqlx"

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
