package models

type Location struct {
	OrderID     int    `json:"order_id"`
	Origin      string `json:"origin"`
	Destination string `json:"destination"`
}
type Users struct {
	UserID   int    `json:"user_id"`
	UserName string `json:"user_name"`
	UserType string `json:"user_type"`
}
type Orders struct {
	ID              int     `json:"order_id"`
	Origin          string  `json:"origin"`
	Destination     string  `json:"destination"`
	Status          string  `json:"status"`
	AssignedDroneID int     `json:"assigned_drone_id"`
	CurrentLat      float64 `json:"current_lat"`
	CurrentLng      float64 `json:"current_lng"`
	CreatedBy       int     `json:"created_by"`
}
type Drone struct {
	ID      int     `json:"drone_id"`
	Status  string  `json:"status"`
	Lat     float64 `json:"lat"`
	Lng     float64 `json:"lng"`
	OrderID int     `json:"order_id"`
}
