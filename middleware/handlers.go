package middleware

import (
	"encoding/json"
	"net/http"
	"strconv"

	"GODAPP/models"

	"github.com/gorilla/mux"

	_ "github.com/lib/pq" //TODO: change to another driver and use sqlx
)

func GetTableKeys(w http.ResponseWriter, r *http.Request) {
	//w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	params := mux.Vars(r)
	data := getTableKeys(params["table"])
	json.NewEncoder(w).Encode(data)
}

func GetSelectedData(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	reqData := models.SelectModel{}
	err := json.NewDecoder(r.Body).Decode(&reqData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	data := getSelectData(reqData)
	w.Write([]byte(data))
}

func GetSelectedDataWithPagination(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	params := mux.Vars(r)
	limit, _ := strconv.Atoi(params["limit"])
	offset, _ := strconv.Atoi(params["offset"])
	reqData := models.SelectModel{}
	err := json.NewDecoder(r.Body).Decode(&reqData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	response := getPagedSelectData(reqData, r.Host, limit, offset)
	w.Write([]byte(response))
}

func MergeJSON(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/text")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	reqData := models.MergeModel{}
	err := json.NewDecoder(r.Body).Decode(&reqData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	/*affected := mergeSelectedData(reqData)
	response, err := json.Marshal(affected)*/
	query := GetMergeQuery(reqData)
	w.Write([]byte(query))
}
