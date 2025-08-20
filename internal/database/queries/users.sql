-- name: GetUsers :many
SELECT 
    id,
    pgp_sym_decrypt(email, @key::text) AS email,
    pgp_sym_decrypt(full_name, @key::text) AS fullname,
    pgp_sym_decrypt(phone_number, @key::text) AS phone_number,
    username,
    role,
    credibility_score,
    status,
    is_email_verified, 
    is_phone_verified,
    last_login_at,
    created_at,
    updated_at
FROM users
WHERE deleted_at IS NULL 
ORDER BY created_at DESC
OFFSET @offset_count LIMIT @limit_count;

-- name: SearchUser :many
SELECT 
    id, 
    username,
    pgp_sym_decrypt(email, @key::text) AS email,
    pgp_sym_decrypt(full_name, @key::text) AS fullname,
    pgp_sym_decrypt(phone_number, @key::text) AS phone_number,
    role, 
    is_email_verified, 
    is_phone_verified,
    last_login_at, 
    created_at, 
    updated_at
FROM users 
WHERE deleted_at IS NULL
AND (
    id = @id::uuid OR
    username ILIKE '%' || @query::text || '%' OR
    pgp_sym_decrypt(email, @key::text) ILIKE '%' || @query::text || '%' OR
    pgp_sym_decrypt(full_name, @key::text) ILIKE '%' || @query::text || '%' OR
    pgp_sym_decrypt(phone_number, @key::text) ILIKE '%' || @query::text || '%'
)
OFFSET @offset_count LIMIT @limit_count;

-- name: CreateUser :one
INSERT INTO users (
    email,
    full_name,
    username,
    password_hash,
    phone_number,
    role
)
VALUES (
    pgp_sym_encrypt(@email::text, @key::text),
    pgp_sym_encrypt(@full_name::text, @key::text),
    @username::text,
    @password_hash::text,
    pgp_sym_encrypt(@phone_number::text, @key::text),
    @role::text
)
RETURNING 
    id,
    pgp_sym_decrypt(email, @key::text) AS email,
    pgp_sym_decrypt(full_name, @key::text) AS full_name,
    pgp_sym_decrypt(phone_number, @key::text) AS phone_number,
    role,
    credibility_score,
    status,
    created_at,
    updated_at;

-- name: UpdateUser :exec
UPDATE users 
SET 
    email = CASE 
        WHEN @email::text IS NOT NULL AND @email::text != '' 
        THEN pgp_sym_encrypt(@email::text, @key::text)
        ELSE email 
    END,
    full_name = CASE
        WHEN @full_name::text IS NOT NULL AND @full_name::text != ''
        THEN pgp_sym_encrypt(@full_name::text, @key::text)
        ELSE full_name
    END,
    phone_number = CASE
        WHEN @phone_number::text IS NOT NULL AND @phone_number::text != ''
        THEN pgp_sym_encrypt(@phone_number::text, @key::text)
        ELSE phone_number
    END,
    role = CASE 
        WHEN @role::text IS NOT NULL AND @role::text IN ('citizen', 'admin', 'superadmin') 
        THEN @role::text 
        ELSE role 
    END,
    status = CASE
        WHEN @status::text IS NOT NULL AND @status::text IN ('probation', 'regular', 'suspended')
        THEN @status::text
        ELSE status
    END,
    credibility_score = COALESCE(@credibility_score::smallint, credibility_score),
    updated_at = CURRENT_TIMESTAMP
WHERE id = @id::uuid;

-- name: DeleteUser :exec
UPDATE users 
SET deleted_at = CURRENT_TIMESTAMP
WHERE id = @id::uuid;

-- name: RestoreUser :exec
UPDATE users 
SET deleted_at = NULL
WHERE id = @id::uuid;

-- name: CheckUserExists :one
SELECT EXISTS(
    SELECT 1
    FROM users
    WHERE 
        (
            @email::text IS NOT NULL
            AND pgp_sym_decrypt(email, @key::text) = @email::text
        )
        OR (
            @username::text IS NOT NULL
            AND username = @username::text
        )
) AS exists;

-- name: GetUserByEmail :one
SELECT 
    id,
    pgp_sym_decrypt(email, @key::text) AS email,
    pgp_sym_decrypt(full_name, @key::text) AS full_name,
    pgp_sym_decrypt(phone_number, @key::text) AS phone_number,
    role,
    credibility_score,
    status,
    created_at,
    updated_at
FROM users
WHERE pgp_sym_decrypt(email, @key::text) = @email::text
LIMIT 1;

-- name: GetUserByID :one
SELECT 
    id,
    pgp_sym_decrypt(email, @key::text) AS email,
    pgp_sym_decrypt(full_name, @key::text) AS full_name,
    pgp_sym_decrypt(phone_number, @key::text) AS phone_number,
    role,
    credibility_score,
    status,
    created_at,
    updated_at
FROM users
WHERE id = @id::uuid
LIMIT 1;

-- name: GetUserByIdentifier :one
SELECT 
    id,
    pgp_sym_decrypt(email, @key::text) AS email,
    pgp_sym_decrypt(full_name, @key::text) AS full_name,
    pgp_sym_decrypt(phone_number, @key::text) AS phone_number,
    username,
    role,
    credibility_score,
    status,
    password_hash,
    auth_provider,
    oauth_id,
    locked_until,
    failed_login_attempts,
    created_at,
    updated_at
FROM users
WHERE 
    (id = @id::uuid)
    OR (NULLIF(@email::text, '') IS NOT NULL AND pgp_sym_decrypt(email, @key::text) = @email::text)
    OR (NULLIF(@username::text, '') IS NOT NULL AND username = @username::text)
LIMIT 1;

-- name: UpdateLastLogin :exec
UPDATE users 
SET last_login_at = NOW()
WHERE id = @id::uuid;

-- name: IncrementFailedLoginCount :exec
UPDATE users 
SET failed_login_attempts = failed_login_attempts + 1
WHERE id = @id::uuid;
