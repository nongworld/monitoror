package models

import (
	coreModels "github.com/monitoror/monitoror/models"
)

type MergeRequest struct {
	ID     int
	Title  string
	Author coreModels.Author

	PipelineID *int

	ProjectID int
	Branch    string
	CommitSHA string
}
