package routes

type Project struct {
	Id string `json:"id,omitempty"`

	Client string `json:"client,omitempty"`

	Type_ string `json:"type,omitempty"`
}

type CreateProjectResponse struct {
	Id string `json:"id,omitempty"`
}
