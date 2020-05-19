package repository

import (
	"fmt"
	"net/http"
	"time"

	"github.com/AlekSi/pointer"

	"github.com/monitoror/monitoror/monitorables/gitlab/api"
	"github.com/monitoror/monitoror/monitorables/gitlab/api/models"
	"github.com/monitoror/monitoror/monitorables/gitlab/config"
	"github.com/monitoror/monitoror/pkg/gogitlab"

	"github.com/xanzy/go-gitlab"
)

type (
	gitlabRepository struct {
		config *config.Gitlab

		pipelinesService     gogitlab.Pipelines
		mergeRequestsService gogitlab.MergeRequests
		projectService       gogitlab.Project
	}
)

func NewGitlabRepository(config *config.Gitlab) api.Repository {
	httpClient := &http.Client{
		Timeout: time.Duration(config.Timeout) * time.Millisecond,
	}
	gitlabApiBaseUrl := fmt.Sprintf("%s/api/v4", config.URL)

	git, err := gitlab.NewClient(config.Token, gitlab.WithBaseURL(gitlabApiBaseUrl), gitlab.WithHTTPClient(httpClient))
	if err != nil {
		// only when gitlabApiBaseUrl is not a valid URL
		panic(fmt.Sprintf("unable to setup Gitlab client\n. %v\n", err))
	}

	return &gitlabRepository{
		config: config,

		pipelinesService:     git.Pipelines,
		mergeRequestsService: git.MergeRequests,
		projectService:       git.Projects,
	}
}

func (gr *gitlabRepository) GetPipeline(projectID, pipelineID int) (*models.Pipeline, error) {
	gitlabPipeline, _, err := gr.pipelinesService.GetPipeline(projectID, pipelineID)
	if err != nil {
		return nil, err
	}

	pipeline := &models.Pipeline{
		ID:         gitlabPipeline.ID,
		Branch:     gitlabPipeline.Ref,
		Status:     gitlabPipeline.Status,
		StartedAt:  gitlabPipeline.CreatedAt,
		FinishedAt: gitlabPipeline.FinishedAt,
	}

	if gitlabPipeline.User != nil {
		pipeline.Author.Name = gitlabPipeline.User.Name
		pipeline.Author.AvatarURL = gitlabPipeline.User.AvatarURL

		if pipeline.Author.Name == "" {
			pipeline.Author.Name = gitlabPipeline.User.Username
		}
	}

	return pipeline, nil
}

func (gr *gitlabRepository) GetPipelines(projectID int, ref string) ([]int, error) {
	var ids []int

	gitlabPipelines, _, err := gr.pipelinesService.ListProjectPipelines(projectID, &gitlab.ListProjectPipelinesOptions{
		Ref:     &ref,
		OrderBy: pointer.ToString("id"),
		Sort:    pointer.ToString("desc"),
	})
	if err != nil {
		return nil, err
	}

	for _, pipeline := range gitlabPipelines {
		ids = append(ids, pipeline.ID)
	}

	return ids, nil
}

func (gr *gitlabRepository) GetMergeRequest(projectID, mergeRequestID int) (*models.MergeRequest, error) {
	gitlabMergeRequest, _, err := gr.mergeRequestsService.GetMergeRequest(projectID, mergeRequestID, &gitlab.GetMergeRequestsOptions{})
	if err != nil {
		return nil, err
	}

	mergeRequest := &models.MergeRequest{
		ID:        gitlabMergeRequest.IID,
		Title:     gitlabMergeRequest.Title,
		ProjectID: gitlabMergeRequest.SourceProjectID,
		Branch:    gitlabMergeRequest.SourceBranch,
		CommitSHA: gitlabMergeRequest.SHA,
	}

	if gitlabMergeRequest.Pipeline != nil {
		mergeRequest.PipelineID = &gitlabMergeRequest.Pipeline.ID
	}

	if gitlabMergeRequest.Author != nil {
		mergeRequest.Author.Name = gitlabMergeRequest.Author.Name
		mergeRequest.Author.AvatarURL = gitlabMergeRequest.Author.AvatarURL

		if mergeRequest.Author.Name == "" {
			mergeRequest.Author.Name = gitlabMergeRequest.Author.Username
		}
	}

	return mergeRequest, nil
}

func (gr *gitlabRepository) GetMergeRequests(projectID int) ([]int, error) {
	var ids []int

	gitlabMergeRequests, _, err := gr.mergeRequestsService.ListProjectMergeRequests(projectID, &gitlab.ListProjectMergeRequestsOptions{
		State: pointer.ToString("opened"),
	})
	if err != nil {
		return nil, err
	}

	for _, mergeRequest := range gitlabMergeRequests {
		ids = append(ids, mergeRequest.IID)
	}

	return ids, nil
}

func (gr *gitlabRepository) GetProject(projectID int) (*models.Project, error) {
	gitlabProject, _, err := gr.projectService.GetProject(projectID, &gitlab.GetProjectOptions{})
	if err != nil {
		return nil, err
	}

	project := &models.Project{
		ID:         gitlabProject.ID,
		Owner:      gitlabProject.Namespace.Path,
		Repository: gitlabProject.Path,
	}

	return project, nil
}
