package router

import (
	"GODAPP/middleware"

	"github.com/gorilla/mux"
)

func Router() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/api/get_columns/{table}", middleware.GetTableKeys).Methods("GET", "OPTIONS")
	//router.HandleFunc("/api/get_data", middleware.GetRows).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/get_data", middleware.GetSelectedData).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/get_data&limit={limit}&offset={offset}", middleware.GetSelectedDataWithPagination).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/merge_data", middleware.MergeJSON).Methods("POST", "OPTIONS")

	/*router.HandleFunc("/api/user/", middleware.GetAllUser).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/newuser", middleware.CreateUser).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/user/{id}", middleware.UpdateUser).Methods("PUT", "OPTIONS")
	router.HandleFunc("/api/deleteuser/{id}", middleware.DeleteUser).Methods("DELETE", "OPTIONS")
	*/
	return router
}
