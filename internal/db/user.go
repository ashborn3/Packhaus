package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"-"` // omit in JSON responses
}

func CreateUser(db *pgxpool.Pool, username, email, hashedPassword string) (User, error) {
	var user User
	err := db.QueryRow(
		context.Background(),
		`INSERT INTO users (username, email, password_hash)
		VALUES ($1, $2, $3)
		RETURNING id, username, email`,
		username, email, hashedPassword,
	).Scan(&user.ID, &user.Username, &user.Email)

	return user, err
}

func GetUserByUsername(db *pgxpool.Pool, username string) (User, error) {
	var user User
	err := db.QueryRow(
		context.Background(),
		`SELECT id, username, email, password_hash FROM users WHERE username = $1`,
		username,
	).Scan(&user.ID, &user.Username, &user.Email, &user.Password)

	return user, err
}

func GetUserByID(db *pgxpool.Pool, id int) (User, error) {
	var user User
	err := db.QueryRow(
		context.Background(),
		`SELECT id, username, email, password_hash FROM users WHERE id = $1`,
		id,
	).Scan(&user.ID, &user.Username, &user.Email, &user.Password)

	return user, err
}
