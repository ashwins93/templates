package routes

import (
	"context"

	"github.com/ashwins93/fiber-sql/db"
	"golang.org/x/crypto/bcrypt"
)

func seedDataIntoDb(q *db.Queries) error {
	userList := []struct {
		id       string
		name     string
		email    string
		password string
	}{
		{"1", "John Doe", "johndoe@example.com", "password"},
		{"2", "Jane Doe", "janedoe@example.com", "password"},
		{"3", "Ashwin", "ashwin@example.com", "password"},
	}

	for _, user := range userList {
		name := db.NullString{}
		name.String = user.name
		name.Valid = true

		_, err := q.CreateUser(context.Background(), db.CreateUserParams{
			ID:       user.id,
			Name:     name,
			Email:    user.email,
			Password: user.email,
		})

		if err != nil {
			return err
		}
	}
	return nil
}

func hashPassword(password string) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash)
}
