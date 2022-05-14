package router

import (
	"ToDoApp/controller"

	"github.com/gorilla/mux"
)

func Router() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/api/tasks", controller.GetAllTask).Methods("GET")
	router.HandleFunc("/api/task", controller.CreateTask).Methods("POST")
	router.HandleFunc("/api/task/{id}", controller.TaskCompleted).Methods("PUT")
	router.HandleFunc("/api/undoTask/{id}", controller.UndoTask).Methods("PUT")
	router.HandleFunc("/api/deleteTask/{id}", controller.DeleteOneTask).Methods("DELETE")
	router.HandleFunc("/api/deleteAll", controller.DeleteAllTask).Methods("DELETE")
	return router
}
