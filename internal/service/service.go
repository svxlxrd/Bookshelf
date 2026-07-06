package service

import "github.com/bookshelf/monolith/internal/repository"

type Service struct {
	User   UserService
	Book   BookService
	Review ReviewService
}

func New(repos *repository.Repository, jwtSecret string) *Service {
	return &Service{
		User: UserService{
			repo:      repos.User,
			jwtSecret: jwtSecret,
		},

		Book: BookService{
			bookRepo: repos.Book,
			userRepo: repos.User,
		},

		Review: ReviewService{
			reviewRepo: repos.Review,
			bookRepo:   repos.Book,
			userRepo:   repos.User,
		},
	}
}
