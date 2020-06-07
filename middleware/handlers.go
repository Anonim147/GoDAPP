package middleware

import (
	"GODAPP/models"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gorilla/mux"
)

func GetTableKeys(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	data, err := getTableKeys(params["table"])
	response := models.BaseResponse{
		Success: true,
		Value:   data,
	}
	if err != nil {
		response = models.BaseResponse{
			Success: false,
			Value:   err.Error,
		}
	}
	json.NewEncoder(w).Encode(response)
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
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	if r.Method == "POST" {
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
}

func UploadFile(w http.ResponseWriter, r *http.Request) {
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
		}
		response, _ := json.Marshal(resData)
		w.Write([]byte(response))
	}
}

func DownloadFileFromLink(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	w.Header().Set("Content-Type", "application/json")

	if r.Method == "POST" {
		var reqData models.DownloadModel
		err := json.NewDecoder(r.Body).Decode(&reqData)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		value, err := downloadFileFromLink(reqData.FilePath)
		if err != nil {
			value = err.Error()
		}
		resData := models.BaseResponse{
			Success: err == nil,
			Value:   value,
		}
		json.NewEncoder(w).Encode(resData)
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
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		affected, err := insertJSONIntoTable(reqData.FilePath, reqData.TableName, false)
		value := affected
		if err != nil {
			value = err.Error()
		}
		resData := models.BaseResponse{
			Success: err == nil,
			Value:   value,
		}
		json.NewEncoder(w).Encode(resData)
	}
}

func GetTableList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	m := map[string]interface{}{}
	tableList, _ := getTableList()
	m["tables"] = tableList
	json.NewEncoder(w).Encode(m)
}

func UpdateTable(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	w.Header().Set("Content-Type", "application/json")

	if r.Method == "POST" {
		fmt.Println("do it")
		var reqData models.UpdateModel
		err := json.NewDecoder(r.Body).Decode(&reqData)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		var value string
		if reqData.Method == "full replace" {
			fmt.Println("replace")
			value, err = replaceTable(reqData)
		} else {
			value, err = updateTable(reqData)
		}
		if err != nil {
			value = err.Error()
		}

		resData := models.BaseResponse{
			Success: err == nil,
			Value:   value,
		}

		json.NewEncoder(w).Encode(resData)
	}
}

func DropTable(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	tablename := params["tablename"]

	db, err := createConnection()
	if err != nil {
		res := models.BaseResponse{
			Success: false,
			Value:   err.Error(),
		}
		json.NewEncoder(w).Encode(res)
		return
	}
	err = dropTable(db, tablename)
	if err != nil {
		res := models.BaseResponse{
			Success: false,
			Value:   err.Error(),
		}
		json.NewEncoder(w).Encode(res)
		return
	}

	res := models.BaseResponse{
		Success: true,
		Value:   "ok",
	}
	json.NewEncoder(w).Encode(res)
}

func ClearTable(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	tablename := params["tablename"]

	db, err := createConnection()
	if err != nil {
		res := models.BaseResponse{
			Success: false,
			Value:   err.Error(),
		}
		json.NewEncoder(w).Encode(res)
		return
	}
	err = clearTable(db, tablename)
	if err != nil {
		res := models.BaseResponse{
			Success: false,
			Value:   err.Error(),
		}
		json.NewEncoder(w).Encode(res)
		return
	}

	res := models.BaseResponse{
		Success: true,
		Value:   "ok",
	}
	json.NewEncoder(w).Encode(res)
}

func GetTableInfo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	tablename := params["tablename"]

	resp := models.BaseResponse{}
	info, err := getTableInfo(tablename)
	if err != nil {
		resp.Success = false
		resp.Value = err.Error()
	} else {
		resp.Success = true
		resp.Value = info
	}
	json.NewEncoder(w).Encode(resp)
}
