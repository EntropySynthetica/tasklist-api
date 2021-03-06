package router

import (
	"tasklist-api/middleware"

	"github.com/gorilla/mux"
)

// Router is exported and used in main.go
func Router() *mux.Router {

	router := mux.NewRouter()

	router.HandleFunc("/api/user/{id}", middleware.GetUser).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/user", middleware.GetAllUser).Methods("GET", "OPTIONS")

	router.HandleFunc("/api/task/{id}", middleware.GetTask).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/task/status/{id}", middleware.GetTaskStatus).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/task", middleware.GetAllTasks).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/task/{id}", middleware.UpdateTask).Methods("PUT", "OPTIONS")
	router.HandleFunc("/api/newtask", middleware.CreateTask).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/deletetask/{id}", middleware.DeleteTask).Methods("DELETE", "OPTIONS")

	return router
}
