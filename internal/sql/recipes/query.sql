-- name: GetRecipe :one
SELECT
    *
FROM
    recipes
WHERE
    id = ?
LIMIT
    1;

-- name: ListRecipes :many
SELECT
    *
FROM
    recipes
ORDER BY
    name DESC;

-- name: CreateRecipe :one
INSERT INTO
    recipes (id, name, description)
VALUES
    (?, ?, ?) RETURNING *;

-- name: UpdateRecipe :exec
UPDATE recipes
SET
    name = ?,
    description = ?
WHERE
    id = ?;

-- name: DeleteRecipe :exec
DELETE FROM recipes
WHERE
    id = ?;
