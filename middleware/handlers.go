package middleware

import (
	"encoding/json"
	"fmt"
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

func MergeJSON(w http.ResponseWriter, r *http.Request) { // TO DO: прочекати
	w.Header().Set("Content-Type", "application/text")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	reqData := models.MergeModel{}
	err := json.NewDecoder(r.Body).Decode(&reqData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	affected := mergeSelectedData(reqData)
	resData := fmt.Sprintf(`{"affected" : %d}`, affected)
	w.Write([]byte(resData))
}

func UploadTable(w http.ResponseWriter, r *http.Request) { // TO DO: прочекати
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	w.Header().Set("Content-Type", "application/json")

	//w.Header().Set("Content-Type", "application/json")
	if r.Method == "POST" {
		src, hdr, err := r.FormFile("file")
		if err != nil {
			w.Write([]byte(err.Error()))
		}
		defer src.Close()

		path := filepath.Join(os.TempDir(), hdr.Filename)
		dst, err := os.Create(path)
		if err != nil {
			//http.Error(w, err.Error(), 500)
			w.Write([]byte(err.Error()))
		}
		defer dst.Close()

		io.Copy(dst, src)
		resData := models.BaseResponse{
			Success: true,
			Value:   path,
			Error:   "",
		}
		response, _ := json.Marshal(resData)
		w.Write([]byte(response))
	}
}

func ImportToNewTable(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	if r.Method == "POST" {
		reqData := models.InsertTableModel{}
		err := json.NewDecoder(r.Body).Decode(&reqData)
		if err != nil {
			fmt.Println(reqData)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		affected, err := insertJSONIntoTable(reqData.FilePath, reqData.TableName)
		fmt.Print("done")
		fmt.Println(err)
		resData := models.BaseResponse{
			Success: err == nil,
			Value:   affected,
			Error:   "",
		}
		response, _ := json.Marshal(resData)
		w.Write(response)
	}
}

func UpdateTable(w http.ResponseWriter, r *http.Request) { // TO DO: зробити норм модельку на вхід і на вихід
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	if r.Method == "POST" {
		reqData := models.InsertTableModel{}
		err := json.NewDecoder(r.Body).Decode(&reqData)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		affected, err := updateJSONIntoTable(reqData.FilePath, reqData.TableName)
		errtext := ""
		if err != nil {
			errtext = err.Error()
		}
		resData := models.BaseResponse{
			Success: err == nil,
			Value:   affected,
			Error:   errtext,
		}
		response, _ := json.Marshal(resData)
		w.Write(response)
	}
}

func GetTableList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	m := map[string]interface{}{}
	tableList, _ := getTableList()
	m["tables"] = tableList
	data, _ := json.Marshal(m)
	w.Write(data)
}
