package main

import (
	"drone/database"
	"drone/handlers"
	"fmt"
	"log"
	"net/http"
)

func main() {

	// init database
	database.InitDB()

	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
	mux.HandleFunc("/get_orderByID", handlers.GetOrderByID)	
	mux.HandleFunc("/get_orders", handlers.GetOrders)

	mux.HandleFunc("/create_order", handlers.CreateOrder)
	mux.HandleFunc("/order_status", handlers.GetOrderStatus)
	mux.HandleFunc("/cancel_order", handlers.CancelOrder)
	//Admin
mux.HandleFunc("/create_drone", handlers.CreateDrone)
	mux.HandleFunc("/update_origin", handlers.UpdateOrderOrigin)
	mux.HandleFunc("/update_destination", handlers.UpdateOrderDestination)
	mux.HandleFunc("/update_drone_status", handlers.UpdateDroneStatus)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	log.Println("Server starting on :8080")

	err := server.ListenAndServe()
	if err != nil {
		fmt.Println("server error:", err)
	}
}
