package service

import (
	"context"
	"database/sql"
	"errors"

	"github.com/bookshelf/monolith/internal/domain"
	"github.com/bookshelf/monolith/internal/repository"
)

var (
	ErrBookNotFound    = errors.New("Book not found")
	ErrNotBookOwner    = errors.New("Not book owner")
	ErrBookTitleEmpty  = errors.New("Book title empty")
	ErrBookAuthorEmpty = errors.New("Book author empty")
)

type BookService struct {
	bookRepo *repository.BookRepository
	userRepo *repository.UserRepository
}

func (s *BookService) Create(ctx context.Context, userID string, req domain.CreateBookRequest) (*domain.BookResponse, error) {
	if req.Author == "" {
		return nil, ErrBookAuthorEmpty
	}

	if req.Title == "" {
		return nil, ErrBookTitleEmpty
	}

	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, repository.ErrUserNotFound
		}
		return nil, err
	}

	book := &domain.Book{
		Title:     req.Title,
		Author:    req.Author,
		CreatedBy: userID,
	}

	if req.ISBN != nil {
		book.ISBN = sql.NullString{
			String: *req.ISBN,
			Valid:  true,
		}
	}

	if req.PublishedYear != nil {
		book.PublishedYear = sql.NullInt32{
			Int32: int32(*req.PublishedYear),
			Valid: true,
		}
	}

	if err := s.bookRepo.Create(ctx, book); err != nil {
		return nil, err
	}

	response := book.ToResponse()

	summary := user.ToSummary()
	response.Creator = &summary

	return &response, nil
}

func (s *BookService) GetByID(ctx context.Context, id string) (*domain.BookResponse, error) {
	book, err := s.bookRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if book == nil {
		return nil, ErrBookNotFound
	}

	user, err := s.userRepo.GetByID(ctx, book.CreatedBy)
	if err != nil {
		return nil, err
	}

	response := book.ToResponse()

	summary := user.ToSummary()
	response.Creator = &summary

	return &response, nil
}

func (s *BookService) List(ctx context.Context, filter domain.BookFilter) (*domain.BookListResponse, error) {
	if filter.Page <= 0 {
		filter.Page = 1
	}

	if filter.Limit <= 0 {
		filter.Limit = 10
	}

	if filter.Sort == "" {
		filter.Sort = "created_at"
	}

	if filter.Order == "" {
		filter.Order = "desc"
	}

	books, total, err := s.bookRepo.List(ctx, &filter)
	if err != nil {
		return nil, err
	}

	responseBooks := make([]domain.BookResponse, 0, len(books))
	for _, book := range books {
		responseBooks = append(responseBooks, book.ToResponse())
	}

	totalPages := 0
	if total > 0 {
		totalPages = (total + filter.Limit - 1) / filter.Limit
	}

	return &domain.BookListResponse{
		Data: responseBooks,
		Pagination: domain.Pagination{
			Page:       filter.Page,
			Limit:      filter.Limit,
			Total:      total,
			TotalPages: totalPages,
		},
	}, nil
}

func (s *BookService) Update(ctx context.Context, userID string, bookID string, req domain.UpdateBookRequest) (*domain.BookResponse, error) {
	book, err := s.bookRepo.GetByID(ctx, bookID)
	if err != nil {
		return nil, err
	}

	if book == nil {
		return nil, ErrBookNotFound
	}

	if book.CreatedBy != userID {
		return nil, ErrNotBookOwner
	}

	if req.Title != nil {
		if *req.Title == "" {
			return nil, ErrBookTitleEmpty
		}

		book.Title = *req.Title
	}

	if req.Author != nil {
		if *req.Author == "" {
			return nil, ErrBookAuthorEmpty
		}

		book.Author = *req.Author
	}

	if req.Description != nil {
		book.Description = sql.NullString{
			String: *req.Description,
			Valid:  true,
		}
	}

	if req.ISBN != nil {
		book.ISBN = sql.NullString{
			String: *req.ISBN,
			Valid:  true,
		}
	}

	if req.PublishedYear != nil {
		book.PublishedYear = sql.NullInt32{
			Int32: int32(*req.PublishedYear),
			Valid: true,
		}
	}

	err = s.bookRepo.Update(ctx, book)
	if err != nil {
		return nil, err
	}

	response := book.ToResponse()

	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	summary := user.ToSummary()
	response.Creator = &summary

	return &response, nil
}

func (s *BookService) Delete(ctx context.Context, userID string, bookID string) error {
	book, err := s.bookRepo.GetByID(ctx, bookID)
	if err != nil {
		return err
	}

	if book == nil {
		return ErrBookNotFound
	}

	if book.CreatedBy != userID {
		return ErrNotBookOwner
	}

	err = s.bookRepo.Delete(ctx, bookID)
	if err != nil {
		return err
	}

	return nil
}
