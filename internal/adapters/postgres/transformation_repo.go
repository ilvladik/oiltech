package postgres

import (
	"context"
	"database/sql"

	"oiltech/internal/domain"

	"github.com/jmoiron/sqlx"
)

type SQLTransformationRepo struct {
	db *sqlx.DB
}

func NewTransformationRepo(db *sqlx.DB) *SQLTransformationRepo {
	return &SQLTransformationRepo{db: db}
}

func (r *SQLTransformationRepo) Add(ctx context.Context, transformation *domain.Transformation, sourceDatasetIDs []string, targetDatasetID string, algorithmID string) error {
	var periodArg any
	if transformation.PeriodSeconds != nil {
		periodArg = *transformation.PeriodSeconds
	}
	query := `
		insert into meta.transformations(code, name, description, algorithm_id, target_dataset_id, period_seconds, enabled)
		values ($1, $2, $3, $4, $5, $6, $7)
		returning id
	`
	if err := TxOrDb(ctx, r.db).QueryRowContext(
		ctx,
		query,
		transformation.Code,
		transformation.Name,
		transformation.Description,
		algorithmID,
		targetDatasetID,
		periodArg,
		transformation.Enabled,
	).Scan(&transformation.ID); err != nil {
		return err
	}

	for _, sourceDatasetID := range sourceDatasetIDs {
		_, err := TxOrDb(ctx, r.db).ExecContext(ctx, `
			insert into meta.transformation_sources(transformation_id, dataset_id)
			values ($1, $2)
		`, transformation.ID, sourceDatasetID)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *SQLTransformationRepo) GetByCode(ctx context.Context, code string) (*domain.Transformation, error) {
	transformation, err := r.getByCode(ctx, code)
	if err != nil || transformation == nil {
		return transformation, err
	}
	sources, err := r.sources(ctx, transformation.ID)
	if err != nil {
		return nil, err
	}
	transformation.SourceDatasetCodes = sources
	return transformation, nil
}

func (r *SQLTransformationRepo) List(ctx context.Context) ([]domain.Transformation, error) {
	rows, err := r.db.QueryContext(ctx, `select code from meta.transformations order by code`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transformations []domain.Transformation
	for rows.Next() {
		var code string
		if err := rows.Scan(&code); err != nil {
			return nil, err
		}
		transformation, err := r.GetByCode(ctx, code)
		if err != nil {
			return nil, err
		}
		transformations = append(transformations, *transformation)
	}
	return transformations, rows.Err()
}

func (r *SQLTransformationRepo) ListDue(ctx context.Context, limit int) ([]domain.Transformation, error) {
	rows, err := r.db.QueryContext(ctx, `
		select t.code
		from meta.transformations t
		where t.enabled = true
		  and t.period_seconds is not null
		  and (
		    t.last_processed_at is null
		    or t.last_processed_at + (t.period_seconds || ' seconds')::interval <= now()
		  )
		  and not exists (
		    select 1
		    from meta.transformation_runs tr
		    where tr.transformation_id = t.id and tr.status = 'running'
		  )
		order by coalesce(t.last_processed_at, '-infinity'::timestamptz)
		limit $1
	`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transformations []domain.Transformation
	for rows.Next() {
		var code string
		if err := rows.Scan(&code); err != nil {
			return nil, err
		}
		transformation, err := r.GetByCode(ctx, code)
		if err != nil {
			return nil, err
		}
		transformations = append(transformations, *transformation)
	}
	return transformations, rows.Err()
}

func (r *SQLTransformationRepo) getByCode(ctx context.Context, code string) (*domain.Transformation, error) {
	transformation := &domain.Transformation{}
	algorithm := &domain.Algorithm{}
	var period sql.NullInt64
	var lastProcessed sql.NullTime
	err := r.db.QueryRowContext(ctx, `
		select t.id, t.code, t.name, t.description, d.code,
		       a.id, a.code, a.name, a.description, a.run_url, a.created_at, a.updated_at,
		       t.period_seconds, t.enabled, t.last_processed_at
		from meta.transformations t
		join meta.datasets d on d.id = t.target_dataset_id
		join meta.algorithms a on a.id = t.algorithm_id
		where t.code = $1
	`, code).Scan(
		&transformation.ID,
		&transformation.Code,
		&transformation.Name,
		&transformation.Description,
		&transformation.TargetDatasetCode,
		&algorithm.ID,
		&algorithm.Code,
		&algorithm.Name,
		&algorithm.Description,
		&algorithm.RunURL,
		&algorithm.CreatedAt,
		&algorithm.UpdatedAt,
		&period,
		&transformation.Enabled,
		&lastProcessed,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	transformation.Algorithm = algorithm
	transformation.AlgorithmCode = algorithm.Code
	if period.Valid {
		value := int(period.Int64)
		transformation.PeriodSeconds = &value
	}
	if lastProcessed.Valid {
		transformation.LastProcessedAt = &lastProcessed.Time
	}
	return transformation, nil
}

func (r *SQLTransformationRepo) sources(ctx context.Context, transformationID string) ([]string, error) {
	rows, err := r.db.QueryContext(ctx, `
		select d.code
		from meta.transformation_sources ts
		join meta.datasets d on d.id = ts.dataset_id
		where ts.transformation_id = $1
		order by d.code
	`, transformationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sources []string
	for rows.Next() {
		var code string
		if err := rows.Scan(&code); err != nil {
			return nil, err
		}
		sources = append(sources, code)
	}
	return sources, rows.Err()
}
