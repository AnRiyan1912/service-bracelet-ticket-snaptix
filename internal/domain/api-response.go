package domain

type ApiResponse struct {
	Error   bool        `json:"Error"`
	Message string      `json:"Message"`
	Data    interface{} `json:"Data,omitempty"`
}

type ApiResponseWithaoutData struct {
	Error   bool   `json:"Error"`
	Message string `json:"Message"`
}
