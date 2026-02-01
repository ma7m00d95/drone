package services

import (
	"drone/database"
	model "drone/models"
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
func IsDroneID(id int64) bool {
	query := ` select id from drones where id =?  `
	var droneID int64
	err := database.DB.QueryRow(query, id).Scan(&droneID)
	if err != nil {
		return false
	}
	return true
}

func CreateOrder(order model.Orders) (int64, error) {
	query := `insert into orders 
	(origin, destination, status, assigned_drone_id, current_lat, current_lng, created_by) 
	values ( ?, ?, "pending", 0, 0.0, 0.0, ?)`

	res, err := database.DB.Exec(query,
		order.Origin,
		order.Destination,
		order.CreatedBy)
	if err != nil {
		log.Println("Failed to create order:", err)
		return 0, err
	}
	id, err := res.LastInsertId()

	return id, nil
}
func AssignDroneToOrder(orderID int64) error {
	query := `SELECT id FROM drones WHERE status = 'fixed' LIMIT 1`
	var droneID int

	err := database.DB.QueryRow(query).Scan(&droneID)
	if err != nil {
		// no available drone → order stays pending
		return nil
	}

	// 2) update order
	_, err = database.DB.Exec(`
		UPDATE orders 
		SET assigned_drone_id = ?, status = 'assigned'
		WHERE id = ?
	`, droneID, orderID)
	if err != nil {
		return err
	}

	// 3) update drone
	_, err = database.DB.Exec(`
		UPDATE drones 
		SET status = 'busy', current_order_id = ?
		WHERE id = ?
	`, orderID, droneID)

	return err
}
func GetOrderStatus(orderID string) (string, error) {
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
	query := `update orders set origin=? where id=?`
	_, err := database.DB.Exec(query, l.Origin, l.OrderID)
	return err
}
func UpdateOrderDestination(d model.Location) error {
	query := `update orders set destination=? where id=?`
	_, err := database.DB.Exec(query, d.Destination, d.OrderID)
	return err
}
func CancelOrder(orderID string) error {
	query := `update orders set status = "cancelled" where id=? and status = "pending"`
	_, err := database.DB.Exec(query, orderID)
	return err
}

// Reserve a job.
// ● Grab an order from a location (origin or broken drone).
// ● Mark an order they have gotten as delivered or failed.
// ● Mark themselves as broken (and in need of an order handoff).
// ● Update their location (use latitude/longitude), and get a status update as a heartbeat.
// ● Get details on the order they are currently assigned.
