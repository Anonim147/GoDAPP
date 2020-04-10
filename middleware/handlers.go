package middleware

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	_ "github.com/lib/pq" //TODO: change to another driver and use sqlx
)

func GetTableKeys(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	params := mux.Vars(r)
	data := getTableKeys(params["table"])

	/*if err != nil {
		log.Fatalf("Unable to get all user. %v", err)
	}*/

	json.NewEncoder(w).Encode(data)
}
