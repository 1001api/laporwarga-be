package auditlogs

import (
	"context"
	db "hubku/lapor_warga_be_v2/internal/database/generated"

	"github.com/jackc/pgx/v5/pgxpool"
)

type LogsRepository interface {
	CreateLog(arg db.CreateAuditLogParams) error
	GetLogs() ([]db.AuditLog, error)
}

type repository struct {
	db *db.Queries
}

func NewLogsRepository(pool *pgxpool.Pool) LogsRepository {
	return &repository{
		db: db.New(pool),
	}
}

func (r *repository) CreateLog(arg db.CreateAuditLogParams) error {
	return r.db.CreateAuditLog(context.Background(), arg)
}

func (r *repository) GetLogs() ([]db.AuditLog, error) {
	return r.db.GetAuditLogs(context.Background())
}
