package dtos

type DatasetColumn struct {
	Code     string `json:"code"`
	Name     string `json:"name"`
	DataType string `json:"data_type"`
	Ordinal  int    `json:"ordinal"`
	Nullable bool   `json:"nullable"`
}

type CreateDatasetRequest struct {
	Code        string          `json:"code"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	MQTTTopic   string          `json:"mqtt_topic"`
	Columns     []DatasetColumn `json:"columns"`
}

type DatasetResponse struct {
	ID          string          `json:"id"`
	Code        string          `json:"code"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	TableName   string          `json:"table_name"`
	MQTTTopic   string          `json:"mqtt_topic"`
	IsHidden    bool            `json:"is_hidden"`
	Columns     []DatasetColumn `json:"columns"`
	CreatedAt   string          `json:"created_at"`
	UpdatedAt   string          `json:"updated_at"`
}

type RowsResponse struct {
	DatasetCode string           `json:"dataset_code"`
	Rows        []map[string]any `json:"rows"`
}

type MQTTDatasetMessage struct {
	EventID string         `json:"event_id"`
	Values  map[string]any `json:"values"`
}
