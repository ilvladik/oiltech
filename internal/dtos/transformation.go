package dtos

type CreateTransformationRequest struct {
	Code               string   `json:"code"`
	Name               string   `json:"name"`
	Description        string   `json:"description"`
	SourceDatasetCodes []string `json:"source_dataset_codes"`
	TargetDatasetCode  string   `json:"target_dataset_code"`
	AlgorithmCode      string   `json:"algorithm_code"`
	PeriodSeconds      *int     `json:"period_seconds"`
	Enabled            bool     `json:"enabled"`
}

type TransformationResponse struct {
	ID                 string             `json:"id"`
	Code               string             `json:"code"`
	Name               string             `json:"name"`
	Description        string             `json:"description"`
	SourceDatasetCodes []string           `json:"source_dataset_codes"`
	TargetDatasetCode  string             `json:"target_dataset_code"`
	AlgorithmCode      string             `json:"algorithm_code"`
	Algorithm          *AlgorithmResponse `json:"algorithm,omitempty"`
	PeriodSeconds      *int               `json:"period_seconds"`
	Enabled            bool               `json:"enabled"`
	LastProcessedAt    string             `json:"last_processed_at,omitempty"`
}

type TransformationRunResponse struct {
	ID                 string     `json:"id"`
	TransformationCode string     `json:"transformation_code"`
	Status             string     `json:"status"`
	StartedAt          string     `json:"started_at"`
	FinishedAt         string     `json:"finished_at,omitempty"`
	ErrorMessage       string     `json:"error_message"`
	RequestBody        string     `json:"request_body"`
	ResponseBody       string     `json:"response_body"`
}

type TreeResponse struct {
	Datasets        []TreeDataset        `json:"datasets"`
	Transformations []TreeTransformation `json:"transformations"`
	Edges           []TreeEdge           `json:"edges"`
}

type TreeDataset struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

type TreeTransformation struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

type TreeEdge struct {
	From string `json:"from"`
	To   string `json:"to"`
	Kind string `json:"kind"`
}
