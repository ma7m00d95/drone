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
	mux.HandleFunc("GET /get_orderByID/{order_id}", handlers.GetOrderByID)	
	mux.HandleFunc("GET /get_orders", handlers.GetOrders)

	mux.HandleFunc("POST /create_order", handlers.CreateOrder)
	mux.HandleFunc("GET /order_status/{order_id}", handlers.GetOrderStatus)
	mux.HandleFunc("GET /cancel_order/{order_id}", handlers.CancelOrder)
	mux.HandleFunc("GET /update_order/{order_id}/{status}", handlers.UpdateOrder)
	//Admin
mux.HandleFunc("POST /create_drone", handlers.CreateDrone)
	mux.HandleFunc("POST /update_origin", handlers.UpdateOrderOrigin)
	mux.HandleFunc("POST /update_destination", handlers.UpdateOrderDestination)
	mux.HandleFunc("POST /update_drone_status", handlers.UpdateDroneStatus)

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
