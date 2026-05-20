package domain

import "time"

type Algorithm struct {
	ID          string
	Code        string
	Name        string
	Description string
	RunURL      string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
