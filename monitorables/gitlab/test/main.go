package main

import (
	"fmt"
	"time"

	"github.com/AlekSi/pointer"
	"github.com/jsdidierlaurent/echo-middleware/cache"

	"github.com/monitoror/monitoror/monitorables/gitlab/api/models"
	"github.com/monitoror/monitoror/monitorables/gitlab/api/repository"
	"github.com/monitoror/monitoror/monitorables/gitlab/api/usecase"
	"github.com/monitoror/monitoror/monitorables/gitlab/config"
)

func main() {
	store := cache.NewGoCacheStore(time.Minute*5, time.Second)

	conf := &config.Gitlab{
		URL:     "https://git.sarbacane.com",
		Token:   "WdG9Hr4Lf_Cd49vxAzKN",
		Timeout: 5000,
	}

	repo := repository.NewGitlabRepository(conf)
	uc := usecase.NewGitlabUsecase(repo, store)

	projectID := 1

	ids, _ := repo.GetMergeRequests(projectID)
	for _, id := range ids {
		mr, err := uc.MergeRequest(&models.MergeRequestParams{
			ProjectID: pointer.ToInt(projectID),
			ID:        pointer.ToInt(id),
		})
		if err != nil {
			panic(err)
		}

		fmt.Printf("Project: %s\n", mr.Label)
		fmt.Printf("Merge Request: #%d\n", mr.Build.MergeRequest.ID)
		fmt.Printf("Branch: %s\n", *mr.Build.Branch)
		fmt.Printf("Author: %s\n", mr.Build.Author)
		fmt.Printf("Pipeline: %s\n", mr.Status)

		fmt.Println()
	}
}
