package repository

import "github.com/jmoiron/sqlx"

type Repository struct {
	User   *UserRepository
	Book   *BookRepository
	Review *ReviewRepository
}

func New(db *sqlx.DB) *Repository {
	return &Repository{
		User:   NewUserRepository(db),
		Book:   NewBookRepository(db),
		Review: NewReviewrepository(db),
	}
}
