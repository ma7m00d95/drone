package handlers

import (
	model "drone/models"
	services "drone/services"
	"encoding/json"
	"net/http"
)

func CreateDrone(w http.ResponseWriter, r *http.Request) {

	err := services.CreateDrone()
	if err != nil {
		http.Error(w, "Failed to create drone", http.StatusInternalServerError)
		return
	}
	w.Write([]byte("Drone created successfully"))

}

func UpdateDroneStatus(w http.ResponseWriter, r *http.Request) {
	status := r.PathValue("status")
	status = string(status) //it has to be lowercase
	id := r.PathValue("drone_id")
	if id == "" || status != "Fixed" && status != "Broken" {
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

func ReserveOrderToDrone(w http.ResponseWriter, r *http.Request) {
	var drone model.Drone
	json.NewDecoder(r.Body).Decode(&drone)
	err := services.ReserveOrderToDrone(drone.ID)
	if err != nil {
		http.Error(w, "Failed to reserve order to drone", http.StatusInternalServerError)
		return
	}
	w.Write([]byte("Order reserved to drone successfully"))

}

func DronePickupOrder(w http.ResponseWriter, r *http.Request) {
	// in case that the drone going to grab order it will grab any order and change it to inprocess

	var drone model.Drone
	json.NewDecoder(r.Body).Decode(&drone)
	err := services.DronePickupOrder(drone.ID)
	if err != nil {
		http.Error(w, "Failed to pickup order by drone", http.StatusInternalServerError)
		return
	}
	w.Write([]byte("Order picked up by drone successfully"))
	
}
func UpdateDroneLocation(w http.ResponseWriter, r *http.Request) {
	var drone model.Drone
	json.NewDecoder(r.Body).Decode(&drone)

	err := services.UpdateDroneLocation(drone)
	if err != nil {
		http.Error(w, "Failed to update drone location", http.StatusInternalServerError)
		return
	}
	w.Write([]byte("Drone location updated successfully"))

}
func GetDroneCurrentOrder(w http.ResponseWriter, r *http.Request) {
	droneID := r.PathValue("drone_id")
	order, err := services.GetDroneCurrentOrder(droneID)
	if err != nil {
		http.Error(w, "Failed to get drone current order", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(order)

}