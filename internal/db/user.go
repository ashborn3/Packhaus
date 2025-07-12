package db

import "database/sql"

type User struct {
	ID           int
	Username     string
	Email        string
	PasswordHash string
}

func CreateUser(db *sql.DB, username, email, pwrdhash string) (User, error) {
	var user User
	err := db.QueryRow(`
		insert into users (username, email, password_hash)
		values ($1, $2, $3)
		returning id, username, email, password_hash
	`, username, email, pwrdhash).Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash)

	return user, err
}
