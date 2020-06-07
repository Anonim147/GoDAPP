package router

import (
	"GODAPP/middleware"

	"github.com/gorilla/mux"
)

func Router() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/api/get_columns/{table}", middleware.GetTableKeys).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/get_table_list", middleware.GetTableList).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/get_table_info/{tablename}", middleware.GetTableInfo).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/get_data", middleware.GetSelectedData).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/get_data&limit={limit}&offset={offset}", middleware.GetSelectedDataWithPagination).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/upload", middleware.UploadFile).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/insert_data", middleware.ImportToTable).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/download_from_link", middleware.DownloadFileFromLink).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/update_table", middleware.UpdateTable).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/drop_table/{tablename}", middleware.DropTable).Methods("DELETE", "OPTIONS")
	router.HandleFunc("/api/clear_table/{tablename}", middleware.ClearTable).Methods("DELETE", "OPTIONS")

	return router
}
