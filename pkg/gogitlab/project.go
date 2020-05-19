package gogitlab

import (
	"github.com/xanzy/go-gitlab"
)

type Project interface {
	GetProject(pid interface{}, opt *gitlab.GetProjectOptions, options ...gitlab.RequestOptionFunc) (*gitlab.Project, *gitlab.Response, error)
}
