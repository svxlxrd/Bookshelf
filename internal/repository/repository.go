package repository

import "github.com/jmoiron/sqlx"

type Repository struct {
	User *UserRepository
}

func New(db *sqlx.DB) *Repository {
	return &Repository{
		User: NewUserRepository(db),
	}
}
