package handlers

import (
	model "drone/models"
	services "drone/services"
	"encoding/json"
	"net/http"
)

func CreateOrder(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	//get the whole order model from user and give it to the function

	var order model.Orders

	json.NewDecoder(r.Body).Decode(&order)
	isDroneID := services.IsDroneID(order.AssignedDroneID)
	if !isDroneID {
		http.Error(w, "Invalid drone ID", http.StatusBadRequest)
		return
	}
	_, err := services.CreateOrder(order)
	if err != nil {
		http.Error(w, "Failed to create order", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(order)

}
func GetOrderByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	request := r.URL.Query().Get("order_id")
	location, err := services.GetOrderByID(request)
	if err != nil {
		http.Error(w, "Failed to get order", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(location)

}
func GetOrders(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	location, err := services.GetOrders()
	if err != nil {
		http.Error(w, "Failed to get order", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(location)

}
func GetOrderStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	req := r.URL.Query().Get("order_id")
	status, err := services.GetOrderStatus(req)
	if err != nil {
		http.Error(w, "Failed to get order status", http.StatusInternalServerError)
		return
	}
	w.Write([]byte("Order is : " + status))

}
func UpdateOrderOrigin(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	var origin model.Location
	json.NewDecoder(r.Body).Decode(&origin)

	err := services.UpdateOrderOrigin(origin)
	if err != nil {
		http.Error(w, "Failed to update order origin", http.StatusInternalServerError)
		return
	}
	w.Write([]byte("Order origin updated successfully"))

}
func UpdateOrderDestination(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	var destination model.Location
	json.NewDecoder(r.Body).Decode(&destination)
	err := services.UpdateOrderDestination(destination)
	if err != nil {
		http.Error(w, "Failed to update order destination", http.StatusInternalServerError)
		return
	}
	w.Write([]byte("Order destination updated successfully"))

}
