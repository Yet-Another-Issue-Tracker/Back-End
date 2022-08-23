package issue

import (
	"issue-service/app/issue-api/routes/models"
	"strings"

	"gorm.io/gorm"
)

type issueRouter struct {
	routes   models.Routes
	database *gorm.DB
}

func NewRouter(db *gorm.DB) models.Router {
	r := &issueRouter{database: db}
	r.initRoutes()
	return r
}

// Routes returns the available routers to the checkpoint controller
func (r *issueRouter) Routes() models.Routes {
	return r.routes
}

func (r *issueRouter) initRoutes() {
	r.routes = models.Routes{
		models.Route{
			Name:        "AddIssue",
			Method:      strings.ToUpper("Post"),
			Pattern:     "/v1/projects/{projectId}/sprints/{sprintId}/issues",
			HandlerFunc: createAddIssueHandler,
		},

		models.Route{
			Name:        "GetIssues",
			Method:      strings.ToUpper("Get"),
			Pattern:     "/v1/projects/{projectId}/sprints/{sprintId}/issues",
			HandlerFunc: createGetIssuesHandler,
		},

		models.Route{
			Name:        "PatchIssue",
			Method:      strings.ToUpper("Patch"),
			Pattern:     "/v1/projects/{projectId}/sprints/{sprintId}/issues/{issueId}",
			HandlerFunc: createPatchIssueHandler,
		},
	}
}
