//go:generate mockery -name Usecase

package api

import (
	uiConfigModels "github.com/monitoror/monitoror/api/config/models"
	coreModels "github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorables/gitlab/api/models"
)

const (
	GitlabPipelineTileType     coreModels.TileType = "GITHUB-PIPELINE"
	GitlabMergeRequestTileType coreModels.TileType = "GITHUB-MERGEREQUEST"
)

type (
	Usecase interface {
		Pipeline(params *models.PipelineParams) (*coreModels.Tile, error)
		MergeRequest(params *models.MergeRequestParams) (*coreModels.Tile, error)

		MergeRequestsGenerator(params interface{}) ([]uiConfigModels.GeneratedTile, error)
	}
)
