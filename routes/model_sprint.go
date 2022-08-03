package routes

type Sprint struct {
	Number string `json:"number,omitempty"`

	StartDate string `json:"startDate,omitempty"`

	EndDate string `json:"endDate,omitempty"`

	Completed bool `json:"completed,omitempty"`
}
