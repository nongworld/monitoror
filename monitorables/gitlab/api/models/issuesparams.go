package models

import (
	"github.com/monitoror/monitoror/internal/pkg/monitorable/params"
)

type (
	IssuesParams struct {
		params.Default

		ProjectID *int   `json:"projectId" query:"projectId"`
		Query     string `json:"query" query:"query" validate:"required"`
	}
)
