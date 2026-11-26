package model

import "time"

// User merepresentasikan user yang dapat login ke API.
// Password tidak pernah ikut di-serialize ke response.
type User struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
