package db

import (
	"context"
	"database/sql"
)

type User struct {
	ID        string     `json:"id"`
	Name      NullString `json:"name"`
	Email     string     `json:"email"`
	Password  string     `json:"-"`
	CreatedAt int64      `json:"created_at" db:"created_at"`
	UpdatedAt int64      `json:"updated_at" db:"updated_at"`
}

type CreateUserParams struct {
	ID       string     `json:"id" validate:"required,min=1,max=36"`
	Name     NullString `json:"name"`
	Email    string     `json:"email" validate:"required,email"`
	Password string     `json:"password" validate:"required,min=8,max=15"`
}

type UpdateUserParams struct {
	ID       string     `json:"id" validate:"required,min=1,max=36"`
	Name     NullString `json:"name"`
	Password NullString `json:"password"`
}

const getUsers = `
SELECT id, name, email, password, created_at, updated_at
FROM users
`

func (q *Queries) GetUsers(ctx context.Context) ([]User, error) {
	var users []User
	err := q.db.SelectContext(ctx, &users, getUsers)
	if err != nil {
		return nil, err
	}
	return users, nil
}

const createUser = `
INSERT INTO users (id, name, email, password, created_at, updated_at)
VALUES (:id, :name, :email, :password, unixepoch(), unixepoch())
RETURNING id, name, email, password, created_at, updated_at
`

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (*User, error) {
	rows, err := q.db.NamedQueryContext(ctx, createUser, arg)
	if err != nil {
		return nil, err
	}
	var i User

	for rows.Next() {
		err = rows.StructScan(&i)
	}
	return &i, err
}

const getUserByEmail = `
SELECT id, name, email, password, created_at, updated_at
FROM users
WHERE email = $1
LIMIT 1
`

func (q *Queries) GetUserByEmail(ctx context.Context, email string) (User, error) {
	var i User
	err := q.db.GetContext(ctx, &i, getUserByEmail, email)
	if err == sql.ErrNoRows {
		return User{}, nil
	}
	return i, err
}

const getUserById = `
SELECT id, name, email, password, created_at, updated_at
FROM users
WHERE id = $1
LIMIT 1
`

func (q *Queries) GetUserByID(ctx context.Context, id string) (User, error) {
	var i User
	err := q.db.GetContext(ctx, &i, getUserById, id)
	if err == sql.ErrNoRows {
		return User{}, nil
	}
	return i, err
}

const updateUser = `
UPDATE users
SET name = coalesce(:name, name), password = coalesce(:password, password), updated_at = unixepoch()
WHERE id = :id 
RETURNING id, name, email, password, created_at, updated_at
`

func (q *Queries) UpdateUser(ctx context.Context, arg UpdateUserParams) (User, error) {
	rows, err := q.db.NamedQueryContext(ctx, updateUser, arg)
	if err != nil {
		return User{}, err
	}

	var i User
	for rows.Next() {
		err = rows.StructScan(&i)
	}

	return i, err
}

const deleteUser = `
DELETE FROM users
WHERE id = $1
`

func (q *Queries) DeleteUser(ctx context.Context, id string) error {
	_, err := q.db.ExecContext(ctx, deleteUser, id)
	return err
}
