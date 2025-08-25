-- name: CreateAuditLog :exec
INSERT INTO audit_logs (
    entity_name,
    entity_id,
    action,
    metadata,
    performed_by
) VALUES (
    @entity_name::text,
    @entity_id::uuid,
    @action::text,
    @metadata::jsonb,
    @performed_by::uuid
);

-- name: GetAuditLogs :many
SELECT * FROM audit_logs
ORDER BY created_at DESC;