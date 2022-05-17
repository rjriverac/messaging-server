// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.13.0
// source: user.sql

package db

import (
	"context"
	"database/sql"
)

const createUser = `-- name: CreateUser :one
INSERT INTO "Users" (
    name,
    email,
    hashed_pw,
    image,
    status
  )
VALUES ($1, $2, $3, $4, $5)
RETURNING $1,
  $2,
  $4,
  $5
`

type CreateUserParams struct {
	Name     string         `json:"name"`
	Email    string         `json:"email"`
	HashedPw string         `json:"hashedPw"`
	Image    sql.NullString `json:"image"`
	Status   sql.NullString `json:"status"`
}

type CreateUserRow struct {
	Column1 interface{} `json:"column1"`
	Column2 interface{} `json:"column2"`
	Column3 interface{} `json:"column3"`
	Column4 interface{} `json:"column4"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (CreateUserRow, error) {
	row := q.db.QueryRowContext(ctx, createUser,
		arg.Name,
		arg.Email,
		arg.HashedPw,
		arg.Image,
		arg.Status,
	)
	var i CreateUserRow
	err := row.Scan(
		&i.Column1,
		&i.Column2,
		&i.Column3,
		&i.Column4,
	)
	return i, err
}

const getUser = `-- name: GetUser :one
SELECT id,
  name,
  email,
  image,
  status
FROM "Users"
WHERE id = $1
LIMIT 1
`

type GetUserRow struct {
	ID     int64          `json:"id"`
	Name   string         `json:"name"`
	Email  string         `json:"email"`
	Image  sql.NullString `json:"image"`
	Status sql.NullString `json:"status"`
}

func (q *Queries) GetUser(ctx context.Context, id int64) (GetUserRow, error) {
	row := q.db.QueryRowContext(ctx, getUser, id)
	var i GetUserRow
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Email,
		&i.Image,
		&i.Status,
	)
	return i, err
}

const listUsers = `-- name: ListUsers :many
SELECT id,
  name,
  email,
  image,
  status
FROM "Users"
ORDER BY id
LIMIT $1 OFFSET $2
`

type ListUsersParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

type ListUsersRow struct {
	ID     int64          `json:"id"`
	Name   string         `json:"name"`
	Email  string         `json:"email"`
	Image  sql.NullString `json:"image"`
	Status sql.NullString `json:"status"`
}

func (q *Queries) ListUsers(ctx context.Context, arg ListUsersParams) ([]ListUsersRow, error) {
	rows, err := q.db.QueryContext(ctx, listUsers, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListUsersRow
	for rows.Next() {
		var i ListUsersRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Email,
			&i.Image,
			&i.Status,
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

const updateUserImage = `-- name: UpdateUserImage :exec
UPDATE "Users"
SET image = $2
WHERE id = $1
`

type UpdateUserImageParams struct {
	ID    int64          `json:"id"`
	Image sql.NullString `json:"image"`
}

func (q *Queries) UpdateUserImage(ctx context.Context, arg UpdateUserImageParams) error {
	_, err := q.db.ExecContext(ctx, updateUserImage, arg.ID, arg.Image)
	return err
}

const updateUserInfo = `-- name: UpdateUserInfo :exec
UPDATE "Users"
SET (name, email, image, status, hashed_pw) = ($2, $3, $4, $5, $6)
where id = $1
returning $1, $2, $3, $4, $5
`

type UpdateUserInfoParams struct {
	ID       int64          `json:"id"`
	Name     string         `json:"name"`
	Email    string         `json:"email"`
	Image    sql.NullString `json:"image"`
	Status   sql.NullString `json:"status"`
	HashedPw string         `json:"hashedPw"`
}

type UpdateUserInfoRow struct {
	Column1 interface{} `json:"column1"`
	Column2 interface{} `json:"column2"`
	Column3 interface{} `json:"column3"`
	Column4 interface{} `json:"column4"`
	Column5 interface{} `json:"column5"`
}

func (q *Queries) UpdateUserInfo(ctx context.Context, arg UpdateUserInfoParams) error {
	_, err := q.db.ExecContext(ctx, updateUserInfo,
		arg.ID,
		arg.Name,
		arg.Email,
		arg.Image,
		arg.Status,
		arg.HashedPw,
	)
	return err
}
