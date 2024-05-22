package service

import (
	"butler/application/domains/services/warehouse/models"
	"context"
)

type IService interface {
	GetOne(ctx context.Context, params *models.GetRequest) (*models.Warehouse, error)
	GetList(ctx context.Context, params *models.GetRequest) ([]*models.Warehouse, error)
	Update(ctx context.Context, obj *models.Warehouse) (*models.Warehouse, error)
}
