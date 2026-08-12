package database

import (
	"fmt"

	"Myjob/models"
)

func CreateUser(user models.User) (int, error) {

	result, err := DB.Exec(`
		INSERT INTO users (
			fullname,
			email,
			phone,
			password,
			verification_method,
			email_verified,
			phone_verified,
			verification_code,
			verification_expires
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`,
		user.FullName,
		user.Email,
		user.Phone,
		user.Password,
		user.VerificationMethod,
		user.EmailVerified,
		user.PhoneVerified,
		user.VerificationCode,
		user.VerificationExpires,
	)

	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func GetUserByID(id int) (models.User, error) {

	var user models.User

	err := DB.QueryRow(`
		SELECT
			id,
			fullname,
			email,
			phone,
			password,
			verification_method,
			email_verified,
			phone_verified,
			COALESCE(verification_code, ''),
			COALESCE(verification_expires, ''),
			COALESCE(reset_token, ''),
			COALESCE(reset_expires, ''),
			created_at,
			updated_at
		FROM users
		WHERE id = ?
	`, id).Scan(
		&user.ID,
		&user.FullName,
		&user.Email,
		&user.Phone,
		&user.Password,
		&user.VerificationMethod,
		&user.EmailVerified,
		&user.PhoneVerified,
		&user.VerificationCode,
		&user.VerificationExpires,
		&user.ResetToken,
		&user.ResetExpires,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	return user, err
}

func GetUserByEmail(email string) (models.User, error) {

	var user models.User

	err := DB.QueryRow(`
		SELECT
			id,
			fullname,
			email,
			phone,
			password,
			verification_method,
			email_verified,
			phone_verified,
			COALESCE(verification_code, ''),
			COALESCE(verification_expires, ''),
			COALESCE(reset_token, ''),
			COALESCE(reset_expires, ''),
			created_at,
			updated_at
		FROM users
		WHERE email = ?
	`, email).Scan(
		&user.ID,
		&user.FullName,
		&user.Email,
		&user.Phone,
		&user.Password,
		&user.VerificationMethod,
		&user.EmailVerified,
		&user.PhoneVerified,
		&user.VerificationCode,
		&user.VerificationExpires,
		&user.ResetToken,
		&user.ResetExpires,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	return user, err
}

func GetUserByPhone(phone string) (models.User, error) {

	var user models.User

	err := DB.QueryRow(`
		SELECT
			id,
			fullname,
			email,
			phone,
			password,
			verification_method,
			email_verified,
			phone_verified,
			COALESCE(verification_code, ''),
			COALESCE(verification_expires, ''),
			COALESCE(reset_token, ''),
			COALESCE(reset_expires, ''),
			created_at,
			updated_at
		FROM users
		WHERE phone = ?
	`, phone).Scan(
		&user.ID,
		&user.FullName,
		&user.Email,
		&user.Phone,
		&user.Password,
		&user.VerificationMethod,
		&user.EmailVerified,
		&user.PhoneVerified,
		&user.VerificationCode,
		&user.VerificationExpires,
		&user.ResetToken,
		&user.ResetExpires,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	return user, err
}

func UpdateUserPassword(userID int, password string) error {

	_, err := DB.Exec(`
		UPDATE users
		SET password = ?
		WHERE id = ?
	`,
		password,
		userID,
	)

	return err
}

func UpdateVerification(
	userID int,
	code string,
	expires string,
) error {

	_, err := DB.Exec(`
		UPDATE users
		SET
			verification_code = ?,
			verification_expires = ?,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`,
		code,
		expires,
		userID,
	)

	return err
}

func VerifyUser(userID int, method string) error {

	var query string

	switch method {

	case "email":

		query = `
			UPDATE users
			SET
				email_verified = 1,
				verification_code = NULL,
				verification_expires = NULL,
				updated_at = CURRENT_TIMESTAMP
			WHERE id = ?
		`

	case "phone":

		query = `
			UPDATE users
			SET
				phone_verified = 1,
				verification_code = NULL,
				verification_expires = NULL,
				updated_at = CURRENT_TIMESTAMP
			WHERE id = ?
		`

	default:
		return fmt.Errorf("invalid verification method")
	}

	_, err := DB.Exec(query, userID)

	return err
}

func UpdateResetToken(
	userID int,
	token string,
	expires string,
) error {

	_, err := DB.Exec(`
		UPDATE users
		SET
			reset_token = ?,
			reset_expires = ?,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`,
		token,
		expires,
		userID,
	)

	return err
}

func GetUserByResetToken(token string) (models.User, error) {

	var user models.User

	err := DB.QueryRow(`
		SELECT
			id,
			fullname,
			email,
			phone,
			password,
			verification_method,
			email_verified,
			phone_verified,
			COALESCE(verification_code, ''),
			COALESCE(verification_expires, ''),
			COALESCE(reset_token, ''),
			COALESCE(reset_expires, ''),
			created_at,
			updated_at
		FROM users
		WHERE reset_token = ?
	`, token).Scan(
		&user.ID,
		&user.FullName,
		&user.Email,
		&user.Phone,
		&user.Password,
		&user.VerificationMethod,
		&user.EmailVerified,
		&user.PhoneVerified,
		&user.VerificationCode,
		&user.VerificationExpires,
		&user.ResetToken,
		&user.ResetExpires,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	return user, err
}

func UpdatePassword(userID int, password string) error {

	_, err := DB.Exec(`
		UPDATE users
		SET
			password = ?,
			reset_token = NULL,
			reset_expires = NULL,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`,
		password,
		userID,
	)

	return err
}

func ClearVerification(userID int) error {

	_, err := DB.Exec(`
		UPDATE users
		SET
			verification_code = NULL,
			verification_expires = NULL,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`, userID)

	return err
}

func ClearResetToken(userID int) error {

	_, err := DB.Exec(`
		UPDATE users
		SET
			reset_token = NULL,
			reset_expires = NULL,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`, userID)

	return err
}
