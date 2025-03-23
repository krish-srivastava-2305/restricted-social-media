package services

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/krish-srivastava-2305/config"
)

type User struct {
	Email            string    `json:"email"`
	Password         string    `json:"password"`
	FirstName        string    `json:"first_name"`
	LastName         string    `json:"last_name"`
	EmergencyContact string    `json:"emergency_contact"`
	DateOfBirth      time.Time `json:"date_of_birth"`
}

func RegisterUser(email, password, firstName, lastName, emergencyContact, dateOfBirth string) (User, error) {
	db := config.DB

	var exists bool
	err := db.QueryRow(context.Background(), `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`, email).Scan(&exists)
	if err != nil {
		return User{}, fmt.Errorf("database error: %w", err)
	}
	if exists {
		return User{}, fmt.Errorf("user already exists, please log in or use a different email")
	}

	_, err = db.Exec(context.Background(), `
		INSERT INTO users (email, password, first_name, last_name, emergency_contact, date_of_birth) 
		VALUES ($1, $2, $3, $4, $5, $6)`,
		email, password, firstName, lastName, emergencyContact, dateOfBirth,
	)

	if err != nil {
		return User{}, fmt.Errorf("failed to register user: %w", err)
	}

	return User{
		Email:            email,
		Password:         password,
		FirstName:        firstName,
		LastName:         lastName,
		EmergencyContact: emergencyContact,
	}, nil
}

func LoginUser(email, password string) (User, error) {
	db := config.DB

	var scannedUser User
	err := db.QueryRow(context.Background(), `SELECT email, password, first_name, last_name, emergency_contact, date_of_birth FROM users WHERE email = $1`, email).Scan(
		&scannedUser.Email,
		&scannedUser.Password,
		&scannedUser.FirstName,
		&scannedUser.LastName,
		&scannedUser.EmergencyContact,
		&scannedUser.DateOfBirth,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return User{}, fmt.Errorf("user not found")
		}
		return User{}, fmt.Errorf("error while scanning user: %w", err)
	}

	if !CheckPasswordHash(password, scannedUser.Password) {
		return User{}, fmt.Errorf("invalid password")
	}

	return scannedUser, nil
}

func GetUser(email string) (User, error) {
	if email == "" {
		return User{}, fmt.Errorf("email is required")
	}

	db := config.DB

	var user User

	err := db.QueryRow(context.Background(), `SELECT email, first_name, last_name, emergency_contact, date_of_birth FROM users WHERE email = $1`, email).Scan(
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.EmergencyContact,
		&user.DateOfBirth,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return User{}, fmt.Errorf("user not found")
		}
		return User{}, fmt.Errorf("error while scanning user: %w", err)
	}

	return user, nil
}
