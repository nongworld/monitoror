package usecase

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/monitoror/monitoror/monitorables/gitlab/api"
	"github.com/monitoror/monitoror/monitorables/gitlab/api/mocks"
	"github.com/monitoror/monitoror/monitorables/gitlab/api/models"

	"github.com/jsdidierlaurent/echo-middleware/cache"
	"github.com/stretchr/testify/mock"
)

func initUsecase(mockRepository api.Repository) api.Usecase {
	store := cache.NewGoCacheStore(time.Minute*5, time.Second)
	pu := NewGitlabUsecase(mockRepository, store)
	return pu
}

func TestUsecase_getProject(t *testing.T) {
	mockRepository := new(mocks.Repository)
	mockRepository.On("GetProject", mock.AnythingOfType("int")).
		Return(&models.Project{Repository: "TEST"}, nil)

	gu := initUsecase(mockRepository)
	castedGu := gu.(*gitlabUsecase)

	project, err := castedGu.getProject(10)
	assert.NoError(t, err)
	assert.Equal(t, "TEST", project.Repository)

	project, err = castedGu.getProject(10)
	assert.NoError(t, err)
	assert.Equal(t, "TEST", project.Repository)

	mockRepository.AssertNumberOfCalls(t, "GetProject", 1)
	mockRepository.AssertExpectations(t)
}
