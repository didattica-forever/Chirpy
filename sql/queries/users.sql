-- name: CreateUser :one
INSERT INTO users (
        id,
        created_at,
        updated_at,
        email,
        hashed_password
    )
VALUES (
        gen_random_uuid(),
        NOW(),
        NOW(),
        $1,
        $2
    )
RETURNING *;

-- name: GetUserByEmail :one
SELECT *
FROM users
WHERE email = $1;

-- name: UpdateUser :one
update users
SET updated_at = NOW(),
    email = $1,
    hashed_password = $2
where id = $3
RETURNING *;

-- name: GetUserById :one
SELECT *
FROM users
WHERE id = $1;

-- name: UpgradeToChirpyRed :one
update users
SET is_chirpy_red = true,
    updated_at = NOW()
WHERE id = $1
RETURNING *;