package postgres

import (
	"context"
	"database/sql"

	"oiltech/internal/domain"

	"github.com/jmoiron/sqlx"
)

type SQLRunRepo struct {
	db *sqlx.DB
}

func NewRunRepo(db *sqlx.DB) *SQLRunRepo {
	return &SQLRunRepo{db: db}
}

func (r *SQLRunRepo) AddRunning(ctx context.Context, transformationID string) (*domain.TransformationRun, error) {
	run := &domain.TransformationRun{TransformationID: transformationID, Status: "running"}
	err := TxOrDb(ctx, r.db).QueryRowContext(ctx, `
		insert into meta.transformation_runs(transformation_id, status)
		values ($1, 'running')
		returning id, started_at
	`, transformationID).Scan(&run.ID, &run.StartedAt)
	if err != nil {
		return nil, err
	}
	return run, nil
}

func (r *SQLRunRepo) Finish(ctx context.Context, runID, status, errorMessage, requestBody, responseBody string) error {
	_, err := TxOrDb(ctx, r.db).ExecContext(ctx, `
		update meta.transformation_runs
		set status=$2, error_message=$3, request_body=$4, response_body=$5, finished_at=now()
		where id=$1
	`, runID, status, errorMessage, requestBody, responseBody)
	return err
}

func (r *SQLRunRepo) FinishSuccessAndTouchTransformation(ctx context.Context, runID, transformationID, requestBody, responseBody string) error {
	if err := r.Finish(ctx, runID, "success", "", requestBody, responseBody); err != nil {
		return err
	}
	_, err := TxOrDb(ctx, r.db).ExecContext(ctx, `
		update meta.transformations
		set last_processed_at=now(), updated_at=now()
		where id=$1
	`, transformationID)
	return err
}

func (r *SQLRunRepo) GetByID(ctx context.Context, runID string) (*domain.TransformationRun, error) {
	run := &domain.TransformationRun{}
	var finished sql.NullTime
	err := r.db.QueryRowContext(ctx, `
		select r.id, r.transformation_id, t.code, r.status, r.started_at, r.finished_at,
		       coalesce(r.error_message,''), coalesce(r.request_body,''), coalesce(r.response_body,'')
		from meta.transformation_runs r
		join meta.transformations t on t.id = r.transformation_id
		where r.id=$1
	`, runID).Scan(
		&run.ID,
		&run.TransformationID,
		&run.TransformationCode,
		&run.Status,
		&run.StartedAt,
		&finished,
		&run.ErrorMessage,
		&run.RequestBody,
		&run.ResponseBody,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if finished.Valid {
		run.FinishedAt = &finished.Time
	}
	return run, nil
}

func (r *SQLRunRepo) ListByTransformationCode(ctx context.Context, transformationCode string, limit int) ([]domain.TransformationRun, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	rows, err := r.db.QueryContext(ctx, `
		select r.id, r.transformation_id, t.code, r.status, r.started_at, r.finished_at,
		       coalesce(r.error_message,''), coalesce(r.request_body,''), coalesce(r.response_body,'')
		from meta.transformation_runs r
		join meta.transformations t on t.id = r.transformation_id
		where t.code=$1
		order by r.started_at desc
		limit $2
	`, transformationCode, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var runs []domain.TransformationRun
	for rows.Next() {
		var run domain.TransformationRun
		var finished sql.NullTime
		if err := rows.Scan(
			&run.ID,
			&run.TransformationID,
			&run.TransformationCode,
			&run.Status,
			&run.StartedAt,
			&finished,
			&run.ErrorMessage,
			&run.RequestBody,
			&run.ResponseBody,
		); err != nil {
			return nil, err
		}
		if finished.Valid {
			run.FinishedAt = &finished.Time
		}
		runs = append(runs, run)
	}
	return runs, rows.Err()
}

func (r *SQLRunRepo) HasRunningRun(ctx context.Context, transformationID string) (bool, error) {
	var exists bool
	err := r.db.QueryRowContext(ctx, `
		select exists(
			select 1 from meta.transformation_runs
			where transformation_id=$1 and status='running'
		)
	`, transformationID).Scan(&exists)
	return exists, err
}
