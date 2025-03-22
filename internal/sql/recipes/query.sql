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
    created_at DESC;

-- name: CreateRecipe :one
INSERT INTO
    recipes (
        id,
        name,
        description,
        ingredients,
        instructions,
        cooking_time
    )
VALUES
    (?, ?, ?, ?, ?, ?) RETURNING *;

-- name: UpdateRecipe :exec
UPDATE recipes
SET
    name = ?,
    description = ?,
    updated_at = ?,
    ingredients = ?,
    instructions = ?,
    cooking_time = ?
WHERE
    id = ?;

-- name: DeleteRecipe :exec
DELETE FROM recipes
WHERE
    id = ?;
