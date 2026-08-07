package database

import (
	"Myjob/models"
)

// Create a new task
func CreateTask(task models.Task) error {

	_, err := DB.Exec(
		"INSERT INTO tasks(title, description, user_id) VALUES(?,?,?)",
		task.Title,
		task.Description,
		task.UserID,
	)

	return err
}

// Get all tasks for a user
func GetTasksByUser(userID int) ([]models.Task, error) {

	rows, err := DB.Query(
		"SELECT id, title, description, completed FROM tasks WHERE user_id=?",
		userID,
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
			&task.Completed,
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
		"UPDATE tasks SET completed = 1 WHERE id=?",
		id,
	)

	return err
}

// Delete a task
func DeleteTask(id int) error {

	_, err := DB.Exec(
		"DELETE FROM tasks WHERE id=?",
		id,
	)

	return err
}
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
func GetTaskByID(id int) (models.Task, error) {

	var task models.Task

	err := DB.QueryRow(
		"SELECT id, title, description, completed, user_id FROM tasks WHERE id = ?",
		id,
	).Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&task.Completed,
		&task.UserID,
	)

	return task, err
}

func UpdateTask(task models.Task) error {

	_, err := DB.Exec(
		"UPDATE tasks SET title = ?, description = ? WHERE id = ?",
		task.Title,
		task.Description,
		task.ID,
	)

	return err
}
