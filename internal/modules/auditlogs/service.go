package auditlogs

import (
	db "hubku/lapor_warga_be_v2/internal/database/generated"
)

type LogsService interface {
	CreateLog(arg db.CreateAuditLogParams) error
	GetLogs() ([]db.AuditLog, error)
}

type service struct {
	repo LogsRepository
}

func NewLogsService(repo LogsRepository) LogsService {
	return &service{
		repo: repo,
	}
}

func (s *service) CreateLog(arg db.CreateAuditLogParams) error {
	return s.repo.CreateLog(arg)
}

func (s *service) GetLogs() ([]db.AuditLog, error) {
	return s.repo.GetLogs()
}
