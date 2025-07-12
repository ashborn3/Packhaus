package db

import "database/sql"

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"-"` // omit in JSON responses
}

func CreateUser(db *sql.DB, username, email, hashedPassword string) (User, error) {
	var user User
	err := db.QueryRow(`
		INSERT INTO users (username, email, password)
		VALUES ($1, $2, $3)
		RETURNING id, username, email`,
		username, email, hashedPassword,
	).Scan(&user.ID, &user.Username, &user.Email)

	return user, err
}
