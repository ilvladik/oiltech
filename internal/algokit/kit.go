package algokit

import (
	"context"
	"fmt"

	"oiltech/internal/domain"
)

type Kit struct {
	datasetRepo        domain.DatasetRepo
	transformationRepo domain.TransformationRepo
	runRepo            domain.RunRepo
}

func NewKit(
	datasetRepo domain.DatasetRepo,
	transformationRepo domain.TransformationRepo,
	runRepo domain.RunRepo,
) *Kit {
	return &Kit{
		datasetRepo:        datasetRepo,
		transformationRepo: transformationRepo,
		runRepo:            runRepo,
	}
}

func (k *Kit) LoadTransformation(ctx context.Context, code string) (*domain.Transformation, error) {
	transformation, err := k.transformationRepo.GetByCode(ctx, code)
	if err != nil {
		return nil, err
	}
	if transformation == nil {
		return nil, domain.NewDomainErrorWithMessage(domain.ErrNotFoundCode, fmt.Sprintf("transformation %s not found", code))
	}
	return transformation, nil
}

func (k *Kit) LoadSourceDatasets(ctx context.Context, transformation *domain.Transformation) ([]domain.Dataset, error) {
	datasets := make([]domain.Dataset, 0, len(transformation.SourceDatasetCodes))
	for _, code := range transformation.SourceDatasetCodes {
		dataset, err := k.datasetRepo.GetByCode(ctx, code)
		if err != nil {
			return nil, err
		}
		if dataset == nil {
			return nil, domain.NewDomainErrorWithMessage(domain.ErrNotFoundCode, fmt.Sprintf("source dataset %s not found", code))
		}
		datasets = append(datasets, *dataset)
	}
	return datasets, nil
}

func (k *Kit) LoadTargetDataset(ctx context.Context, transformation *domain.Transformation) (*domain.Dataset, error) {
	dataset, err := k.datasetRepo.GetByCode(ctx, transformation.TargetDatasetCode)
	if err != nil {
		return nil, err
	}
	if dataset == nil {
		return nil, domain.NewDomainErrorWithMessage(domain.ErrNotFoundCode, fmt.Sprintf("target dataset %s not found", transformation.TargetDatasetCode))
	}
	return dataset, nil
}
