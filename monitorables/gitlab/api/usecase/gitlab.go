package usecase

import (
	"fmt"
	"time"

	"github.com/AlekSi/pointer"
	uuid "github.com/satori/go.uuid"

	uiConfigModels "github.com/monitoror/monitoror/api/config/models"
	monitorableCache "github.com/monitoror/monitoror/internal/pkg/monitorable/cache"
	coreModels "github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorables/gitlab/api"
	"github.com/monitoror/monitoror/monitorables/gitlab/api/models"
	"github.com/monitoror/monitoror/pkg/git"

	"github.com/jsdidierlaurent/echo-middleware/cache"
)

type (
	gitlabUsecase struct {
		repository api.Repository
		// Used to generate store key by repository
		repositoryUID string

		// store is used to store persistent data (project,
		store cache.Store

		// builds cache. used for save small history of build for stats
		buildsCache *monitorableCache.BuildCache
	}
)

const (
	buildCacheSize = 5

	GitlabProjectStoreKeyPrefix = "monitoror.gitlab.project.store"
)

func NewGitlabUsecase(repository api.Repository, store cache.Store) api.Usecase {
	return &gitlabUsecase{
		repository:    repository,
		repositoryUID: uuid.NewV4().String(),
		store:         store,
		buildsCache:   monitorableCache.NewBuildCache(buildCacheSize),
	}
}

func (gu *gitlabUsecase) getProject(projectId int) (*models.Project, error) {
	project := &models.Project{}

	storeKey := gu.getProjectStoreKey(projectId)
	if err := gu.store.Get(storeKey, project); err != nil {
		if project, err = gu.repository.GetProject(projectId); err != nil {
			return nil, err
		}

		_ = gu.store.Set(storeKey, *project, cache.NEVER)
	}

	return project, nil
}

func (gu *gitlabUsecase) Pipeline(params *models.PipelineParams) (*coreModels.Tile, error) {
	tile := coreModels.NewTile(api.GitlabPipelineTileType).WithBuild()
	tile.Label = fmt.Sprintf("%d", params.ProjectID)
	tile.Build.Branch = pointer.ToString(git.HumanizeBranch(params.Ref))

	// Load Project and cache it
	project, err := gu.getProject(*params.ProjectID)
	if err != nil {
		return nil, &coreModels.MonitororError{Err: err, Tile: tile, Message: "unable to load project"}
	}
	tile.Label = project.Repository

	// Load pipelines for given ref
	pipelines, err := gu.repository.GetPipelines(*params.ProjectID, params.Ref)
	if err != nil {
		return nil, &coreModels.MonitororError{Err: err, Tile: tile, Message: "unable to load pipelines"}
	}
	if len(pipelines) == 0 {
		// Warning because request was correct but there is no build
		return nil, &coreModels.MonitororError{Tile: tile, Message: "no pipelines found", ErrorStatus: coreModels.UnknownStatus}
	}

	// Load pipeline detail
	pipeline, err := gu.repository.GetPipeline(*params.ProjectID, pipelines[0])
	if err != nil {
		return nil, &coreModels.MonitororError{Err: err, Tile: tile, Message: "unable to load pipeline"}
	}

	gu.computePipeline(params, tile, pipeline)

	// Author
	if tile.Status == coreModels.FailedStatus {
		tile.Build.Author = &pipeline.Author
	}

	return tile, nil
}

func (gu *gitlabUsecase) MergeRequest(params *models.MergeRequestParams) (*coreModels.Tile, error) {
	tile := coreModels.NewTile(api.GitlabMergeRequestTileType).WithBuild()
	tile.Label = fmt.Sprintf("%d", params.ProjectID)

	// Load Project and cache it
	project, err := gu.getProject(*params.ProjectID)
	if err != nil {
		return nil, &coreModels.MonitororError{Err: err, Tile: tile, Message: "unable to load project"}
	}
	tile.Label = project.Repository

	// Load MergeRequest
	mergeRequest, err := gu.repository.GetMergeRequest(*params.ProjectID, *params.ID)
	if err != nil {
		return nil, &coreModels.MonitororError{Err: err, Tile: tile, Message: "unable to load merge request"}
	}

	// Load MergeRequest project
	mergeRequestProject, err := gu.getProject(*params.ProjectID)
	if err != nil {
		return nil, &coreModels.MonitororError{Err: err, Tile: tile, Message: "unable to load project"}
	}

	tile.Build.Branch = pointer.ToString(git.HumanizeBranch(mergeRequest.Branch))
	if project.Owner != mergeRequestProject.Owner {
		tile.Build.Branch = pointer.ToString(fmt.Sprintf("%s:%s", mergeRequestProject.Owner, *tile.Build.Branch))
	}
	tile.Build.MergeRequest = &coreModels.TileMergeRequest{
		ID:    mergeRequest.ID,
		Title: mergeRequest.Title,
	}

	// Load Pipeline
	var pipeline *models.Pipeline
	if mergeRequest.PipelineID != nil {
		// Load pipeline detail
		pipeline, err = gu.repository.GetPipeline(*params.ProjectID, *mergeRequest.PipelineID)
		if err != nil {
			return nil, &coreModels.MonitororError{Err: err, Tile: tile, Message: "unable to load pipelines"}
		}
	} else if project.Owner != mergeRequestProject.Owner {
		// Load pipelines for given ref in case of fork
		pipelines, err := gu.repository.GetPipelines(*params.ProjectID, mergeRequest.Branch)
		if err != nil {
			return nil, &coreModels.MonitororError{Err: err, Tile: tile, Message: "unable to load pipelines"}
		}
		if len(pipelines) == 0 {
			// Warning because request was correct but there is no build
			return nil, &coreModels.MonitororError{Tile: tile, Message: "no pipelines found", ErrorStatus: coreModels.UnknownStatus}
		}

		// Load pipeline detail
		pipeline, err = gu.repository.GetPipeline(*params.ProjectID, pipelines[0])
		if err != nil {
			return nil, &coreModels.MonitororError{Err: err, Tile: tile, Message: "unable to load pipeline"}
		}
	} else {
		return nil, &coreModels.MonitororError{Tile: tile, Message: "no pipelines found", ErrorStatus: coreModels.UnknownStatus}
	}

	gu.computePipeline(params, tile, pipeline)

	// Author
	if tile.Status == coreModels.FailedStatus {
		tile.Build.Author = &mergeRequest.Author
	}

	return tile, nil
}

func (gu *gitlabUsecase) computePipeline(params interface{}, tile *coreModels.Tile, pipeline *models.Pipeline) {
	tile.Status = parseStatus(pipeline.Status)

	// Set Previous Status
	strPipelineID := fmt.Sprintf("%d", pipeline.ID)
	previousStatus := gu.buildsCache.GetPreviousStatus(params, strPipelineID)
	if previousStatus != nil {
		tile.Build.PreviousStatus = *previousStatus
	} else {
		tile.Build.PreviousStatus = coreModels.UnknownStatus
	}

	// StartedAt / FinishedAt
	tile.Build.StartedAt = pipeline.StartedAt
	if tile.Status != coreModels.RunningStatus && tile.Status != coreModels.QueuedStatus {
		tile.Build.FinishedAt = pipeline.FinishedAt
	}

	// Duration
	if tile.Status == coreModels.RunningStatus {
		tile.Build.Duration = pointer.ToInt64(int64(time.Since(*tile.Build.StartedAt).Seconds()))

		estimatedDuration := gu.buildsCache.GetEstimatedDuration(params)
		if estimatedDuration != nil {
			tile.Build.EstimatedDuration = pointer.ToInt64(int64(estimatedDuration.Seconds()))
		} else {
			tile.Build.EstimatedDuration = pointer.ToInt64(int64(0))
		}
	}

	// Cache Duration when success / failed
	if tile.Status == coreModels.SuccessStatus || tile.Status == coreModels.FailedStatus {
		gu.buildsCache.Add(params, strPipelineID, tile.Status, tile.Build.FinishedAt.Sub(*tile.Build.StartedAt))
	}
}

func (gu *gitlabUsecase) MergeRequestsGenerator(params interface{}) ([]uiConfigModels.GeneratedTile, error) {
	panic("implement me")
}

func (gu *gitlabUsecase) getProjectStoreKey(projectId int) string {
	return fmt.Sprintf("%s:%s-%d", GitlabProjectStoreKeyPrefix, gu.repositoryUID, projectId)
}

func parseStatus(status string) coreModels.TileStatus {
	// See: https://docs.gitlab.com/ee/api/pipelines.html#list-project-pipelines
	switch status {
	case "running":
		return coreModels.RunningStatus
	case "pending":
		return coreModels.QueuedStatus
	case "success":
		return coreModels.SuccessStatus
	case "failed":
		return coreModels.FailedStatus
	case "canceled":
		return coreModels.CanceledStatus
	case "skipped":
		return coreModels.CanceledStatus
	case "created":
		return coreModels.QueuedStatus
	case "manual":
		return coreModels.ActionRequiredStatus
	default:
		return coreModels.UnknownStatus
	}
}
