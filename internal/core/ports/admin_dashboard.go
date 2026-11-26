package ports

import (
	"context"

	"ForgeFit_Servidor/internal/core/domain"
)

type AdminDashboardRepository interface {
	GetAdminDashboardSummary(ctx context.Context) (*domain.AdminDashboardSummary, error)
}

type AdminDashboardService interface {
	GetAdminDashboardSummary(ctx context.Context) (*domain.AdminDashboardSummary, error)
}
