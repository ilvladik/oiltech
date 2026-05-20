package domain

import "context"

type DatasetRepo interface {
	Add(ctx context.Context, dataset *Dataset) error
	AddColumns(ctx context.Context, columns []DatasetColumn) error
	CreatePhysicalTable(ctx context.Context, dataset *Dataset, columns []DatasetColumn) error
	GetByCode(ctx context.Context, code string) (*Dataset, error)
	GetByID(ctx context.Context, id string) (*Dataset, error)
	List(ctx context.Context) ([]Dataset, error)
	GetColumns(ctx context.Context, datasetID string) ([]DatasetColumn, error)
	InsertMessage(ctx context.Context, dataset *Dataset, msg MQTTDatasetMessage) error
	QueryRows(ctx context.Context, dataset *Dataset, filterField, filterValue string, limit int) ([]map[string]any, error)
}

type AlgorithmRepo interface {
	Add(ctx context.Context, algorithm *Algorithm) error
	GetByCode(ctx context.Context, code string) (*Algorithm, error)
	List(ctx context.Context) ([]Algorithm, error)
}

type TransformationRepo interface {
	Add(ctx context.Context, transformation *Transformation, sourceDatasetIDs []string, targetDatasetID string, algorithmID string) error
	GetByCode(ctx context.Context, code string) (*Transformation, error)
	List(ctx context.Context) ([]Transformation, error)
	ListDue(ctx context.Context, limit int) ([]Transformation, error)
}

type RunRepo interface {
	AddRunning(ctx context.Context, transformationID string) (*TransformationRun, error)
	Finish(ctx context.Context, runID, status, errorMessage, requestBody, responseBody string) error
	FinishSuccessAndTouchTransformation(ctx context.Context, runID, transformationID, requestBody, responseBody string) error
	GetByID(ctx context.Context, runID string) (*TransformationRun, error)
	ListByTransformationCode(ctx context.Context, transformationCode string, limit int) ([]TransformationRun, error)
	HasRunningRun(ctx context.Context, transformationID string) (bool, error)
}

type TransactionManager interface {
	Do(ctx context.Context, fn func(context.Context) error) error
}
