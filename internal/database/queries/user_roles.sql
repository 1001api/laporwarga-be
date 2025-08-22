-- name: CreateUserRole :exec
INSERT INTO user_roles (
    user_id,
    role_type,
    created_by
) VALUES (
    @user_id::uuid,
    @role_type,
    @created_by::uuid
) ON CONFLICT (user_id, role_type) DO NOTHING;

-- name: GetUserRoleByUserID :one
SELECT * FROM user_roles
WHERE user_id = @user_id AND deleted_at IS NULL;

-- name: UpdateUserRole :one
UPDATE user_roles
SET
    role_type = @role_type,
    last_updated_by = @last_updated_by::uuid,
    last_updated_at = NOW()
WHERE user_id = @user_id::uuid
RETURNING *;

-- name: DeleteUserRole :exec
UPDATE user_roles
SET
    deleted_at = NOW(),
    deleted_by = @deleted_by::uuid
WHERE user_id = @user_id::uuid;