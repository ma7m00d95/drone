package services

import (
	"drone/database"
	model "drone/models"
	"fmt"
	"log"
)

func GetOrderByID(orderID string) (model.Orders, error) {
	query := `select * from orders where id = ?`
	order := model.Orders{}
	err := database.DB.QueryRow(query, orderID).Scan(&order.ID, &order.Origin, &order.Destination, &order.Status, &order.AssignedDroneID, &order.CurrentLat, &order.CurrentLng, &order.CreatedBy)
	if err != nil {
		log.Println("Failed to get order:", err)
		return model.Orders{}, err
	}

	return order, nil
}
func GetOrders() ([]model.Orders, error) {
	query := `select * from orders`

	rows, err := database.DB.Query(query)
	if err != nil {
		log.Println("Failed to get order:", err)
		return []model.Orders{}, err
	}
	defer rows.Close()

	orderList := []model.Orders{}

	for rows.Next() {
		var order model.Orders
		err := rows.Scan(
			&order.ID,
			&order.Origin,
			&order.Destination,
			&order.Status,
			&order.AssignedDroneID,
			&order.CurrentLat,
			&order.CurrentLng,
			&order.CreatedBy,
		)
		if err != nil {
			log.Println("Scan error:", err)
			return []model.Orders{}, err
		}

		orderList = append(orderList, order)
	}

	return orderList, nil
}
func GetUserOrders(userID string) ([]model.Orders, error) {
	query := `select * from orders where created_by = ?`

	rows, err := database.DB.Query(query, userID)
	if err != nil {
		log.Println("Failed to get order:", err)
		return []model.Orders{}, err
	}
	defer rows.Close()

	orderList := []model.Orders{}

	for rows.Next() {
		var order model.Orders
		err := rows.Scan(
			&order.ID,
			&order.Origin,
			&order.Destination,
			&order.Status,
			&order.AssignedDroneID,
			&order.CurrentLat,
			&order.CurrentLng,
			&order.CreatedBy,
		)
		if err != nil {
			log.Println("Scan error:", err)
			return []model.Orders{}, err
		}

		orderList = append(orderList, order)
	}
	return orderList, nil
}

func GetInProcessOrders() ([]model.Orders, error) {
	query := `select * from orders where status = "InProcess"`

	rows, err := database.DB.Query(query)
	if err != nil {
		log.Println("Failed to get order:", err)
		return []model.Orders{}, err
	}
	defer rows.Close()

	orderList := []model.Orders{}

	for rows.Next() {
		var order model.Orders
		err := rows.Scan(
			&order.ID,
			&order.Origin,
			&order.Destination,
			&order.Status,
			&order.AssignedDroneID,
			&order.CurrentLat,
			&order.CurrentLng,
			&order.CreatedBy,
		)
		if err != nil {
			log.Println("Scan error:", err)
			return []model.Orders{}, err
		}

		orderList = append(orderList, order)
	}

	return orderList, nil
}
func IsDroneID(id int) bool {
	query := ` select id from drones where id =?  `
	var droneID int
	err := database.DB.QueryRow(query, id).Scan(&droneID)
	if err != nil {
		return false
	}
	return true
}
func IsOrderID(id int) bool {
	query := ` select id from orders where id =?  `
	var orderID int64
	err := database.DB.QueryRow(query, id).Scan(&orderID)
	if err != nil {
		return false
	}
	return true
}
func CreateOrder(order model.Orders) (int, error) {
	query := `insert into orders 
	(origin, destination, status, assigned_drone_id, current_lat, current_lng, created_by) 
	values ( ?, ?, "Pending", 0, 0.0, 0.0, ?)`

	res, err := database.DB.Exec(query,
		order.Origin,
		order.Destination,
		order.CreatedBy)
	if err != nil {
		log.Println("Failed to create order:", err)
		return 0, err
	}
	id, err := res.LastInsertId()

	return int(id), err
}
func AssignDroneToOrder(orderID int) error {
	if !IsOrderID(orderID) {
		return fmt.Errorf("Invalid order ID")
	}
	query := `SELECT id FROM drones WHERE status = 'Fixed' LIMIT 1`
	var droneID int

	err := database.DB.QueryRow(query).Scan(&droneID)
	if err != nil {
		return nil
	}

	//update order
	_, err = database.DB.Exec(`
		UPDATE orders 
		SET assigned_drone_id = ?, status = 'Assigned'
		WHERE id = ?
	`, droneID, orderID)
	if err != nil {
		return err
	}

	//update drone
	_, err = database.DB.Exec(`
		UPDATE drones 
		SET status = 'Busy', current_order_id = ?
		WHERE id = ?
	`, orderID, droneID)

	return err
}
func GetOrderStatus(orderID int) (string, error) {
	if !IsOrderID(orderID) {
		return "", fmt.Errorf("Invalid order ID")
	}
	query := ` select status from orders where id =?  `
	status := ""
	err := database.DB.QueryRow(query, orderID).Scan(&status)
	if err != nil {
		log.Println("Failed to check delivery status:", err)
		return status, err
	}
	return status, nil
}
func UpdateOrderOrigin(l model.Location) error {
	if !IsOrderID(l.OrderID) {
		return fmt.Errorf("Invalid order ID")
	}
	query := `update orders set origin=? where id=?`
	_, err := database.DB.Exec(query, l.Origin, l.OrderID)
	return err
}
func UpdateOrderDestination(d model.Location) error {
	if !IsOrderID(d.OrderID) {
		return fmt.Errorf("Invalid order ID")
	}
	query := `update orders set destination=? where id=?`
	_, err := database.DB.Exec(query, d.Destination, d.OrderID)
	return err
}
func CancelOrder(orderID int) error {
	if !IsOrderID(orderID) {
		return fmt.Errorf("Invalid order ID")
	}
	query := `update orders set status = "Cancelled" where id=? and status = "Pending"`
	_, err := database.DB.Exec(query, orderID)
	return err
}
func UpdateOrder(status string, orderID int) error {
	if !IsOrderID(orderID) {
		return fmt.Errorf("Invalid order ID")
	}
	//only the order that are picked in DB can be changed to these status
	if status != "Delivered" &&
		status != "Failed" {
		return fmt.Errorf("Invalid status")
	}
	query := `update orders set status=? where id=? and status = "InProcess"`
	id, err := database.DB.Exec(query, status, orderID)
	if err != nil {
		return err
	}
	if i, _ := id.RowsAffected(); i == 0 {
		return fmt.Errorf("Order not in 'InProcess' status")
	}
	return err
}

func DeliverOrder(orderID int) error {
	if !IsOrderID(orderID) {
		return fmt.Errorf("Invalid order ID")
	}
	query := `update orders set status = "Delivered" where id=? and status = "InProcess"`
	_, err := database.DB.Exec(query, orderID)
	return err

}

// internal API
func FreeDrone(droneID int) error {
	if !IsDroneID(droneID) {
		return fmt.Errorf("Invalid drone ID")
	}
	query := `update drones set status = "Available", current_order_id = 0 where id=?`
	_, err := database.DB.Exec(query, droneID)

	return err
}
func FailOrder(orderID int) error {
	if !IsOrderID(orderID) {
		return fmt.Errorf("Invalid order ID")
	}
	query := `update orders set status = "Failed", assigned_drone_id = 0 where id=? and status = "Assigned"`
	_, err := database.DB.Exec(query, orderID)
	return err
}

func IsUser(userID int) error {
	query := ` select id from users where id =?  `
	var id string
	err := database.DB.QueryRow(query, userID).Scan(&id)
	if err != nil {
		return fmt.Errorf("User not found")
	}
	return nil
}
