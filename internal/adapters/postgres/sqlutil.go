package postgres

import (
	"fmt"
	"regexp"
	"strings"
)

var datasetCodeRe = regexp.MustCompile(`^[a-z][a-z0-9_]{2,62}$`)
var columnCodeRe = regexp.MustCompile(`^[A-Za-z][A-Za-z0-9_]{0,62}$`)

func ValidateDatasetCode(code string) error {
	if !datasetCodeRe.MatchString(code) {
		return fmt.Errorf("dataset code must match %s", datasetCodeRe.String())
	}
	return nil
}

func QuoteIdent(identifier string) (string, error) {
	if !columnCodeRe.MatchString(identifier) && !datasetCodeRe.MatchString(identifier) {
		return "", fmt.Errorf("unsafe sql identifier: %s", identifier)
	}
	return `"` + strings.ReplaceAll(identifier, `"`, `""`) + `"`, nil
}

func QuoteTable(schema, table string) (string, error) {
	quotedSchema, err := QuoteIdent(schema)
	if err != nil {
		return "", err
	}
	quotedTable, err := QuoteIdent(table)
	if err != nil {
		return "", err
	}
	return quotedSchema + "." + quotedTable, nil
}

func PgType(dataType string) (string, error) {
	switch strings.ToLower(dataType) {
	case "string", "text":
		return "text", nil
	case "int", "long", "int64":
		return "bigint", nil
	case "decimal", "numeric":
		return "numeric", nil
	case "float", "float64", "double":
		return "double precision", nil
	case "datetime", "timestamp":
		return "timestamptz", nil
	case "bool", "boolean":
		return "boolean", nil
	default:
		return "", fmt.Errorf("unsupported data type: %s", dataType)
	}
}

func Nullability(nullable bool) string {
	if nullable {
		return ""
	}
	return "not null"
}
