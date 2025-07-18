package repository

import (
	"batch/interface/repository"
)

type BatchRepositoryImpl struct {
	ReportRepo repository.ReportRepository
}

func (r *BatchRepositoryImpl) Report() repository.ReportRepository {
	return r.ReportRepo
}
