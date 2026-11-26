package services

import (
	"context"

	"ForgeFit_Servidor/internal/core/domain"
	"ForgeFit_Servidor/internal/core/ports"
)

type adminDashboardService struct {
	repository ports.AdminDashboardRepository
}

func NewAdminDashboardService(repository ports.AdminDashboardRepository) ports.AdminDashboardService {
	return &adminDashboardService{repository: repository}
}

func (s *adminDashboardService) GetAdminDashboardSummary(ctx context.Context) (*domain.AdminDashboardSummary, error) {
	return s.repository.GetAdminDashboardSummary(ctx)
}
