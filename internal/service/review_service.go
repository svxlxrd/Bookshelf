package service

import (
	"context"
	"database/sql"
	"errors"

	"github.com/bookshelf/monolith/internal/domain"
	"github.com/bookshelf/monolith/internal/repository"
)

var (
	ErrReviewNotFound        = errors.New("review not found")
	ErrNotReviewOwner        = errors.New("not review owner")
	ErrAlreadyReviewed       = errors.New("already reviewed")
	ErrInvalidRating         = errors.New("invalid rating")
	ErrReviewContentTooShort = errors.New("review content too short")
)

type ReviewService struct {
	reviewRepo *repository.ReviewRepository
	bookRepo   *repository.BookRepository
	userRepo   *repository.UserRepository
}

func (s *ReviewService) Create(ctx context.Context, userID string, bookID string, req domain.CreateReviewRequest) (*domain.ReviewResponse, error) {
	book, err := s.bookRepo.GetByID(ctx, bookID)
	if err != nil {
		return nil, err
	}

	if book == nil {
		return nil, ErrBookNotFound
	}

	exists, err := s.reviewRepo.UserHasReviewedBook(ctx, userID, bookID)
	if err != nil {
		return nil, err
	}

	if exists {
		return nil, ErrAlreadyReviewed
	}

	if req.Rating < 1 || req.Rating > 5 {
		return nil, ErrInvalidRating
	}

	if len(req.Content) < 10 {
		return nil, ErrReviewContentTooShort
	}

	review := &domain.Review{
		UserID:  userID,
		BookID:  bookID,
		Rating:  req.Rating,
		Content: req.Content,
	}

	if req.Title != nil {
		review.Title = sql.NullString{
			String: *req.Title,
			Valid:  true,
		}
	}

	err = s.reviewRepo.Create(ctx, review)
	if err != nil {
		return nil, err
	}

	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	response := review.ToResponse(user)

	return &response, nil
}

func (s *ReviewService) GetByID(ctx context.Context, id string) (*domain.ReviewResponse, error) {
	review, err := s.reviewRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if review == nil {
		return nil, ErrReviewNotFound
	}

	user, err := s.userRepo.GetByID(ctx, review.UserID)
	if err != nil {
		return nil, err
	}

	response := review.ToResponse(user)

	return &response, nil
}

func (s *ReviewService) ListByBookID(ctx context.Context, bookID string, page, limit int) (*domain.ReviewListResponse, error) {
	if page <= 0 {
		page = 1
	}

	if limit <= 0 {
		limit = 10
	}

	reviews, total, err := s.reviewRepo.ListByBookID(
		ctx,
		bookID,
		page,
		limit,
	)

	if err != nil {
		return nil, err
	}

	userIDs := make([]string, 0, len(reviews))

	for _, review := range reviews {
		userIDs = append(userIDs, review.UserID)
	}

	users, err := s.userRepo.GetByIDs(ctx, userIDs)
	if err != nil {
		return nil, err
	}

	reviewsResponse := make([]domain.ReviewResponse, 0, len(reviews))

	for _, review := range reviews {

		user := users[review.UserID]

		response := review.ToResponse(user)

		reviewsResponse = append(
			reviewsResponse,
			response,
		)
	}

	totalPages := 0

	if total > 0 {
		totalPages = (total + limit - 1) / limit
	}

	return &domain.ReviewListResponse{
		Data: reviewsResponse,
		Pagination: domain.Pagination{
			Page:       page,
			Limit:      limit,
			Total:      total,
			TotalPages: totalPages,
		},
	}, nil
}

func (s *ReviewService) Update(ctx context.Context, userID string, reviewID string, req domain.UpdateReviewRequest) (*domain.ReviewResponse, error) {
	review, err := s.reviewRepo.GetByID(ctx, reviewID)

	if err != nil {
		return nil, err
	}

	if review == nil {
		return nil, ErrReviewNotFound
	}

	if review.UserID != userID {
		return nil, ErrNotReviewOwner
	}

	if req.Rating != nil {

		if *req.Rating < 1 || *req.Rating > 5 {
			return nil, ErrInvalidRating
		}

		review.Rating = *req.Rating
	}

	if req.Content != nil {

		if len(*req.Content) < 10 {
			return nil, ErrReviewContentTooShort
		}

		review.Content = *req.Content
	}

	if req.Title != nil {

		review.Title = sql.NullString{
			String: *req.Title,
			Valid:  true,
		}
	}

	err = s.reviewRepo.Update(ctx, review)

	if err != nil {
		return nil, err
	}

	user, err := s.userRepo.GetByID(ctx, review.UserID)

	if err != nil {
		return nil, err
	}

	response := review.ToResponse(user)

	return &response, nil
}

func (s *ReviewService) Delete(ctx context.Context, userID string, reviewID string) error {
	review, err := s.reviewRepo.GetByID(ctx, reviewID)

	if err != nil {
		return err
	}

	if review == nil {
		return ErrReviewNotFound
	}

	if review.UserID != userID {
		return ErrNotReviewOwner
	}

	return s.reviewRepo.Delete(ctx, reviewID)
}
