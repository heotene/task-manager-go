package models

type PageData struct {
	Title string

	Name string

	User User

	Search string

	// Dashboard
	TotalTasks     int
	CompletedTasks int
	PendingTasks   int

	HighPriorityTasks int
	DueTodayTasks     int
	CompletionRate    int

	// Tasks
	Tasks []Task

	Task Task
}
