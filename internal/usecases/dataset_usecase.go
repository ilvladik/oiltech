package usecases

import (
	"context"
	"fmt"
	"time"

	"oiltech/internal/domain"
	"oiltech/internal/dtos"
	"oiltech/internal/adapters/postgres"
)

type DatasetUsecase struct {
	datasetRepo domain.DatasetRepo
	trm         domain.TransactionManager
}

func NewDatasetUsecase(datasetRepo domain.DatasetRepo, trm domain.TransactionManager) *DatasetUsecase {
	return &DatasetUsecase{
		datasetRepo: datasetRepo,
		trm:         trm,
	}
}

func (u *DatasetUsecase) CreateDataset(ctx context.Context, req dtos.CreateDatasetRequest) (*dtos.DatasetResponse, error) {
	if err := postgres.ValidateDatasetCode(req.Code); err != nil {
		return nil, domain.NewDomainErrorWithMessage(domain.ErrValidationCode, err.Error())
	}
	if req.Name == "" {
		return nil, domain.NewDomainErrorWithMessage(domain.ErrValidationCode, "dataset name is required")
	}
	if len(req.Columns) == 0 {
		return nil, domain.NewDomainErrorWithMessage(domain.ErrValidationCode, "dataset columns are required")
	}

	existing, err := u.datasetRepo.GetByCode(ctx, req.Code)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, domain.NewDomainErrorWithMessage(domain.ErrAlreadyExistsCode, fmt.Sprintf("dataset %s already exists", req.Code))
	}

	now := time.Now().UTC()
	dataset := &domain.Dataset{
		Code:        req.Code,
		Name:        req.Name,
		Description: req.Description,
		TableName:   "ds_" + req.Code,
		MQTTTopic:   req.MQTTTopic,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	columns := make([]domain.DatasetColumn, 0, len(req.Columns))
	for i, column := range req.Columns {
		if _, err := postgres.QuoteIdent(column.Code); err != nil {
			return nil, domain.NewDomainErrorWithMessage(domain.ErrValidationCode, err.Error())
		}
		if _, err := postgres.PgType(column.DataType); err != nil {
			return nil, domain.NewDomainErrorWithMessage(domain.ErrValidationCode, err.Error())
		}
		columns = append(columns, domain.DatasetColumn{
			Code:     column.Code,
			Name:     column.Name,
			DataType: column.DataType,
			Ordinal:  i + 1,
			Nullable: column.Nullable,
		})
	}

	err = u.trm.Do(ctx, func(ctx context.Context) error {
		if err := u.datasetRepo.Add(ctx, dataset); err != nil {
			return err
		}
		for i := range columns {
			columns[i].DatasetID = dataset.ID
		}
		if err := u.datasetRepo.AddColumns(ctx, columns); err != nil {
			return err
		}
		if err := u.datasetRepo.CreatePhysicalTable(ctx, dataset, columns); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	created, err := u.datasetRepo.GetByCode(ctx, req.Code)
	if err != nil {
		return nil, err
	}
	return datasetResponse(created), nil
}

func (u *DatasetUsecase) ListDatasets(ctx context.Context) ([]dtos.DatasetResponse, error) {
	datasets, err := u.datasetRepo.List(ctx)
	if err != nil {
		return nil, err
	}
	response := make([]dtos.DatasetResponse, 0, len(datasets))
	for _, dataset := range datasets {
		item := dataset
		response = append(response, *datasetResponse(&item))
	}
	return response, nil
}

func (u *DatasetUsecase) GetDataset(ctx context.Context, code string) (*dtos.DatasetResponse, error) {
	dataset, err := u.datasetRepo.GetByCode(ctx, code)
	if err != nil {
		return nil, err
	}
	if dataset == nil {
		return nil, domain.NewDomainErrorWithMessage(domain.ErrNotFoundCode, fmt.Sprintf("dataset %s not found", code))
	}
	return datasetResponse(dataset), nil
}

func (u *DatasetUsecase) GetRows(ctx context.Context, code, filterField, filterValue string, limit int) (*dtos.RowsResponse, error) {
	dataset, err := u.datasetRepo.GetByCode(ctx, code)
	if err != nil {
		return nil, err
	}
	if dataset == nil {
		return nil, domain.NewDomainErrorWithMessage(domain.ErrNotFoundCode, fmt.Sprintf("dataset %s not found", code))
	}
	rows, err := u.datasetRepo.QueryRows(ctx, dataset, filterField, filterValue, limit)
	if err != nil {
		return nil, err
	}
	return &dtos.RowsResponse{
		DatasetCode: code,
		Rows:        rows,
	}, nil
}

func (u *DatasetUsecase) StoreMQTTMessage(ctx context.Context, datasetCode string, req dtos.MQTTDatasetMessage) error {
	dataset, err := u.datasetRepo.GetByCode(ctx, datasetCode)
	if err != nil {
		return err
	}
	if dataset == nil {
		return domain.NewDomainErrorWithMessage(domain.ErrNotFoundCode, fmt.Sprintf("dataset %s not found", datasetCode))
	}
	return u.datasetRepo.InsertMessage(ctx, dataset, domain.MQTTDatasetMessage{
		EventID: req.EventID,
		Values:  req.Values,
	})
}
