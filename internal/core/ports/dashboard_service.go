package ports

import (
	"context"
	"ForgeFit_Servidor/internal/core/domain"
)

type DashboardService interface {
	GetDashboardSummary(ctx context.Context, userID string) (*domain.DashboardSummary, error)
}
