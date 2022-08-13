package project

import (
	"issue-service/app/issue-api/routes/models"
	"strings"

	"gorm.io/gorm"
)

type projectRouter struct {
	routes   models.Routes
	database *gorm.DB
}

func NewRouter(db *gorm.DB) models.Router {
	r := &projectRouter{database: db}
	r.initRoutes()
	return r
}

// Routes returns the available routers to the checkpoint controller
func (r *projectRouter) Routes() models.Routes {
	return r.routes
}

func (r *projectRouter) initRoutes() {
	r.routes = models.Routes{
		models.Route{
			Name:        "AddProject",
			Method:      strings.ToUpper("Post"),
			Pattern:     "/v1/projects",
			HandlerFunc: CreateAddProjectHandler,
		},

		models.Route{
			Name:        "GetProjects",
			Method:      strings.ToUpper("Get"),
			Pattern:     "/v1/projects",
			HandlerFunc: CreateGetProjectsHandler,
		},
	}
}
