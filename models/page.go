package models

type PageData struct {
	Title          string
	Name           string
	User           User
	Task           Task
	Tasks          []Task
	TotalTasks     int
	CompletedTasks int
	PendingTasks   int
	Message        string
	Error          string
}
