package database

import (
	"Myjob/models"
)

// Create a new task
func CreateTask(task models.Task) error {

	_, err := DB.Exec(`
		INSERT INTO tasks
		(title, description, priority, category, due_date, user_id)
		VALUES (?, ?, ?, ?, ?, ?)
	`,
		task.Title,
		task.Description,
		task.Priority,
		task.Category,
		task.DueDate,
		task.UserID,
	)

	return err
}

// Get all tasks for a user
func GetTasksByUser(userID int) ([]models.Task, error) {

	rows, err := DB.Query(`
		SELECT
			id,
			title,
			description,
			priority,
			category,
			due_date,
			completed,
			created_at,
			updated_at,
			user_id
		FROM tasks
		WHERE user_id = ?
		ORDER BY created_at DESC
	`, userID)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var tasks []models.Task

	for rows.Next() {

		var task models.Task

		err := rows.Scan(
			&task.ID,
			&task.Title,
			&task.Description,
			&task.Priority,
			&task.Category,
			&task.DueDate,
			&task.Completed,
			&task.CreatedAt,
			&task.UpdatedAt,
			&task.UserID,
		)

		if err != nil {
			return nil, err
		}

		tasks = append(tasks, task)
	}

	return tasks, nil
}

// Mark a user's task as completed
func CompleteTask(id int, userID int) error {

	_, err := DB.Exec(`
		UPDATE tasks
		SET completed = 1,
		    updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
		AND user_id = ?
	`,
		id,
		userID,
	)

	return err
}

// Delete a user's task
func DeleteTask(id int, userID int) error {

	_, err := DB.Exec(`
		DELETE FROM tasks
		WHERE id = ?
		AND user_id = ?
	`,
		id,
		userID,
	)

	return err
}

// Get task statistics
func GetTaskStats(userID int) (int, int, int, error) {

	var total int
	var completed int

	err := DB.QueryRow(
		"SELECT COUNT(*) FROM tasks WHERE user_id = ?",
		userID,
	).Scan(&total)

	if err != nil {
		return 0, 0, 0, err
	}

	err = DB.QueryRow(
		"SELECT COUNT(*) FROM tasks WHERE user_id = ? AND completed = 1",
		userID,
	).Scan(&completed)

	if err != nil {
		return 0, 0, 0, err
	}

	pending := total - completed

	return total, completed, pending, nil
}

// Get one task by ID
func GetTaskByID(id int) (models.Task, error) {

	var task models.Task

	err := DB.QueryRow(`
		SELECT
			id,
			title,
			description,
			priority,
			category,
			due_date,
			completed,
			created_at,
			updated_at,
			user_id
		FROM tasks
		WHERE id = ?
	`, id).Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&task.Priority,
		&task.Category,
		&task.DueDate,
		&task.Completed,
		&task.CreatedAt,
		&task.UpdatedAt,
		&task.UserID,
	)

	return task, err
}

// Update task
func UpdateTask(task models.Task) error {

	_, err := DB.Exec(`
		UPDATE tasks
		SET
			title = ?,
			description = ?,
			priority = ?,
			category = ?,
			due_date = ?,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
		AND user_id = ?
	`,
		task.Title,
		task.Description,
		task.Priority,
		task.Category,
		task.DueDate,
		task.ID,
		task.UserID,
	)

	return err
}

func SearchTasks(userID int, search string) ([]models.Task, error) {

	rows, err := DB.Query(`
		SELECT 
			id,
			title,
			description,
			priority,
			category,
			due_date,
			completed,
			user_id
		FROM tasks
		WHERE user_id = ?
		AND (
			title LIKE ?
			OR description LIKE ?
		)
	`,
		userID,
		"%"+search+"%",
		"%"+search+"%",
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var tasks []models.Task

	for rows.Next() {

		var task models.Task

		err := rows.Scan(
			&task.ID,
			&task.Title,
			&task.Description,
			&task.Priority,
			&task.Category,
			&task.DueDate,
			&task.Completed,
			&task.UserID,
		)

		if err != nil {
			return nil, err
		}

		tasks = append(tasks, task)
	}

	return tasks, nil
}
func GetDashboardAnalytics(userID int) (int, int, int, error) {

	var highPriority int
	var upcoming int
	var completionRate int

	// High-priority tasks
	err := DB.QueryRow(`
		SELECT COUNT(*)
		FROM tasks
		WHERE user_id = ?
		AND priority = 'High'
		AND completed = 0
	`, userID).Scan(&highPriority)

	if err != nil {
		return 0, 0, 0, err
	}

	// Upcoming tasks
	err = DB.QueryRow(`
		SELECT COUNT(*)
		FROM tasks
		WHERE user_id = ?
		AND completed = 0
		AND due_date >= date('now')
	`, userID).Scan(&upcoming)

	if err != nil {
		return 0, 0, 0, err
	}

	// Completion rate
	var total int
	var completed int

	err = DB.QueryRow(`
		SELECT COUNT(*)
		FROM tasks
		WHERE user_id = ?
	`, userID).Scan(&total)

	if err != nil {
		return 0, 0, 0, err
	}

	err = DB.QueryRow(`
		SELECT COUNT(*)
		FROM tasks
		WHERE user_id = ?
		AND completed = 1
	`, userID).Scan(&completed)

	if err != nil {
		return 0, 0, 0, err
	}

	if total > 0 {

		completionRate = (completed * 100) / total

	}

	return highPriority, upcoming, completionRate, nil
}

// Get the 5 most recent tasks for a user
func GetRecentTasks(userID int) ([]models.Task, error) {

	rows, err := DB.Query(`
		SELECT
			id,
			title,
			description,
			completed,
			priority,
			category,
			due_date,
			user_id
		FROM tasks
		WHERE user_id = ?
		ORDER BY id DESC
		LIMIT 5
	`, userID)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var tasks []models.Task

	for rows.Next() {

		var task models.Task

		err := rows.Scan(
			&task.ID,
			&task.Title,
			&task.Description,
			&task.Completed,
			&task.Priority,
			&task.Category,
			&task.DueDate,
			&task.UserID,
		)

		if err != nil {
			return nil, err
		}

		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

// Search and filter tasks for a user
func SearchAndFilterTasks(
	userID int,
	search string,
	priority string,
	category string,
	status string,
	sort string,
) ([]models.Task, error) {

	query := `
		SELECT
			id,
			title,
			description,
			completed,
			priority,
			category,
			due_date,
			user_id
		FROM tasks
		WHERE user_id = ?
	`

	args := []interface{}{userID}

	// Search
	if search != "" {
		query += `
			AND (
				title LIKE ?
				OR description LIKE ?
			)
		`

		searchValue := "%" + search + "%"

		args = append(
			args,
			searchValue,
			searchValue,
		)
	}

	// Priority
	if priority != "" {
		query += ` AND priority = ?`
		args = append(args, priority)
	}

	// Category
	if category != "" {
		query += ` AND category = ?`
		args = append(args, category)
	}

	// Status
	if status == "completed" {
		query += ` AND completed = 1`
	}

	if status == "pending" {
		query += ` AND completed = 0`
	}

	// Sorting
	switch sort {

	case "due_asc":
		query += ` ORDER BY due_date ASC`

	case "due_desc":
		query += ` ORDER BY due_date DESC`

	case "title":
		query += ` ORDER BY title ASC`

	default:
		query += ` ORDER BY id DESC`
	}

	rows, err := DB.Query(query, args...)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var tasks []models.Task

	for rows.Next() {

		var task models.Task

		err := rows.Scan(
			&task.ID,
			&task.Title,
			&task.Description,
			&task.Completed,
			&task.Priority,
			&task.Category,
			&task.DueDate,
			&task.UserID,
		)

		if err != nil {
			return nil, err
		}

		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

// Get a task by ID belonging to a specific user
func GetTaskByIDAndUser(taskID int, userID int) (models.Task, error) {

	var task models.Task

	err := DB.QueryRow(`
		SELECT
			id,
			title,
			description,
			priority,
			category,
			due_date,
			completed,
			created_at,
			updated_at,
			user_id
		FROM tasks
		WHERE id = ?
		AND user_id = ?
	`,
		taskID,
		userID,
	).Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&task.Priority,
		&task.Category,
		&task.DueDate,
		&task.Completed,
		&task.CreatedAt,
		&task.UpdatedAt,
		&task.UserID,
	)

	return task, err
}
