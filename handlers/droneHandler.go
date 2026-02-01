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

func UpdateDroneStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	status := r.URL.Query().Get("status")
	status = string(status) //it has to be lowercase
	id := r.URL.Query().Get("drone_id")
	if id == "" || status != "fixed" && status != "broken" {
		http.Error(w, "Invalid drone_id or status", http.StatusBadRequest)
		return
	}

	err := services.UpdateDroneStatus(id, status)
	services.TryAssignPendingOrder()

	if err != nil {
		http.Error(w, "Failed to update drone status", http.StatusInternalServerError)
		return
	}
	w.Write([]byte("drone status updated"))

}
