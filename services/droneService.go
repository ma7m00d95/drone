package services

import (
	"drone/database"
	model "drone/models"
	"fmt"
	"log"
)

func GetDroneByID(droneID int) (model.Drone, error) {
	query := `select * from drones where id = ?`
	drone := model.Drone{}
	err := database.DB.QueryRow(query, droneID).Scan(&drone.ID, &drone.Status, &drone.Lat, &drone.Lng, &drone.OrderID)
	if err != nil {
		log.Println("Failed to get drone:", err)
		return model.Drone{}, err
	}

	return drone, nil
}

func GetDroneByOrderID(orderID int) (model.Drone, error) {
	query := `select id from drones where current_order_id = ?`
	var drone model.Drone
	err := database.DB.QueryRow(query, orderID).Scan(&drone.ID)
	if err != nil {
		log.Println("Failed to get drone:", err)
		return model.Drone{}, err
	}

	return drone, nil
}
func GetDrones() ([]model.Drone, error) {
	query := `select * from drones`

	rows, err := database.DB.Query(query)
	if err != nil {
		log.Println("Failed to get drones:", err)
		return []model.Drone{}, err
	}
	defer rows.Close()

	droneList := []model.Drone{}
	for rows.Next() {
		var drone model.Drone
		err := rows.Scan(
			&drone.ID,
			&drone.Status,
			&drone.Lat,
			&drone.Lng,
			&drone.OrderID,
		)
		if err != nil {
			log.Println("Scan error:", err)
			return []model.Drone{}, err
		}
		droneList = append(droneList, drone)
	}
	return droneList, nil
}

func GetDroneStatus(droneID int) (string, error) {
	query := ` select status from drones where id =?  `
	status := ""
	err := database.DB.QueryRow(query, droneID).Scan(&status)
	if err != nil {
		return status, err
	}
	return status, nil
}
func UpdateDroneStatus(droneID string, status string) error {
	query := ` update drones set status=? where id =?  `
	_, err := database.DB.Exec(query, status, droneID)
	if err != nil {
		return err
	}
	if status == "Fixed" {
		return nil
	}
	// get the ID
	var orderID string
	err = database.DB.QueryRow(`select current_order_id from drones where id=?`, droneID).Scan(&orderID)
	if err != nil {
		return err
	}
	//update order status
	query = ` update orders set status=? where id =?  `
	_, err = database.DB.Exec(query, "Pickup", orderID)
	if err != nil {
		return err
	}
	return nil
}

func UpdateDroneOrigin(l model.Location) error {
	query := `update orders set origin=? where id=?`
	_, err := database.DB.Exec(query, l.Origin, l.OrderID)
	return err
}

func CreateDrone() error {
	query := `insert into drones 
	(status, lat, lng, current_order_id) 
	values ( "Available", 0.0, 0.0, 0)`

	_, err := database.DB.Exec(query)
	if err != nil {
		log.Println("Failed to create drone:", err)
		return err
	}

	return nil
}
func TryAssignPendingOrder() error {
	// 1) find pending order
	var orderID string
	err := database.DB.QueryRow(`
		SELECT id FROM orders 
		WHERE status = 'Pending' AND assigned_drone_id IS 0
		ORDER BY id ASC
		LIMIT 1
	`).Scan(&orderID)

	if err != nil {
		// no pending orders
		return nil
	}

	// 2) find available drone
	var droneID string
	err = database.DB.QueryRow(`
		SELECT id FROM drones 
		WHERE status = 'Available' OR status = 'Fixed' AND current_order_id IS 0
		LIMIT 1
	`).Scan(&droneID)

	if err != nil {
		// no available drone
		return nil
	}

	// 3) assign
	_, err = database.DB.Exec(`
		UPDATE orders 
		SET assigned_drone_id = ?, status = 'Assigned'
		WHERE id = ?
	`, droneID, orderID)
	if err != nil {
		return err
	}

	_, err = database.DB.Exec(`
		UPDATE drones 
		SET current_order_id = ?, status = 'Busy'
		WHERE id = ?
	`, orderID, droneID)

	return err
}
func ReserveOrderToDrone(droneID int) error {
	// find pending order
	var orderID string
	err := database.DB.QueryRow(`
		SELECT id FROM orders 
		WHERE status = 'Pending' AND assigned_drone_id = 0 
		ORDER BY id ASC
		LIMIT 1
	`).Scan(&orderID)
	if err != nil {
		return err
	}
	//we will get one order or 0
	_, err = database.DB.Exec(`
		UPDATE orders 
		SET assigned_drone_id = ?, status = 'Assigned'
		WHERE id = ?
	`, droneID, orderID)
	if err != nil {
		return err
	}

	_, err = database.DB.Exec(`
		UPDATE drones
		SET current_order_id = ?, status = 'Busy'
		WHERE id = ?
	`, orderID, droneID)

	return err
}
func DronePickupOrder(droneID int) error {
	// get current order assigned to drone
	var orderID string
	err := database.DB.QueryRow(`
		SELECT current_order_id FROM drones 
		WHERE id = ?
	`, droneID).Scan(&orderID)
	if err != nil {
		return err
	}

	// update order status to inprocess
	_, err = database.DB.Exec(`
		UPDATE orders 
		SET status = 'InProcess'
		WHERE id = ?
	`, orderID)
	if err != nil {
		return err
	}

	return nil
}

func UpdateDroneLocation(drone model.Drone) error {
	query := ` update drones set lat=?, lng=? where id =?  `
	record, err := database.DB.Exec(query, drone.Lat, drone.Lng, drone.ID)
	if err != nil {
		return err
	}
	affected, err := record.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return fmt.Errorf("No drone found with ID %s", drone.ID)
	}
	return nil
}
func GetDroneCurrentOrder(droneID string) (model.Orders, error) {
	var order model.Orders
	query := ` select o.id, o.origin, o.destination, o.status from orders o
	JOIN drones d ON o.id = d.current_order_id
	where d.id =?  `
	err := database.DB.QueryRow(query, droneID).Scan(&order.ID, &order.Origin, &order.Destination, &order.Status)
	if err != nil {
		return model.Orders{}, err
	}
	return order, nil
}
