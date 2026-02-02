package handlers

import (
	model "drone/models"
	services "drone/services"
	"encoding/json"
	"net/http"
)

func CreateOrder(w http.ResponseWriter, r *http.Request) {
	var order model.Orders

	json.NewDecoder(r.Body).Decode(&order)

	id, err := services.CreateOrder(order)
	if err != nil {
		http.Error(w, "Failed to create order", http.StatusInternalServerError)
		return
	}
	services.TryAssignPendingOrder()

	err = services.AssignDroneToOrder(id)
	if err != nil {
		http.Error(w, "Failed to assign drone to order", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode("Order created successfully")

}

func CancelOrder(w http.ResponseWriter, r *http.Request) {

	orderID := r.PathValue("order_id")

	err := services.CancelOrder(orderID)
	if err != nil {
		http.Error(w, "Failed to cancel order", http.StatusInternalServerError)
		return
	}
	services.TryAssignPendingOrder()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode("Order cancelled successfully")

}
func GetOrderByID(w http.ResponseWriter, r *http.Request) {
	request := r.PathValue("order_id")
	location, err := services.GetOrderByID(request)
	if err != nil {
		http.Error(w, "Failed to get order", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(location)

}
func GetOrders(w http.ResponseWriter, r *http.Request) {

	location, err := services.GetOrders()
	if err != nil {
		http.Error(w, "Failed to get order", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(location)

}
func GetOrderStatus(w http.ResponseWriter, r *http.Request) {

	req := r.PathValue("order_id")
	status, err := services.GetOrderStatus(req)
	if err != nil {
		http.Error(w, "Failed to get order status", http.StatusInternalServerError)
		return
	}
	w.Write([]byte("Order is : " + status))

}
func UpdateOrderOrigin(w http.ResponseWriter, r *http.Request) {

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

	var destination model.Location
	json.NewDecoder(r.Body).Decode(&destination)
	err := services.UpdateOrderDestination(destination)
	if err != nil {
		http.Error(w, "Failed to update order destination", http.StatusInternalServerError)
		return
	}
	w.Write([]byte("Order destination updated successfully"))

}
func UpdateOrder(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	status := r.PathValue("status")
	orderID := r.PathValue("order_id")

	err := services.UpdateOrder(status, orderID)
	if err != nil {
		http.Error(w, "Failed to update order", http.StatusInternalServerError)
		return
	}
	w.Write([]byte("Order updated successfully"))

}
