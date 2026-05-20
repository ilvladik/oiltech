package usecases

import (
	"context"
	"fmt"

	"oiltech/internal/domain"
	"oiltech/internal/dtos"
	"oiltech/internal/adapters/postgres"
)

type TransformationUsecase struct {
	datasetRepo        domain.DatasetRepo
	algorithmRepo      domain.AlgorithmRepo
	transformationRepo domain.TransformationRepo
	runUsecase         *RunUsecase
	trm                domain.TransactionManager
}

func NewTransformationUsecase(
	datasetRepo domain.DatasetRepo,
	algorithmRepo domain.AlgorithmRepo,
	transformationRepo domain.TransformationRepo,
	runUsecase *RunUsecase,
	trm domain.TransactionManager,
) *TransformationUsecase {
	return &TransformationUsecase{
		datasetRepo:        datasetRepo,
		algorithmRepo:      algorithmRepo,
		transformationRepo: transformationRepo,
		runUsecase:         runUsecase,
		trm:                trm,
	}
}

func (u *TransformationUsecase) CreateTransformation(ctx context.Context, req dtos.CreateTransformationRequest) (*dtos.TransformationResponse, error) {
	if err := postgres.ValidateDatasetCode(req.Code); err != nil {
		return nil, domain.NewDomainErrorWithMessage(domain.ErrValidationCode, err.Error())
	}
	if req.Name == "" {
		return nil, domain.NewDomainErrorWithMessage(domain.ErrValidationCode, "transformation name is required")
	}
	if len(req.SourceDatasetCodes) == 0 {
		return nil, domain.NewDomainErrorWithMessage(domain.ErrValidationCode, "source datasets are required")
	}

	existing, err := u.transformationRepo.GetByCode(ctx, req.Code)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, domain.NewDomainErrorWithMessage(domain.ErrAlreadyExistsCode, fmt.Sprintf("transformation %s already exists", req.Code))
	}

	algorithm, err := u.algorithmRepo.GetByCode(ctx, req.AlgorithmCode)
	if err != nil {
		return nil, err
	}
	if algorithm == nil {
		return nil, domain.NewDomainErrorWithMessage(domain.ErrNotFoundCode, fmt.Sprintf("algorithm %s not found", req.AlgorithmCode))
	}

	targetDataset, err := u.datasetRepo.GetByCode(ctx, req.TargetDatasetCode)
	if err != nil {
		return nil, err
	}
	if targetDataset == nil {
		return nil, domain.NewDomainErrorWithMessage(domain.ErrNotFoundCode, fmt.Sprintf("target dataset %s not found", req.TargetDatasetCode))
	}

	sourceDatasetIDs := make([]string, 0, len(req.SourceDatasetCodes))
	for _, code := range req.SourceDatasetCodes {
		sourceDataset, err := u.datasetRepo.GetByCode(ctx, code)
		if err != nil {
			return nil, err
		}
		if sourceDataset == nil {
			return nil, domain.NewDomainErrorWithMessage(domain.ErrNotFoundCode, fmt.Sprintf("source dataset %s not found", code))
		}
		sourceDatasetIDs = append(sourceDatasetIDs, sourceDataset.ID)
	}

	transformation := &domain.Transformation{
		Code:               req.Code,
		Name:               req.Name,
		Description:        req.Description,
		SourceDatasetCodes: req.SourceDatasetCodes,
		TargetDatasetCode:  req.TargetDatasetCode,
		AlgorithmCode:      req.AlgorithmCode,
		Algorithm:          algorithm,
		PeriodSeconds:      req.PeriodSeconds,
		Enabled:            req.Enabled,
	}

	err = u.trm.Do(ctx, func(ctx context.Context) error {
		return u.transformationRepo.Add(ctx, transformation, sourceDatasetIDs, targetDataset.ID, algorithm.ID)
	})
	if err != nil {
		return nil, err
	}

	created, err := u.transformationRepo.GetByCode(ctx, req.Code)
	if err != nil {
		return nil, err
	}
	return transformationResponse(created), nil
}

func (u *TransformationUsecase) GetTransformation(ctx context.Context, code string) (*dtos.TransformationResponse, error) {
	transformation, err := u.transformationRepo.GetByCode(ctx, code)
	if err != nil {
		return nil, err
	}
	if transformation == nil {
		return nil, domain.NewDomainErrorWithMessage(domain.ErrNotFoundCode, fmt.Sprintf("transformation %s not found", code))
	}
	return transformationResponse(transformation), nil
}

func (u *TransformationUsecase) ListTransformations(ctx context.Context) ([]dtos.TransformationResponse, error) {
	transformations, err := u.transformationRepo.List(ctx)
	if err != nil {
		return nil, err
	}
	response := make([]dtos.TransformationResponse, 0, len(transformations))
	for _, transformation := range transformations {
		item := transformation
		response = append(response, *transformationResponse(&item))
	}
	return response, nil
}

func (u *TransformationUsecase) GetTree(ctx context.Context) (*dtos.TreeResponse, error) {
	datasets, err := u.datasetRepo.List(ctx)
	if err != nil {
		return nil, err
	}
	transformations, err := u.transformationRepo.List(ctx)
	if err != nil {
		return nil, err
	}

	response := &dtos.TreeResponse{}
	for _, dataset := range datasets {
		response.Datasets = append(response.Datasets, dtos.TreeDataset{Code: dataset.Code, Name: dataset.Name})
	}
	for _, transformation := range transformations {
		response.Transformations = append(response.Transformations, dtos.TreeTransformation{Code: transformation.Code, Name: transformation.Name})
		for _, sourceDatasetCode := range transformation.SourceDatasetCodes {
			response.Edges = append(response.Edges, dtos.TreeEdge{
				From: "dataset:" + sourceDatasetCode,
				To:   "transformation:" + transformation.Code,
				Kind: "input",
			})
		}
		response.Edges = append(response.Edges, dtos.TreeEdge{
			From: "transformation:" + transformation.Code,
			To:   "dataset:" + transformation.TargetDatasetCode,
			Kind: "output",
		})
	}
	return response, nil
}

func (u *TransformationUsecase) RunDueTransformations(ctx context.Context) {
	transformations, err := u.transformationRepo.ListDue(ctx, 10)
	if err != nil {
		return
	}
	for _, transformation := range transformations {
		item := transformation
		go func() {
			_, _ = u.runUsecase.RunTransformation(context.Background(), item.Code)
		}()
	}
}
