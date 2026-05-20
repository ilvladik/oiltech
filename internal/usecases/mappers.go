package usecases

import (
	"time"

	"oiltech/internal/domain"
	"oiltech/internal/dtos"
)

func datasetResponse(dataset *domain.Dataset) *dtos.DatasetResponse {
	columns := make([]dtos.DatasetColumn, 0, len(dataset.Columns))
	for _, column := range dataset.Columns {
		columns = append(columns, dtos.DatasetColumn{
			Code:     column.Code,
			Name:     column.Name,
			DataType: column.DataType,
			Ordinal:  column.Ordinal,
			Nullable: column.Nullable,
		})
	}
	return &dtos.DatasetResponse{
		ID:          dataset.ID,
		Code:        dataset.Code,
		Name:        dataset.Name,
		Description: dataset.Description,
		TableName:   dataset.TableName,
		MQTTTopic:   dataset.MQTTTopic,
		IsHidden:    dataset.IsHidden,
		Columns:     columns,
		CreatedAt:   dataset.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   dataset.UpdatedAt.Format(time.RFC3339),
	}
}

func algorithmResponse(algorithm *domain.Algorithm) *dtos.AlgorithmResponse {
	if algorithm == nil {
		return nil
	}
	return &dtos.AlgorithmResponse{
		ID:          algorithm.ID,
		Code:        algorithm.Code,
		Name:        algorithm.Name,
		Description: algorithm.Description,
		RunURL:      algorithm.RunURL,
		CreatedAt:   algorithm.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   algorithm.UpdatedAt.Format(time.RFC3339),
	}
}

func transformationResponse(transformation *domain.Transformation) *dtos.TransformationResponse {
	lastProcessedAt := ""
	if transformation.LastProcessedAt != nil {
		lastProcessedAt = transformation.LastProcessedAt.Format(time.RFC3339)
	}
	return &dtos.TransformationResponse{
		ID:                 transformation.ID,
		Code:               transformation.Code,
		Name:               transformation.Name,
		Description:        transformation.Description,
		SourceDatasetCodes: transformation.SourceDatasetCodes,
		TargetDatasetCode:  transformation.TargetDatasetCode,
		AlgorithmCode:      transformation.AlgorithmCode,
		Algorithm:          algorithmResponse(transformation.Algorithm),
		PeriodSeconds:      transformation.PeriodSeconds,
		Enabled:            transformation.Enabled,
		LastProcessedAt:    lastProcessedAt,
	}
}

func runResponse(run *domain.TransformationRun) *dtos.TransformationRunResponse {
	finishedAt := ""
	if run.FinishedAt != nil {
		finishedAt = run.FinishedAt.Format(time.RFC3339)
	}
	return &dtos.TransformationRunResponse{
		ID:                 run.ID,
		TransformationCode: run.TransformationCode,
		Status:             run.Status,
		StartedAt:          run.StartedAt.Format(time.RFC3339),
		FinishedAt:         finishedAt,
		ErrorMessage:       run.ErrorMessage,
		RequestBody:        run.RequestBody,
		ResponseBody:       run.ResponseBody,
	}
}
