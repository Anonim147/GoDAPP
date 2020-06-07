package middleware

import (
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"GODAPP/models"

	"github.com/jackc/pgx"
	"github.com/joho/godotenv"
)

func createConnection() (*pgx.Conn, error) {
	err := godotenv.Load(".env")
	if err != nil {
		return nil, err
	}

	connConfig, err := pgx.ParseURI(os.Getenv("POSTGRES_URL"))
	if err != nil {
		return nil, err
	}

	conn, err := pgx.Connect(connConfig)

	if err != nil {
		return nil, err
	}

	return conn, nil
}

func createTable(db *pgx.Conn, tableName string) error {
	query := GetQueryForCreateTable(tableName)
	_, err := db.Exec(query)
	if err != nil {
		return err
	}
	return nil
}

func getTableKeys(tableName string) ([]models.TableKey, error) {
	db, err := createConnection()
	if err != nil {
		return nil, err
	}
	defer db.Close()
	query := GetJSONKeysQuery(tableName)
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	data := []models.TableKey{}
	for rows.Next() {
		var tk models.TableKey
		err = rows.Scan(&tk.KeyName, &tk.KeyType)
		if err != nil {
			return nil, err
		}
		data = append(data, tk)
	}

	data = TransformKeys(data)
	return data, nil
}

func getSelectData(data models.SelectModel) string {
	db, err := createConnection()
	defer db.Close()
	query := GetSelectQuery(data, 0, 0)
	rows, err := db.Query(query)
	if err != nil {
		panic(err)
	}
	queryText := `{"data" :[`
	for rows.Next() {
		var data []byte
		err = rows.Scan(&data)
		if err != nil {
			panic(err)
		}
		queryText += string(data) + ","
	}
	queryText = queryText[:len(queryText)-1] + "]}"
	return queryText
}

func getPagedSelectData(data models.SelectModel, host string, limit int, offset int) string {
	db, err := createConnection()
	if err != nil {
		return ""
	}
	defer db.Close()

	countOfRows := getDataCount(db, data)

	pag := models.Pagination{}
	pag.SelfLink = getLinkForPagination(host, limit, offset)
	if offset == 0 {
		pag.PrevLink = ""
	} else {
		pag.PrevLink = getLinkForPagination(host, limit, offset-limit)
	}
	if countOfRows > limit+offset {
		pag.NextLink = getLinkForPagination(host, limit, offset+limit)
	} else {
		pag.NextLink = ""
	}
	pagData, _ := json.Marshal(pag)

	query := GetSelectQuery(data, limit, offset)
	rows, err := db.Query(query)
	if err != nil {
		panic(err)
	}

	queryText := `{"data" :[ `
	for rows.Next() {
		var data []byte
		err = rows.Scan(&data)
		if err != nil {
			panic(err)
		}
		queryText += string(data) + ","
	}
	queryText = queryText[:len(queryText)-1] + "]," + "\n \"pagination\" :" + string(pagData) + "}"
	return queryText
}

func getDataCount(db *pgx.Conn, data models.SelectModel) int { //TODO: err handling
	query := GetCountQuery(data)
	var count int
	row := db.QueryRow(query)
	err := row.Scan(&count)
	if err != nil {
		log.Fatal(err)
	}

	return count
}

func getAllDataCount(db *pgx.Conn, tableName string) (int, error) {
	query := GetAllDataCountQuery(tableName)
	var count int
	row := db.QueryRow(query)
	err := row.Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func getMD5(data string) string {
	byteData := []byte(data)
	return fmt.Sprintf("%x", md5.Sum(byteData))
}

func insertIntoTable(filePath string, tablename string, tableExists bool) (string, error) {

	var data []byte
	var err error
	if filepath.Ext(filePath) == `.json` {
		data, err = getDataFromFile(filePath)
		if err != nil {
			return "", err
		}
	} else if filepath.Ext(filePath) == `.csv` {
		data, err = FormatCSVtoJSON(filePath)
		if err != nil {
			return "", err
		}
	} else {
		return "", errors.New("Unrecognized file extention")
	}

	db, err := createConnection()
	if err != nil {
		return "", err
	}
	defer db.Close()

	if !tableExists {
		err = createTable(db, tablename)
		if err != nil {
			return "", err
		}
	}

	parseQuery := GetQueryForParseJSON(tablename, string(data))
	index := strings.Index(string(data), `[`) < strings.Index(string(data), `{`)
	if index {
		parseQuery = GetQueryForParseJSONARRAY(tablename, string(data))
	}
	res, err := db.Exec(parseQuery)
	if err != nil {
		return "", err
	}

	err = os.Remove(filePath)
	if err != nil {
		return "", err
	}

	affected := res.RowsAffected()
	result := fmt.Sprintf("Inserted %s rows", strconv.FormatInt(affected, 10))

	return result, nil
}

func updateTable(update models.UpdateModel) (string, error) {
	//To do: temporary file and temporary table to .env
	tempTable := getRandomString()
	_, err := insertIntoTable(update.FilePath, tempTable, false)

	db, err := createConnection()
	if err != nil {
		return "", err
	}
	defer db.Close()

	fmt.Println(GetQueryForUpdateTable(update.TableName, tempTable, update.Columns))

	res, err := db.Exec(GetQueryForUpdateTable(update.TableName, tempTable, update.Columns))
	if err != nil {
		return "", err
	}

	err = dropTable(db, tempTable)
	if err != nil {
		return "", err
	}

	affected := res.RowsAffected()

	result := fmt.Sprintf("Updated %s rows", strconv.FormatInt(affected, 10))
	return result, nil
}

func replaceTable(update models.UpdateModel) (string, error) {
	db, err := createConnection()
	if err != nil {
		return "", err
	}
	defer db.Close()

	_, err = db.Exec(GetQueryForClearTable(update.TableName))
	if err != nil {
		return "", err
	}

	res, err := insertIntoTable(update.FilePath, update.TableName, true)
	if err != nil {
		return "", err
	}

	return res, nil
}

func getTableList() ([]string, error) {
	db, err := createConnection()
	if err != nil {
		return []string{}, err
	}
	defer db.Close()
	rows, err := db.Query(GetQueryForTableList())
	if err != nil {
		return []string{}, err
	}
	defer rows.Close()

	var tableNames []string

	for rows.Next() {
		var tableName string
		err := rows.Scan(&tableName)
		if err != nil {
			return []string{}, err
		}

		tableNames = append(tableNames, tableName)
	}

	return tableNames, nil
}

func dropTable(db *pgx.Conn, tableName string) error {
	_, err := db.Exec(GetQueryForDropTable(tableName))
	if err != nil {
		return err
	}
	return nil
}

//TO DO : додати експорт в json

func clearTable(db *pgx.Conn, tablename string) error {
	_, err := db.Exec(GetQueryForClearTable(tablename))
	if err != nil {
		return err
	}
	return nil
}

func getTableInfo(tablename string) (models.DBInfo, error) {
	var dbinfo models.DBInfo
	db, err := createConnection()
	if err != nil {
		return dbinfo, err
	}

	count, err := getAllDataCount(db, tablename)
	if err != nil {
		return dbinfo, err
	}

	dbinfo.TableName = tablename
	dbinfo.RecordsCount = count

	return dbinfo, nil
}
