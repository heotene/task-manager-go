package models

import "database/sql"

type User struct {
	ID       int
	FullName string
	Email    string
	Phone    string
	Password string

	VerificationMethod string

	EmailVerified bool
	PhoneVerified bool

	VerificationCode    string
	VerificationExpires string

	ResetToken   sql.NullString
	ResetExpires sql.NullString

	CreatedAt string
	UpdatedAt string
}
