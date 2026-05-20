package domain

import "time"

type Transformation struct {
	ID                 string
	Code               string
	Name               string
	Description        string
	SourceDatasetCodes []string
	TargetDatasetCode  string
	AlgorithmCode      string
	Algorithm          *Algorithm
	PeriodSeconds      *int
	Enabled            bool
	LastProcessedAt    *time.Time
}

type TransformationRun struct {
	ID                 string
	TransformationID   string
	TransformationCode string
	Status             string
	StartedAt          time.Time
	FinishedAt         *time.Time
	ErrorMessage       string
	RequestBody        string
	ResponseBody       string
}

type TreeResponse struct {
	Datasets        []TreeDataset
	Transformations []TreeTransformation
	Edges           []TreeEdge
}

type TreeDataset struct {
	Code string
	Name string
}

type TreeTransformation struct {
	Code string
	Name string
}

type TreeEdge struct {
	From string
	To   string
	Kind string
}
