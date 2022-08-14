package sprint

import (
	"issue-service/app/issue-api/routes/models"
	"strings"

	"gorm.io/gorm"
)

type sprintRouter struct {
	routes   models.Routes
	database *gorm.DB
}

func NewRouter(db *gorm.DB) models.Router {
	r := &sprintRouter{database: db}
	r.initRoutes()
	return r
}

// Routes returns the available routers to the checkpoint controller
func (r *sprintRouter) Routes() models.Routes {
	return r.routes
}

func (r *sprintRouter) initRoutes() {
	r.routes = models.Routes{
		models.Route{
			Name:        "AddSprint",
			Method:      strings.ToUpper("Post"),
			Pattern:     "/v1/projects/{projectId}/sprints",
			HandlerFunc: createAddSprintHandler,
		},

		models.Route{
			Name:        "PatchSprint",
			Method:      strings.ToUpper("Patch"),
			Pattern:     "/v1/projects/{projectId}/sprints/{sprintId}",
			HandlerFunc: createPatchSprintHandler,
		},
	}
}
