package main

import (
	"drone/database"
	"drone/handlers"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

func main() {

	// init database
	database.InitDB()
	// 1. Load the specific config.env file
	err := godotenv.Load("config.env")
	if err != nil {
		log.Fatal("Error loading config.env file")
	}

	// 2. Now you can grab the secret
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET not found in config.env")
	}

	mux := http.NewServeMux()
	publicMux := http.NewServeMux()

	publicMux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
	publicMux.HandleFunc("POST /login", handlers.Login) // 👈 ADD THIS LINE HERE

	mux.HandleFunc("GET /get_orderByID/{order_id}", handlers.GetOrderByID)
	mux.HandleFunc("GET /get_orders", handlers.GetOrders)
	mux.HandleFunc("GET /get_inProcess_orders", handlers.GetInProcessOrders)

	mux.HandleFunc("GET /get_user_orders/{user_id}", handlers.GetUserOrders)

	mux.HandleFunc("POST /create_order", handlers.CreateOrder)
	mux.HandleFunc("GET /order_status/{order_id}", handlers.GetOrderStatus)
	mux.HandleFunc("GET /cancel_order/{order_id}", handlers.CancelOrder)
	mux.HandleFunc("GET /update_order/{order_id}/{status}", handlers.UpdateOrder)
	mux.HandleFunc("POST /orders/deliver", handlers.DeliverOrder)
	mux.HandleFunc("POST /orders/fail", handlers.FailOrder)
	mux.HandleFunc("POST /drones/location", handlers.UpdateDroneLocation)

	//Admin
	mux.HandleFunc("GET /create_drone", handlers.CreateDrone)
	mux.HandleFunc("POST /update_origin", handlers.UpdateOrderOrigin)
	mux.HandleFunc("POST /update_destination", handlers.UpdateOrderDestination)
	mux.HandleFunc("GET /update_drone_status/{drone_id}/{status}", handlers.UpdateDroneStatus)

	//drone
	mux.HandleFunc("POST /drones/reserve", handlers.ReserveOrderToDrone)
	mux.HandleFunc("POST /drones/pickup", handlers.DronePickupOrder)
	mux.HandleFunc("GET /drones/current-order/{drone_id}", handlers.GetDroneCurrentOrder)

	rootMux := http.NewServeMux()
	rootMux.Handle("/health", publicMux)

	rootMux.Handle("/login", publicMux)

	rootMux.Handle("/",
		AuthMiddleware(mux, jwtSecret),
	)

	server := &http.Server{
		Addr:    ":8080",
		Handler: rootMux,
	}
	log.Println("Server starting on :8080")

	err = server.ListenAndServe()
	if err != nil {
		fmt.Println("server error:", err)
	}
}

func AuthMiddleware(next http.Handler, secret string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Allow health without auth
		if r.URL.Path == "/health" {
			next.ServeHTTP(w, r)
			return
		}

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		claims, err := handlers.ValidateJWT(tokenString)
		if err != nil {
			http.Error(w, "Invalid token: "+err.Error(), http.StatusUnauthorized)
			return
		}

		// optional context injection later
		log.Printf("Authenticated user: %s | role: %s\n", claims.UserID, claims.Role)

		next.ServeHTTP(w, r)
	})
}
