package models

const (
	Good Response = "Good"
	Bad  Response = "Bad"
)

type Response string

type MessageResponse struct {
	Response Response `json:"response"`
}

type MessageRequest struct {
	from   string
	to     string
	amount float64
}
