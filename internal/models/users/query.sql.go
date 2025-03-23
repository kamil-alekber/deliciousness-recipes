// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: query.sql

package users

import (
	"context"
	"time"
)

const createUser = `-- name: CreateUser :one
INSERT INTO
    users (id, email, name, given_name, family_name, picture)
VALUES
    (?, ?, ?, ?, ?, ?) RETURNING id, email, name, given_name, family_name, picture, created_at, updated_at
`

type CreateUserParams struct {
	ID         string
	Email      string
	Name       string
	GivenName  string
	FamilyName string
	Picture    string
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (*User, error) {
	row := q.db.QueryRowContext(ctx, createUser,
		arg.ID,
		arg.Email,
		arg.Name,
		arg.GivenName,
		arg.FamilyName,
		arg.Picture,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.Name,
		&i.GivenName,
		&i.FamilyName,
		&i.Picture,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return &i, err
}

const deleteUser = `-- name: DeleteUser :exec
DELETE FROM users
WHERE
    id = ?
`

func (q *Queries) DeleteUser(ctx context.Context, id string) error {
	_, err := q.db.ExecContext(ctx, deleteUser, id)
	return err
}

const getUser = `-- name: GetUser :one
SELECT
    id, email, name, given_name, family_name, picture, created_at, updated_at
FROM
    users
WHERE
    id = ?
LIMIT
    1
`

func (q *Queries) GetUser(ctx context.Context, id string) (*User, error) {
	row := q.db.QueryRowContext(ctx, getUser, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.Name,
		&i.GivenName,
		&i.FamilyName,
		&i.Picture,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return &i, err
}

const listUsers = `-- name: ListUsers :many
SELECT
    id, email, name, given_name, family_name, picture, created_at, updated_at
FROM
    users
ORDER BY
    created_at DESC
`

func (q *Queries) ListUsers(ctx context.Context) ([]*User, error) {
	rows, err := q.db.QueryContext(ctx, listUsers)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*User
	for rows.Next() {
		var i User
		if err := rows.Scan(
			&i.ID,
			&i.Email,
			&i.Name,
			&i.GivenName,
			&i.FamilyName,
			&i.Picture,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateUser = `-- name: UpdateUser :exec
UPDATE users
SET
    email = ?,
    name = ?,
    given_name = ?,
    family_name = ?,
    picture = ?,
    updated_at = ?
WHERE
    id = ?
`

type UpdateUserParams struct {
	Email      string
	Name       string
	GivenName  string
	FamilyName string
	Picture    string
	UpdatedAt  time.Time
	ID         string
}

func (q *Queries) UpdateUser(ctx context.Context, arg UpdateUserParams) error {
	_, err := q.db.ExecContext(ctx, updateUser,
		arg.Email,
		arg.Name,
		arg.GivenName,
		arg.FamilyName,
		arg.Picture,
		arg.UpdatedAt,
		arg.ID,
	)
	return err
}
