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
	CreatedAt int64      `json:"created_at"`
	UpdatedAt int64      `json:"updated_at"`
}

type CreateUserParams struct {
	ID       string     `json:"id" validate:"required,min=1,max=36"`
	Name     NullString `json:"name"`
	Email    string     `json:"email" validate:"required,email"`
	Password string     `json:"password" validate:"required,min=8,max=15"`
}

type UpdateUserParams struct {
	Name     NullString `json:"name"`
	Password NullString `json:"password"`
}

const getUsers = `
SELECT name, email, password, created_at, updated_at
FROM users
`

func (q *Queries) GetUsers(ctx context.Context) ([]User, error) {
	rows, err := q.db.QueryContext(ctx, getUsers)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []User{}
	for rows.Next() {
		var i User
		if err := rows.Scan(
			&i.Name,
			&i.Email,
			&i.Password,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const createUser = `
INSERT INTO users (id, name, email, password, created_at, updated_at)
VALUES ($1, $2, $3, $4, unixepoch(), unixepoch())
RETURNING id, name, email, password, created_at, updated_at
`

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, createUser, arg.ID, arg.Name, arg.Email, arg.Password)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Email,
		&i.Password,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getUserByEmail = `
SELECT id, name, email, password, created_at, updated_at
FROM users
WHERE email = $1;
`

func (q *Queries) GetUserByEmail(ctx context.Context, email string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserByEmail, email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Email,
		&i.Password,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return User{}, nil
	}
	return i, err
}

const getUserById = `
SELECT id, name, email, password, created_at, updated_at
FROM users
WHERE id = $1;
`

func (q *Queries) GetUserByID(ctx context.Context, id string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserById, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Email,
		&i.Password,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return User{}, nil
	}
	return i, err
}

const updateUser = `
UPDATE users
SET name = coalesce($1, name), password = coalesce($2, password), updated_at = unixepoch()
WHERE id = $3
RETURNING id, name, email, password, created_at, updated_at
`

func (q *Queries) UpdateUser(ctx context.Context, arg UpdateUserParams, id string) (User, error) {
	row := q.db.QueryRowContext(ctx, updateUser, arg.Name, arg.Password, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Email,
		&i.Password,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
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
