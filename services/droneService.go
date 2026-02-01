package services

import (
	"drone/database"
	model "drone/models"
	"log"
)

func GetDroneByID(droneID string) (model.Drone, error) {
	query := `select * from drones where id = ?`
	drone := model.Drone{}
	err := database.DB.QueryRow(query, droneID).Scan(&drone.ID, &drone.Status, &drone.Lat, &drone.Lng, &drone.OrderID)
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

func GetDroneStatus(droneID string) (string, error) {
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
	if status == "fixed" {
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
	_, err = database.DB.Exec(query, "pickup", orderID)
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
	values ( "fixed", 0.0, 0.0, null)`

	_, err := database.DB.Exec(query)
	if err != nil {
		log.Println("Failed to create drone:", err)
		return err
	}

	return nil
}
func TryAssignPendingOrder() error {
	// 1) find pending order
	var orderID int64
	err := database.DB.QueryRow(`
		SELECT id FROM orders 
		WHERE status = 'pending' AND assigned_drone_id IS NULL
		ORDER BY id ASC
		LIMIT 1
	`).Scan(&orderID)

	if err != nil {
		// no pending orders
		return nil
	}

	// 2) find available drone
	var droneID int64
	err = database.DB.QueryRow(`
		SELECT id FROM drones 
		WHERE status = 'fixed' AND current_order_id IS NULL
		LIMIT 1
	`).Scan(&droneID)

	if err != nil {
		// no available drone
		return nil
	}

	// 3) assign
	_, err = database.DB.Exec(`
		UPDATE orders 
		SET assigned_drone_id = ?, status = 'assigned'
		WHERE id = ?
	`, droneID, orderID)
	if err != nil {
		return err
	}

	_, err = database.DB.Exec(`
		UPDATE drones 
		SET current_order_id = ?, status = 'busy'
		WHERE id = ?
	`, orderID, droneID)

	return err
}
