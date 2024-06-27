package domain

type Sensor struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	IsActive bool   `json:"is_active"`
}
