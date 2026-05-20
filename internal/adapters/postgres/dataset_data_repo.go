package postgres

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"

	"oiltech/internal/domain"
)

func (r *SQLDatasetRepo) InsertMessage(ctx context.Context, dataset *domain.Dataset, msg domain.MQTTDatasetMessage) error {
	table, err := QuoteTable("data", dataset.TableName)
	if err != nil {
		return err
	}

	allowed := make(map[string]domain.DatasetColumn)
	for _, column := range dataset.Columns {
		allowed[column.Code] = column
	}

	keys := make([]string, 0, len(msg.Values))
	for key := range msg.Values {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	columns := []string{"event_id"}
	values := []any{eventID(msg)}
	for _, key := range keys {
		column, ok := allowed[key]
		if !ok {
			return domain.NewDomainErrorWithMessage(domain.ErrValidationCode, fmt.Sprintf("field %s is not declared in dataset %s", key, dataset.Code))
		}
		value, err := coerceValue(column.DataType, msg.Values[key])
		if err != nil {
			return err
		}
		columns = append(columns, key)
		values = append(values, value)
	}

	quotedColumns := make([]string, 0, len(columns))
	placeholders := make([]string, 0, len(columns))
	for i, column := range columns {
		quoted, err := QuoteIdent(column)
		if err != nil {
			return err
		}
		quotedColumns = append(quotedColumns, quoted)
		placeholders = append(placeholders, fmt.Sprintf("$%d", i+1))
	}

	query := fmt.Sprintf(
		`insert into %s (%s) values (%s) on conflict ("event_id") do nothing`,
		table,
		strings.Join(quotedColumns, ","),
		strings.Join(placeholders, ","),
	)
	_, err = r.db.ExecContext(ctx, query, values...)
	return err
}

func (r *SQLDatasetRepo) QueryRows(ctx context.Context, dataset *domain.Dataset, filterField, filterValue string, limit int) ([]map[string]any, error) {
	if limit <= 0 || limit > 500 {
		limit = 100
	}

	table, err := QuoteTable("data", dataset.TableName)
	if err != nil {
		return nil, err
	}

	query := fmt.Sprintf(`select * from %s`, table)
	args := []any{limit}
	if filterField != "" {
		if !datasetHasColumn(dataset, filterField) && filterField != "id" && filterField != "event_id" && filterField != "received_at" {
			return nil, domain.NewDomainErrorWithMessage(domain.ErrValidationCode, fmt.Sprintf("unknown filter field: %s", filterField))
		}
		quotedField, err := QuoteIdent(filterField)
		if err != nil {
			return nil, err
		}
		query += fmt.Sprintf(` where %s::text = $1 order by "received_at" desc limit $2`, quotedField)
		args = []any{filterValue, limit}
	} else {
		query += ` order by "received_at" desc limit $1`
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	columnNames, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	var result []map[string]any
	for rows.Next() {
		values := make([]any, len(columnNames))
		valuePointers := make([]any, len(columnNames))
		for i := range values {
			valuePointers[i] = &values[i]
		}
		if err := rows.Scan(valuePointers...); err != nil {
			return nil, err
		}
		item := map[string]any{}
		for i, name := range columnNames {
			item[name] = normalizeDBValue(values[i])
		}
		result = append(result, item)
	}
	return result, rows.Err()
}

func eventID(msg domain.MQTTDatasetMessage) string {
	if msg.EventID != "" {
		return msg.EventID
	}
	sum := sha256.Sum256([]byte(fmt.Sprintf("%v", msg.Values)))
	return hex.EncodeToString(sum[:])
}

func datasetHasColumn(dataset *domain.Dataset, field string) bool {
	for _, column := range dataset.Columns {
		if column.Code == field {
			return true
		}
	}
	return false
}

func coerceValue(dataType string, value any) (any, error) {
	switch strings.ToLower(dataType) {
	case "int", "long", "int64":
		switch v := value.(type) {
		case float64:
			return int64(v), nil
		case int:
			return int64(v), nil
		case int64:
			return v, nil
		default:
			return v, nil
		}
	case "string", "text":
		return fmt.Sprintf("%v", value), nil
	default:
		return value, nil
	}
}

func normalizeDBValue(value any) any {
	switch v := value.(type) {
	case []byte:
		return string(v)
	default:
		return v
	}
}
