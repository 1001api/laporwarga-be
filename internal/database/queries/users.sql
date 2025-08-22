-- name: GetUsers :many
SELECT 
    u.id,
    u.email_enc as email,
    u.fullname_enc as fullname,
    u.phone_enc as phone,
    u.username,
    ur.role_type AS role,
    u.credibility_score,
    u.status,
    u.is_email_verified, 
    u.is_phone_verified,
    u.last_login_at,
    u.created_at,
    u.last_updated_at AS updated_at
FROM users u
LEFT JOIN user_roles ur ON u.id = ur.user_id
WHERE u.deleted_at IS NULL 
ORDER BY u.created_at DESC
OFFSET @offset_count LIMIT @limit_count;

-- name: SearchUser :many
SELECT 
    u.id, 
    u.username,
    u.email_enc as email,
    u.fullname_enc as fullname,
    u.phone_enc as phone,
    ur.role_type AS role,
    u.credibility_score,
    u.status,
    u.is_email_verified, 
    u.is_phone_verified,
    u.last_login_at,
    u.created_at,
    u.last_updated_at AS updated_at
FROM users u
LEFT JOIN user_roles ur ON u.id = ur.user_id
WHERE u.deleted_at IS NULL
AND (
    u.id = @id OR
    u.username ILIKE '%' || @query::text || '%'
)
OFFSET @offset_count LIMIT @limit_count;

-- name: CreateUser :one
INSERT INTO users (
    email_hash,
    email_enc,
    fullname_hash,
    fullname_enc,
    username,
    password_hash,
    phone_hash,
    phone_enc
)
VALUES (
    @email_hash::text,
    @email_enc,
    @fullname_hash::text,
    @fullname_enc,
    @username,
    @password_hash::text,
    @phone_hash::text,
    @phone_enc
) RETURNING id;

-- name: UpdateUser :exec
UPDATE users
SET
    username = CASE
        WHEN @username::text IS NOT NULL
            AND @username::text != ''
            AND @username::text != username
        THEN @username::text
        ELSE username
    END,
    email_hash = CASE
        WHEN @email_hash::text IS NOT NULL
            AND @email_hash::text != ''
            AND @email_hash::text != email_hash
        THEN @email_hash::text
        ELSE email_hash
    END,
    email_enc = CASE
        WHEN @email_enc::bytea IS NOT NULL
            AND @email_enc::bytea != email_enc
        THEN @email_enc::bytea
        ELSE email_enc
    END,
    fullname_hash = CASE
        WHEN @fullname_hash::text IS NOT NULL
            AND @fullname_hash::text != ''
            AND @fullname_hash::text != fullname_hash
        THEN @fullname_hash::text
        ELSE fullname_hash
    END,
    fullname_enc = CASE
        WHEN @fullname_enc::bytea IS NOT NULL
            AND @fullname_enc::bytea != fullname_enc
        THEN @fullname_enc::bytea
        ELSE fullname_enc
    END,
    phone_hash = CASE
        WHEN @phone_hash::text IS NOT NULL
            AND @phone_hash::text != ''
            AND @phone_hash::text != phone_hash
        THEN @phone_hash::text
        ELSE phone_hash
    END,
    phone_enc = CASE
        WHEN @phone_enc::bytea IS NOT NULL
            AND @phone_enc::bytea != phone_enc
        THEN @phone_enc::bytea
        ELSE phone_enc
    END,
    status = CASE
        WHEN @status::text IS NOT NULL
            AND @status::text != ''
            AND @status::text != status
        THEN @status::text
        ELSE status
    END,
    last_updated_at = CASE
        WHEN (
            (@username::text IS NOT NULL AND @username::text != '' AND @username::text != username)
            OR (@email_hash::text IS NOT NULL AND @email_hash::text != '' AND @email_hash::text != email_hash)
            OR (@email_enc::bytea IS NOT NULL AND @email_enc::bytea != email_enc)
            OR (@fullname_hash::text IS NOT NULL AND @fullname_hash::text != '' AND @fullname_hash::text != fullname_hash)
            OR (@fullname_enc::bytea IS NOT NULL AND @fullname_enc::bytea != fullname_enc)
            OR (@phone_hash::text IS NOT NULL AND @phone_hash::text != '' AND @phone_hash::text != phone_hash)
            OR (@phone_enc::bytea IS NOT NULL AND @phone_enc::bytea != phone_enc)
            OR (@status::text IS NOT NULL AND @status::text != '' AND @status::text != status)
            OR (@credibility_score::smallint IS NOT NULL AND @credibility_score::smallint != credibility_score)
        ) THEN NOW()
        ELSE last_updated_at
    END,
    last_updated_by = @updated_by::uuid
WHERE id = @id;

-- name: DeleteUser :exec
UPDATE users 
SET deleted_at = NOW()
WHERE id = @id;

-- name: RestoreUser :exec
UPDATE users 
SET deleted_at = NULL
WHERE id = @id;

-- name: CheckUserExists :one
SELECT EXISTS (
    SELECT 1
    FROM users
    WHERE
        (@email_hash::text IS NOT NULL AND @email_hash != '' AND email_hash = @email_hash)
        OR (@username::text IS NOT NULL AND @username != '' AND username = @username)
) AS exists;

-- name: GetUserByEmail :one
SELECT 
    u.id,
    u.username,
    u.email_enc as email,
    u.fullname_enc as fullname,
    u.phone_enc as phone,
    ur.role_type AS role,
    u.credibility_score,
    u.status,
    u.created_at,
    u.last_updated_at AS updated_at
FROM users u
LEFT JOIN user_roles ur ON u.id = ur.user_id
WHERE u.email_hash = @email_hash
LIMIT 1;

-- name: GetUserByID :one
SELECT 
    u.id,
    u.username,
    u.email_enc as email,
    u.fullname_enc as fullname,
    u.phone_enc as phone,
    ur.role_type AS role,
    u.credibility_score,
    u.status,
    u.is_email_verified, 
    u.is_phone_verified,
    u.last_login_at,
    u.created_at,
    u.last_updated_at AS updated_at
FROM users u
LEFT JOIN user_roles ur ON u.id = ur.user_id
WHERE u.id = @id
LIMIT 1;

-- name: GetUserByIdentifier :one
SELECT 
    u.id,
    u.username,
    u.email_enc as email,
    u.fullname_enc as fullname,
    u.phone_enc as phone,
    ur.role_type AS role,
    u.credibility_score,
    u.status,
    u.password_hash,
    u.auth_provider,
    u.oauth_id,
    u.locked_until,
    u.failed_login_attempts,
    u.created_at,
    u.last_updated_at AS updated_at
FROM users u
LEFT JOIN user_roles ur ON u.id = ur.user_id
WHERE 
    (u.id = @id)
    OR (NULLIF(@email_hash, '') IS NOT NULL AND u.email_hash = @email_hash)
    OR (NULLIF(@username, '') IS NOT NULL AND u.username = @username)
LIMIT 1;

-- name: UpdateLastLogin :exec
UPDATE users 
SET last_login_at = NOW()
WHERE id = @id;

-- name: IncrementFailedLoginCount :exec
UPDATE users 
SET failed_login_attempts = failed_login_attempts + 1
WHERE id = @id;
