package algokit

type RunRequest struct {
	RunID              string `json:"run_id"`
	TransformationCode string `json:"transformation_code"`
}

type RunResponse struct {
	Status       string `json:"status"`
	Message      string `json:"message,omitempty"`
	RowsRead     int64  `json:"rows_read,omitempty"`
	RowsWritten  int64  `json:"rows_written,omitempty"`
	ErrorMessage string `json:"error_message,omitempty"`
}
