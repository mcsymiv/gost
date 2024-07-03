package models

type DriverStatus struct {
	Message string `json:"message"`
	Ready   bool   `json:"ready"`
}

type Session struct {
	Id string `json:"sessionId"`
}

type Url struct {
	Url string `json:"url"`
}
