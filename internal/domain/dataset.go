package domain

import "time"

type Dataset struct {
	ID          string
	Code        string
	Name        string
	Description string
	TableName   string
	MQTTTopic   string
	IsHidden    bool
	Columns     []DatasetColumn
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type DatasetColumn struct {
	ID        string
	DatasetID string
	Code      string
	Name      string
	DataType  string
	Ordinal   int
	Nullable  bool
}

type MQTTDatasetMessage struct {
	EventID string
	Values  map[string]any
}
