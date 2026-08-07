package models

type Task struct {
	ID          int
	Title       string
	Description string
	Completed   bool
	UserID      int
}
