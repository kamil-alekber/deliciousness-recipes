-- name: GetToken :one
SELECT
    *
FROM
    tokens
WHERE
    user_id = ?
    and vendor = ?
    and expiry > NOW ()
ORDER BY
    created_at DESC
LIMIT
    1;

-- name: ListTokens :many
SELECT
    *
FROM
    tokens
ORDER BY
    created_at DESC;

-- name: CreateToken :one
INSERT INTO
    tokens (
        access_token,
        refresh_token,
        token_type,
        expiry,
        expires_in,
        vendor,
        user_id
    )
VALUES
    (?, ?, ?, ?, ?, ?, ?) RETURNING *;

-- name: DeleteToken :exec
DELETE FROM tokens
WHERE
    expiry < NOW ();
