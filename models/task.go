package models

type Task struct {
	ID          int
	Title       string
	Description string
	Priority    string
	Category    string
	DueDate     string
	Completed   bool
	CreatedAt   string
	UpdatedAt   string
	UserID      int
}
