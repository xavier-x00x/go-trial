package repository

import (
	"context"

	"go-trial/internal/domain/entity"
	domainRepo "go-trial/internal/domain/repository"

	"gorm.io/gorm"
)

type numberSequenceRepository struct {
	db *gorm.DB
}

func NewNumberSequenceRepository(db *gorm.DB) domainRepo.NumberSequenceRepository {
	return &numberSequenceRepository{db: db}
}

func (r *numberSequenceRepository) FindByPrefixAndPeriod(ctx context.Context, prefix, period string) (*entity.NumberSequence, error) {
	var ns entity.NumberSequence
	err := r.db.WithContext(ctx).Where("prefix = ? AND period = ?", prefix, period).First(&ns).Error
	if err != nil {
		return nil, err
	}
	return &ns, nil
}

func (r *numberSequenceRepository) Create(ctx context.Context, ns *entity.NumberSequence) error {
	return r.db.WithContext(ctx).Create(ns).Error
}

func (r *numberSequenceRepository) GetNextNumber(ctx context.Context, prefix, period string) (int, error) {
	var ns entity.NumberSequence

	err := r.db.WithContext(ctx).Where("prefix = ? AND period = ?", prefix, period).First(&ns).Error
	if err == gorm.ErrRecordNotFound {
		ns = entity.NumberSequence{
			Prefix:     prefix,
			Period:     period,
			LastNumber: 1,
		}
		if err := r.db.WithContext(ctx).Create(&ns).Error; err != nil {
			return 0, err
		}
		return 1, nil
	}
	if err != nil {
		return 0, err
	}

	ns.LastNumber++
	if err := r.db.WithContext(ctx).Save(&ns).Error; err != nil {
		return 0, err
	}

	return ns.LastNumber, nil
}