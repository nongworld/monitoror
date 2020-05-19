//go:generate mockery -name Repository

package api

import "github.com/monitoror/monitoror/monitorables/gitlab/api/models"

type (
	Repository interface {
		GetPipeline(projectID, pipelineID int) (*models.Pipeline, error)
		GetPipelines(projectID int, ref string) ([]int, error)
		GetMergeRequest(projectID, mergeRequestID int) (*models.MergeRequest, error)
		GetMergeRequests(projectID int) ([]int, error)
		GetProject(projectID int) (*models.Project, error)
	}
)
