package database

import (
	"Myjob/models"
)

func CreateUser(user models.User) error {
	_, err := DB.Exec(
		"INSERT INTO users(fullname,email,password) VALUES(?,?,?)",
		user.FullName,
		user.Email,
		user.Password,
	)

	return err
}

func GetUserByID(id int) (models.User, error) {

	var user models.User

	err := DB.QueryRow(
		"SELECT id, fullname, email FROM users WHERE id = ?",
		id,
	).Scan(
		&user.ID,
		&user.FullName,
		&user.Email,
	)

	return user, err
}

func GetUserByEmail(email string) (models.User, error) {

	var user models.User

	err := DB.QueryRow(
		"SELECT id, fullname, email, password FROM users WHERE email = ?",
		email,
	).Scan(
		&user.ID,
		&user.FullName,
		&user.Email,
		&user.Password,
	)

	return user, err
}
