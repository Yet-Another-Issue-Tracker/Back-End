package routes

type Issue struct {
	Id string `json:"id,omitempty"`

	Type_ string `json:"type,omitempty"`

	Title string `json:"title,omitempty"`

	Description string `json:"description,omitempty"`

	Status string `json:"status,omitempty"`

	Assignee string `json:"assignee,omitempty"`
}
