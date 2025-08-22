-- name: CreateRole :one
INSERT INTO roles (
    id,
    name,
    description,
    created_at
) VALUES (
    uuid_generate_v4(),
    @name::text,
    @description::text,
    CURRENT_TIMESTAMP
) RETURNING id;

-- name: AssignRoleToUser :exec
UPDATE users
SET 
    role_id = (
        SELECT id 
        FROM roles 
        WHERE name = @role_name::text 
        AND deleted_at IS NULL
    ),
    last_updated_at = CURRENT_TIMESTAMP
WHERE id = @user_id::uuid
AND deleted_at IS NULL;

-- name: RemoveUserRole :exec
UPDATE users
SET 
    role_id = NULL,
    last_updated_at = CURRENT_TIMESTAMP
WHERE id = @user_id::uuid
AND deleted_at IS NULL;

-- name: GetUsersByRoleName :many
SELECT
    u.id,
    u.username,
    u.email_enc AS email,
    u.fullname_enc AS fullname,
    r.name AS role_name,
    u.created_at
FROM users u
JOIN roles r ON u.role_id = r.id
WHERE r.name = @role_name::text
AND u.deleted_at IS NULL
AND r.deleted_at IS NULL
ORDER BY u.created_at;

-- name: GetRoleByName :one
SELECT r.*
FROM roles r
WHERE r.name = @name::text
AND r.deleted_at IS NULL;

-- name: GetRoleByID :one
SELECT r.*
FROM roles r
WHERE r.id = @id::uuid
AND r.deleted_at IS NULL;

-- name: ListAllRoles :many
SELECT r.*
FROM roles r
WHERE r.deleted_at IS NULL
ORDER BY r.name;

-- name: HasRole :one
SELECT EXISTS (
    SELECT 1
    FROM users u
    JOIN roles r ON u.role_id = r.id
    WHERE u.id = @user_id::uuid
    AND r.name = @role_name::text
    AND u.deleted_at IS NULL
    AND r.deleted_at IS NULL
) AS has_role;

-- name: CheckRoleExists :one
SELECT EXISTS (
    SELECT 1
    FROM roles r
    WHERE r.name = @name::text
    AND r.deleted_at IS NULL
) AS exists;

-- name: DeleteRole :exec
UPDATE roles
SET 
    deleted_at = CURRENT_TIMESTAMP
WHERE id = @id::uuid;

-- name: UpdateRole :exec
UPDATE roles
SET
    name = COALESCE(NULLIF(@name::text, ''), name),
    description = COALESCE(@description::text, description),
    last_updated_at = CURRENT_TIMESTAMP
WHERE id = @id::uuid
AND deleted_at IS NULL;