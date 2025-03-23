-- name: GetUser :one
SELECT
    *
FROM
    users
WHERE
    id = ?
LIMIT
    1;

-- name: ListUsers :many
SELECT
    *
FROM
    users
ORDER BY
    created_at DESC;

-- name: CreateUser :one
INSERT INTO
    users (id, email, name, given_name, family_name, picture)
VALUES
    (?, ?, ?, ?, ?, ?) RETURNING *;

-- name: UpdateUser :exec
UPDATE users
SET
    email = ?,
    name = ?,
    given_name = ?,
    family_name = ?,
    picture = ?,
    updated_at = ?
WHERE
    id = ?;

-- name: DeleteUser :exec
DELETE FROM users
WHERE
    id = ?;
