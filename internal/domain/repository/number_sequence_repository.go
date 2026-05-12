package repository

import (
	"context"

	"go-trial/internal/domain/entity"
)

type NumberSequenceRepository interface {
	FindByPrefixAndPeriod(ctx context.Context, prefix, period string) (*entity.NumberSequence, error)
	Create(ctx context.Context, ns *entity.NumberSequence) error
	GetNextNumber(ctx context.Context, prefix, period string) (int, error)
}