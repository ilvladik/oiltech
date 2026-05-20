package postgres

import (
	"context"
	"database/sql"

	"oiltech/internal/domain"

	"github.com/jmoiron/sqlx"
)

type SQLAlgorithmRepo struct {
	db *sqlx.DB
}

func NewAlgorithmRepo(db *sqlx.DB) *SQLAlgorithmRepo {
	return &SQLAlgorithmRepo{db: db}
}

func (r *SQLAlgorithmRepo) Add(ctx context.Context, algorithm *domain.Algorithm) error {
	query := `
		insert into meta.algorithms(code, name, description, run_url)
		values ($1, $2, $3, $4)
		returning id, created_at, updated_at
	`
	return TxOrDb(ctx, r.db).QueryRowContext(
		ctx,
		query,
		algorithm.Code,
		algorithm.Name,
		algorithm.Description,
		algorithm.RunURL,
	).Scan(&algorithm.ID, &algorithm.CreatedAt, &algorithm.UpdatedAt)
}

func (r *SQLAlgorithmRepo) GetByCode(ctx context.Context, code string) (*domain.Algorithm, error) {
	algorithm := &domain.Algorithm{}
	err := r.db.QueryRowContext(ctx, `
		select id, code, name, description, run_url, created_at, updated_at
		from meta.algorithms
		where code = $1
	`, code).Scan(
		&algorithm.ID,
		&algorithm.Code,
		&algorithm.Name,
		&algorithm.Description,
		&algorithm.RunURL,
		&algorithm.CreatedAt,
		&algorithm.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return algorithm, nil
}

func (r *SQLAlgorithmRepo) List(ctx context.Context) ([]domain.Algorithm, error) {
	rows, err := r.db.QueryContext(ctx, `
		select id, code, name, description, run_url, created_at, updated_at
		from meta.algorithms
		order by code
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var algorithms []domain.Algorithm
	for rows.Next() {
		var algorithm domain.Algorithm
		if err := rows.Scan(
			&algorithm.ID,
			&algorithm.Code,
			&algorithm.Name,
			&algorithm.Description,
			&algorithm.RunURL,
			&algorithm.CreatedAt,
			&algorithm.UpdatedAt,
		); err != nil {
			return nil, err
		}
		algorithms = append(algorithms, algorithm)
	}
	return algorithms, rows.Err()
}
