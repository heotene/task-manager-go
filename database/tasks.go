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

// Mark task as completed
func CompleteTask(id int) error {

	_, err := DB.Exec(
		"UPDATE tasks SET completed = 1 WHERE id = ?",
		id,
	)

	return err
}

// Delete a task
func DeleteTask(id int) error {

	_, err := DB.Exec(
		"DELETE FROM tasks WHERE id = ?",
		id,
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
	`,
		task.Title,
		task.Description,
		task.Priority,
		task.Category,
		task.DueDate,
		task.ID,
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
