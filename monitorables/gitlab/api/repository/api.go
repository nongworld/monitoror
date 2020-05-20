package repository

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/monitoror/monitoror/monitorables/gitlab/api"
	"github.com/monitoror/monitoror/monitorables/gitlab/api/models"
	"github.com/monitoror/monitoror/monitorables/gitlab/config"
	"github.com/monitoror/monitoror/pkg/gogitlab"

	"github.com/AlekSi/pointer"
	"github.com/xanzy/go-gitlab"
)

type (
	gitlabRepository struct {
		config *config.Gitlab

		issuesService        gogitlab.IssuesService
		pipelinesService     gogitlab.PipelinesService
		mergeRequestsService gogitlab.MergeRequestsService
		projectService       gogitlab.ProjectService
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

		issuesService:        git.Issues,
		pipelinesService:     git.Pipelines,
		mergeRequestsService: git.MergeRequests,
		projectService:       git.Projects,
	}
}

func (gr *gitlabRepository) GetIssues(projectID *int, query string) (int, error) {
	values, err := url.ParseQuery(query)
	if err != nil {
		return 0, err
	}

	var resp *gitlab.Response
	if projectID != nil {
		listProjectIssueOption := &gitlab.ListProjectIssuesOptions{}
		listProjectIssueOption.State = valueToString(values, "state")
		listProjectIssueOption.Labels = valueToLabels(values, "labels")
		listProjectIssueOption.Milestone = valueToString(values, "milestone")
		listProjectIssueOption.Scope = valueToString(values, "scope")
		listProjectIssueOption.MyReactionEmoji = valueToString(values, "my_reaction_emoji")
		listProjectIssueOption.Search = valueToString(values, "search")
		listProjectIssueOption.AuthorID, err = valueToInt(values, "author_id")
		if err != nil {
			return 0, err
		}
		listProjectIssueOption.AssigneeID, err = valueToInt(values, "assignee_id")
		if err != nil {
			return 0, err
		}

		_, resp, err = gr.issuesService.ListProjectIssues(*projectID, listProjectIssueOption)
		if err != nil {
			return 0, err
		}
	} else {
		listIssueOption := &gitlab.ListIssuesOptions{}
		listIssueOption.State = valueToString(values, "state")
		listIssueOption.Labels = valueToLabels(values, "labels")
		listIssueOption.Milestone = valueToString(values, "milestone")
		listIssueOption.Scope = valueToString(values, "scope")
		listIssueOption.MyReactionEmoji = valueToString(values, "my_reaction_emoji")
		listIssueOption.Search = valueToString(values, "search")
		listIssueOption.AuthorID, err = valueToInt(values, "author_id")
		if err != nil {
			return 0, err
		}
		listIssueOption.AssigneeID, err = valueToInt(values, "assignee_id")
		if err != nil {
			return 0, err
		}

		_, resp, err = gr.issuesService.ListIssues(listIssueOption)
		if err != nil {
			return 0, err
		}
	}
	return resp.TotalItems, nil
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
		ID:              gitlabMergeRequest.IID,
		Title:           gitlabMergeRequest.Title,
		SourceProjectID: gitlabMergeRequest.SourceProjectID,
		SourceBranch:    gitlabMergeRequest.SourceBranch,
		CommitSHA:       gitlabMergeRequest.SHA,
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

func (gr *gitlabRepository) GetMergeRequests(projectID int) ([]models.MergeRequest, error) {
	var mergeRequests []models.MergeRequest

	gitlabMergeRequests, _, err := gr.mergeRequestsService.ListProjectMergeRequests(projectID, &gitlab.ListProjectMergeRequestsOptions{
		// If needed by users, use pagination.
		ListOptions: gitlab.ListOptions{
			Page:    1,
			PerPage: 100, // Maximum par_page allowed.
		},
		State: pointer.ToString("opened"),
	})
	if err != nil {
		return nil, err
	}

	for _, gitlabMergeRequest := range gitlabMergeRequests {
		mergeRequest := models.MergeRequest{
			ID:              gitlabMergeRequest.IID,
			Title:           gitlabMergeRequest.Title,
			SourceProjectID: gitlabMergeRequest.SourceProjectID,
			SourceBranch:    gitlabMergeRequest.SourceBranch,
			CommitSHA:       gitlabMergeRequest.SHA,
		}

		if gitlabMergeRequest.Author != nil {
			mergeRequest.Author.Name = gitlabMergeRequest.Author.Name
			mergeRequest.Author.AvatarURL = gitlabMergeRequest.Author.AvatarURL

			if mergeRequest.Author.Name == "" {
				mergeRequest.Author.Name = gitlabMergeRequest.Author.Username
			}
		}

		mergeRequests = append(mergeRequests, mergeRequest)
	}

	return mergeRequests, nil
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
