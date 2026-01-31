package handlers

import (
	services "drone/services"
	"net/http"
)

func CreateDrone(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	err := services.CreateDrone()
	if err != nil {
		http.Error(w, "Failed to create drone", http.StatusInternalServerError)
		return
	}
	w.Write([]byte("Drone created successfully"))

}

// func GetOrderByID(w http.ResponseWriter, r *http.Request) {
// 	if r.Method != "GET" {
// 		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
// 		return
// 	}
// 	request := r.URL.Query().Get("order_id")
// 	location, err := services.GetOrderByID(request)
// 	if err != nil {
// 		http.Error(w, "Failed to get order", http.StatusInternalServerError)
// 		return
// 	}
// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(location)

// }
// func GetOrders(w http.ResponseWriter, r *http.Request) {
// 	if r.Method != "GET" {
// 		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
// 		return
// 	}
// 	location, err := services.GetOrders()
// 	if err != nil {
// 		http.Error(w, "Failed to get order", http.StatusInternalServerError)
// 		return
// 	}
// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(location)

// }
func UpdateDroneStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	id := r.URL.Query().Get("drone_id")
	status := r.URL.Query().Get("status")

	err := services.UpdateDroneStatus(id, status)
	if err != nil {
		http.Error(w, "Failed to update drone status", http.StatusInternalServerError)
		return
	}
	w.Write([]byte("drone status updated"))

}

// func UpdateOrderOrigin(w http.ResponseWriter, r *http.Request) {
// 	if r.Method != "POST" {
// 		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
// 		return
// 	}
// 	var origin model.Location
// 	json.NewDecoder(r.Body).Decode(&origin)

// 	err := services.UpdateOrderOrigin(origin)
// 	if err != nil {
// 		http.Error(w, "Failed to update order origin", http.StatusInternalServerError)
// 		return
// 	}
// 	w.Write([]byte("Order origin updated successfully"))

// }
// func UpdateOrderDestination(w http.ResponseWriter, r *http.Request) {
// 	if r.Method != "POST" {
// 		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
// 		return
// 	}
// 	var destination model.Location
// 	json.NewDecoder(r.Body).Decode(&destination)
// 	err := services.UpdateOrderDestination(destination)
// 	if err != nil {
// 		http.Error(w, "Failed to update order destination", http.StatusInternalServerError)
// 		return
// 	}
// 	w.Write([]byte("Order destination updated successfully"))

// }
