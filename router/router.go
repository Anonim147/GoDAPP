package router

import (
	"GODAPP/middleware"

	"github.com/gorilla/mux"
)

func Router() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/api/get_columns/{table}", middleware.GetTableKeys).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/get_data", middleware.GetSelectedData).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/get_data&limit={limit}&offset={offset}", middleware.GetSelectedDataWithPagination).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/upload", middleware.UploadTable).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/insert_data", middleware.ImportToNewTable).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/get_table_list", middleware.GetTableList).Methods("GET", "OPTIONS")

	return router
}
