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
	_, err = database.DB.Exec(query, "stopped", orderID)
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

// Reserve a job.
// ● Grab an order from a location (origin or broken drone).
// ● Mark an order they have gotten as delivered or failed.
// ● Mark themselves as broken (and in need of an order handoff).
// ● Update their location (use latitude/longitude), and get a status update as a heartbeat.
// ● Get details on the order they are currently assigned.
