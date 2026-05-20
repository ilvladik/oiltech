package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"oiltech/internal/domain"

	"github.com/jmoiron/sqlx"
)

type SQLDatasetRepo struct {
	db *sqlx.DB
}

func NewDatasetRepo(db *sqlx.DB) *SQLDatasetRepo {
	return &SQLDatasetRepo{db: db}
}

func (r *SQLDatasetRepo) Add(ctx context.Context, dataset *domain.Dataset) error {
	query := `
		insert into meta.datasets(code, name, description, table_name, mqtt_topic)
		values ($1, $2, $3, $4, $5)
		returning id, created_at, updated_at
	`
	return TxOrDb(ctx, r.db).QueryRowContext(
		ctx,
		query,
		dataset.Code,
		dataset.Name,
		dataset.Description,
		dataset.TableName,
		nullableString(dataset.MQTTTopic),
	).Scan(&dataset.ID, &dataset.CreatedAt, &dataset.UpdatedAt)
}

func (r *SQLDatasetRepo) AddColumns(ctx context.Context, columns []domain.DatasetColumn) error {
	query := `
		insert into meta.dataset_columns(dataset_id, code, name, data_type, ordinal, nullable)
		values ($1, $2, $3, $4, $5, $6)
		returning id
	`
	for i := range columns {
		err := TxOrDb(ctx, r.db).QueryRowContext(
			ctx,
			query,
			columns[i].DatasetID,
			columns[i].Code,
			columns[i].Name,
			strings.ToLower(columns[i].DataType),
			columns[i].Ordinal,
			columns[i].Nullable,
		).Scan(&columns[i].ID)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *SQLDatasetRepo) CreatePhysicalTable(ctx context.Context, dataset *domain.Dataset, columns []domain.DatasetColumn) error {
	table, err := QuoteTable("data", dataset.TableName)
	if err != nil {
		return err
	}
	defs := []string{
		`"id" uuid primary key default gen_random_uuid()`,
		`"event_id" text unique`,
		`"received_at" timestamptz not null default now()`,
	}
	for _, column := range columns {
		identifier, err := QuoteIdent(column.Code)
		if err != nil {
			return err
		}
		pgType, err := PgType(column.DataType)
		if err != nil {
			return err
		}
		defs = append(defs, strings.TrimSpace(fmt.Sprintf("%s %s %s", identifier, pgType, Nullability(column.Nullable))))
	}
	query := fmt.Sprintf(`create table %s (%s)`, table, strings.Join(defs, ","))
	_, err = TxOrDb(ctx, r.db).ExecContext(ctx, query)
	return err
}

func (r *SQLDatasetRepo) GetByCode(ctx context.Context, code string) (*domain.Dataset, error) {
	query := `
		select id, code, name, description, table_name, coalesce(mqtt_topic,''), is_hidden, created_at, updated_at
		from meta.datasets
		where code = $1
	`
	dataset := &domain.Dataset{}
	err := r.db.QueryRowContext(ctx, query, code).Scan(
		&dataset.ID,
		&dataset.Code,
		&dataset.Name,
		&dataset.Description,
		&dataset.TableName,
		&dataset.MQTTTopic,
		&dataset.IsHidden,
		&dataset.CreatedAt,
		&dataset.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	columns, err := r.GetColumns(ctx, dataset.ID)
	if err != nil {
		return nil, err
	}
	dataset.Columns = columns
	return dataset, nil
}

func (r *SQLDatasetRepo) GetByID(ctx context.Context, id string) (*domain.Dataset, error) {
	query := `
		select id, code, name, description, table_name, coalesce(mqtt_topic,''), is_hidden, created_at, updated_at
		from meta.datasets
		where id = $1
	`
	dataset := &domain.Dataset{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&dataset.ID,
		&dataset.Code,
		&dataset.Name,
		&dataset.Description,
		&dataset.TableName,
		&dataset.MQTTTopic,
		&dataset.IsHidden,
		&dataset.CreatedAt,
		&dataset.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	columns, err := r.GetColumns(ctx, dataset.ID)
	if err != nil {
		return nil, err
	}
	dataset.Columns = columns
	return dataset, nil
}

func (r *SQLDatasetRepo) List(ctx context.Context) ([]domain.Dataset, error) {
	rows, err := r.db.QueryContext(ctx, `
		select id, code, name, description, table_name, coalesce(mqtt_topic,''), is_hidden, created_at, updated_at
		from meta.datasets
		order by code
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var datasets []domain.Dataset
	for rows.Next() {
		var dataset domain.Dataset
		if err := rows.Scan(
			&dataset.ID,
			&dataset.Code,
			&dataset.Name,
			&dataset.Description,
			&dataset.TableName,
			&dataset.MQTTTopic,
			&dataset.IsHidden,
			&dataset.CreatedAt,
			&dataset.UpdatedAt,
		); err != nil {
			return nil, err
		}
		datasets = append(datasets, dataset)
	}
	return datasets, rows.Err()
}

func (r *SQLDatasetRepo) GetColumns(ctx context.Context, datasetID string) ([]domain.DatasetColumn, error) {
	rows, err := r.db.QueryContext(ctx, `
		select id, dataset_id, code, name, data_type, ordinal, nullable
		from meta.dataset_columns
		where dataset_id = $1
		order by ordinal
	`, datasetID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var columns []domain.DatasetColumn
	for rows.Next() {
		var column domain.DatasetColumn
		if err := rows.Scan(
			&column.ID,
			&column.DatasetID,
			&column.Code,
			&column.Name,
			&column.DataType,
			&column.Ordinal,
			&column.Nullable,
		); err != nil {
			return nil, err
		}
		columns = append(columns, column)
	}
	return columns, rows.Err()
}

func nullableString(value string) any {
	if value == "" {
		return nil
	}
	return value
}
