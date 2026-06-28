package domain

import "time"

type User struct {
	ID           string    `db:"id" json:"id"`
	Username     string    `db:"username" json:"username"`
	Email        string    `db:"email" json:"email"`
	PasswordHash string    `db:"password_hash" json:"-"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time `db:"updated_at" json:"updated_at"`
}

type UserPublic struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type UserSummary struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UpdateUserRequest struct {
	Username string `json:"username,omitempty"`
}

type AuthResponse struct {
	AccessToken string     `json:"access_token"`
	TokenType   string     `json:"token_type"`
	ExpiresIn   int        `json:"expires_in"`
	User        UserPublic `json:"user"`
}

func (u *User) ToPublic() UserPublic {
	return UserPublic{
		ID:       u.ID,
		Username: u.Username,
		Email:    u.Email,
	}
}

func (u *User) ToSummary() UserSummary {
	return UserSummary{
		ID:       u.ID,
		Username: u.Username,
	}
}