package usecases

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"oiltech/internal/algokit"
	"oiltech/internal/domain"
	"oiltech/internal/dtos"
)

type RunUsecase struct {
	transformationRepo domain.TransformationRepo
	runRepo            domain.RunRepo
	trm                domain.TransactionManager
	httpClient         *http.Client
}

func NewRunUsecase(
	transformationRepo domain.TransformationRepo,
	runRepo domain.RunRepo,
	trm domain.TransactionManager,
) *RunUsecase {
	return &RunUsecase{
		transformationRepo: transformationRepo,
		runRepo:            runRepo,
		trm:                trm,
		httpClient:         &http.Client{Timeout: 30 * time.Second},
	}
}

func (u *RunUsecase) RunTransformation(ctx context.Context, transformationCode string) (*dtos.TransformationRunResponse, error) {
	transformation, err := u.transformationRepo.GetByCode(ctx, transformationCode)
	if err != nil {
		return nil, err
	}
	if transformation == nil {
		return nil, domain.NewDomainErrorWithMessage(domain.ErrNotFoundCode, fmt.Sprintf("transformation %s not found", transformationCode))
	}
	if transformation.Algorithm == nil {
		return nil, domain.NewDomainErrorWithMessage(domain.ErrValidationCode, "transformation has no algorithm")
	}

	var run *domain.TransformationRun
	err = u.trm.Do(ctx, func(ctx context.Context) error {
		hasRunningRun, err := u.runRepo.HasRunningRun(ctx, transformation.ID)
		if err != nil {
			return err
		}
		if hasRunningRun {
			return domain.NewDomainErrorWithMessage(domain.ErrValidationCode, "transformation already has running run")
		}
		run, err = u.runRepo.AddRunning(ctx, transformation.ID)
		return err
	})
	if err != nil {
		return nil, err
	}

	requestBody, err := json.Marshal(algokit.RunRequest{
		RunID:              run.ID,
		TransformationCode: transformation.Code,
	})
	if err != nil {
		return nil, err
	}

	status := "success"
	errorMessage := ""
	responseBody := ""
	response, err := u.httpClient.Post(transformation.Algorithm.RunURL, "application/json", bytes.NewReader(requestBody))
	if err != nil {
		status = "failed"
		errorMessage = err.Error()
	} else {
		defer response.Body.Close()
		responseBytes, _ := io.ReadAll(io.LimitReader(response.Body, 1_000_000))
		responseBody = string(responseBytes)
		if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
			status = "failed"
			errorMessage = fmt.Sprintf("algorithm returned status %d", response.StatusCode)
		}
	}

	err = u.trm.Do(ctx, func(ctx context.Context) error {
		if status == "success" {
			return u.runRepo.FinishSuccessAndTouchTransformation(ctx, run.ID, transformation.ID, string(requestBody), responseBody)
		}
		return u.runRepo.Finish(ctx, run.ID, status, errorMessage, string(requestBody), responseBody)
	})
	if err != nil {
		return nil, err
	}

	updatedRun, err := u.runRepo.GetByID(ctx, run.ID)
	if err != nil {
		return nil, err
	}
	return runResponse(updatedRun), nil
}

func (u *RunUsecase) ListRuns(ctx context.Context, transformationCode string, limit int) ([]dtos.TransformationRunResponse, error) {
	runs, err := u.runRepo.ListByTransformationCode(ctx, transformationCode, limit)
	if err != nil {
		return nil, err
	}
	response := make([]dtos.TransformationRunResponse, 0, len(runs))
	for _, run := range runs {
		item := run
		response = append(response, *runResponse(&item))
	}
	return response, nil
}
