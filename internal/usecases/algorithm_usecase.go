package usecases

import (
	"context"
	"fmt"

	"oiltech/internal/domain"
	"oiltech/internal/dtos"
	"oiltech/internal/adapters/postgres"
)

type AlgorithmUsecase struct {
	algorithmRepo domain.AlgorithmRepo
}

func NewAlgorithmUsecase(algorithmRepo domain.AlgorithmRepo) *AlgorithmUsecase {
	return &AlgorithmUsecase{algorithmRepo: algorithmRepo}
}

func (u *AlgorithmUsecase) CreateAlgorithm(ctx context.Context, req dtos.CreateAlgorithmRequest) (*dtos.AlgorithmResponse, error) {
	if err := postgres.ValidateDatasetCode(req.Code); err != nil {
		return nil, domain.NewDomainErrorWithMessage(domain.ErrValidationCode, err.Error())
	}
	if req.Name == "" {
		return nil, domain.NewDomainErrorWithMessage(domain.ErrValidationCode, "algorithm name is required")
	}
	if req.RunURL == "" {
		return nil, domain.NewDomainErrorWithMessage(domain.ErrValidationCode, "algorithm run_url is required")
	}

	existing, err := u.algorithmRepo.GetByCode(ctx, req.Code)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, domain.NewDomainErrorWithMessage(domain.ErrAlreadyExistsCode, fmt.Sprintf("algorithm %s already exists", req.Code))
	}

	algorithm := &domain.Algorithm{
		Code:        req.Code,
		Name:        req.Name,
		Description: req.Description,
		RunURL:      req.RunURL,
	}
	if err := u.algorithmRepo.Add(ctx, algorithm); err != nil {
		return nil, err
	}
	return algorithmResponse(algorithm), nil
}

func (u *AlgorithmUsecase) GetAlgorithm(ctx context.Context, code string) (*dtos.AlgorithmResponse, error) {
	algorithm, err := u.algorithmRepo.GetByCode(ctx, code)
	if err != nil {
		return nil, err
	}
	if algorithm == nil {
		return nil, domain.NewDomainErrorWithMessage(domain.ErrNotFoundCode, fmt.Sprintf("algorithm %s not found", code))
	}
	return algorithmResponse(algorithm), nil
}

func (u *AlgorithmUsecase) ListAlgorithms(ctx context.Context) ([]dtos.AlgorithmResponse, error) {
	algorithms, err := u.algorithmRepo.List(ctx)
	if err != nil {
		return nil, err
	}
	response := make([]dtos.AlgorithmResponse, 0, len(algorithms))
	for _, algorithm := range algorithms {
		item := algorithm
		response = append(response, *algorithmResponse(&item))
	}
	return response, nil
}
