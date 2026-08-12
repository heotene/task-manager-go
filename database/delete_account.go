package database

// DeleteAccount permanently deletes a user's account
// and all tasks belonging to that user.
func DeleteAccount(userID int) error {

	// Delete the user's tasks first.
	_, err := DB.Exec(
		"DELETE FROM tasks WHERE user_id = ?",
		userID,
	)

	if err != nil {
		return err
	}

	// Delete the user account.
	_, err = DB.Exec(
		"DELETE FROM users WHERE id = ?",
		userID,
	)

	return err
}
