package handlers

import (
	model "drone/models"
	services "drone/services"
	"encoding/json"
	"net/http"
	"strconv"
)

func CreateOrder(w http.ResponseWriter, r *http.Request) {
	var order model.Orders

	json.NewDecoder(r.Body).Decode(&order)
	if err := services.IsUser(order.CreatedBy); err != nil {
		http.Error(w, "Invalid user", http.StatusBadRequest)
		return
	}
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

	orderID, err := strconv.Atoi(r.PathValue("order_id"))
	if err != nil {
		http.Error(w, "Invalid order ID", http.StatusBadRequest)
		return
	}

	err = services.CancelOrder(orderID)
	if err != nil {
		http.Error(w, "Failed to cancel order", http.StatusInternalServerError)
		return
	}
	drone, err := services.GetDroneByOrderID(orderID)
	if err != nil {
		http.Error(w, "Failed to get drone by order ID", http.StatusInternalServerError)
		return
	}
	err = services.FreeDrone(drone.ID)
	if err != nil {
		http.Error(w, "Failed to free drone", http.StatusInternalServerError)
		return
	}
	services.TryAssignPendingOrder()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode("Order cancelled successfully")

}
func GetOrderByID(w http.ResponseWriter, r *http.Request) {
	request := r.PathValue("order_id")
	order, err := services.GetOrderByID(request)
	if err != nil {
		http.Error(w, "Failed to get order", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(order)

}
func GetOrders(w http.ResponseWriter, r *http.Request) {

	orders, err := services.GetOrders()
	if err != nil {
		http.Error(w, "Failed to get order", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orders)

}
func GetUserOrders(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("user_id")
	orders, err := services.GetUserOrders(userID)
	if err != nil {
		http.Error(w, "Failed to get user orders", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orders)

}
func GetInProcessOrders(w http.ResponseWriter, r *http.Request) {

	orders, err := services.GetInProcessOrders()
	if err != nil {
		http.Error(w, "Failed to get order", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orders)

}
func GetOrderStatus(w http.ResponseWriter, r *http.Request) {

	req := r.PathValue("order_id")
	orderID, err := strconv.Atoi(req)
	if err != nil {
		http.Error(w, "Invalid order ID", http.StatusBadRequest)
		return
	}
	status, err := services.GetOrderStatus(orderID)
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
	status := r.PathValue("status")
	orderID, err := strconv.Atoi(r.PathValue("order_id"))
	if err != nil {
		http.Error(w, "Invalid order ID", http.StatusBadRequest)
		return
	}

	err = services.UpdateOrder(status, orderID)
	if err != nil {
		http.Error(w, "Failed to update order", http.StatusInternalServerError)
		return
	}
	w.Write([]byte("Order updated successfully"))

}

func DeliverOrder(w http.ResponseWriter, r *http.Request) {
	var order model.Orders
	json.NewDecoder(r.Body).Decode(&order)
	err := services.DeliverOrder(order.ID)
	if err != nil {
		http.Error(w, "Failed to deliver order by drone", http.StatusInternalServerError)
		return
	}
	services.FreeDrone(order.AssignedDroneID)
	w.Write([]byte("Order delivered by drone successfully"))

}
func FailOrder(w http.ResponseWriter, r *http.Request) {
	var order model.Orders
	json.NewDecoder(r.Body).Decode(&order)
	err := services.FailOrder(order.ID)
	if err != nil {
		http.Error(w, "Failed to mark order as failed by drone", http.StatusInternalServerError)
		return
	}
	drone, err := services.GetDroneByOrderID(order.ID)
	if err != nil {
		http.Error(w, "Failed to get drone by order ID", http.StatusInternalServerError)
		return
	}
	services.FreeDrone(drone.ID)
	services.TryAssignPendingOrder()

	w.Write([]byte("Order marked as failed by drone successfully"))

}
