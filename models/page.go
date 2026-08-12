package models

type PageData struct {
	Title string
	Name  string
	User  User

	// Messages
	Error            string
	Success          string
	VerificationCode string

	VerificationMethod string

	// Search and filters
	Search   string
	Priority string
	Category string
	Status   string
	Sort     string

	// Dashboard
	TotalTasks     int
	CompletedTasks int
	PendingTasks   int

	HighPriorityTasks int
	DueTodayTasks     int
	CompletionRate    int

	// Tasks
	Tasks []Task
	Task  Task

	// Recent dashboard tasks
	RecentTasks []Task

	ResetToken string
}
