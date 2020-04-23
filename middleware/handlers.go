package middleware

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
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

func MergeJSON(w http.ResponseWriter, r *http.Request) { // TO DO: доробити до кінця цю шнягу
	w.Header().Set("Content-Type", "application/text") // TO DO: поміняти на json
	w.Header().Set("Access-Control-Allow-Origin", "*")

	m := map[string]interface{}{}
	reqData := models.MergeModel{}
	err := json.NewDecoder(r.Body).Decode(&reqData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	m["affected_rows"] = mergeSelectedData(reqData)
	jsonResp, _ := json.Marshal(m)
	w.Write([]byte(jsonResp))
}

func UploadTable(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	m := map[string]string{}
	if r.Method == "POST" {
		src, hdr, err := r.FormFile("file")
		if err != nil {
			http.Error(w, err.Error(), 500)
		}
		defer src.Close()

		path := filepath.Join(os.TempDir(), hdr.Filename)
		dst, err := os.Create(path)
		if err != nil {
			http.Error(w, err.Error(), 500)
		}
		defer dst.Close()

		io.Copy(dst, src)
		m["path"] = path
		json, _ := json.Marshal(m)
		w.Write([]byte(json))
	}
}
