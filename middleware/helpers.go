package middleware

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"GODAPP/models"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq" //TODO: change to another driver and use sqlx
)

func createConnection() *sql.DB {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	db, err := sql.Open("postgres", os.Getenv("POSTGRES_URL"))

	if err != nil {
		panic(err)
	}

	return db
}

func createTable(tableName string, db *sql.DB) {
	query := `CREATE TABLE ` + tableName + ` (
		id int GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
		data jsonb NOT NULL
	)`
	_, err := db.Exec(query)
	if err != nil {
		panic(err)
	}
}

func copyToTable(path string, tableName string, db *sql.DB) int64 {

	query := `copy ` + tableName + `(data) from ` + path
	res, err := db.Exec(query)
	if err != nil {
		panic(err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		panic(err)
	}
	return rows
}

func downloadFile(url string, filepath string) error {

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}
	return nil
}

func getTableKeys(tableName string) map[string]string {
	db := createConnection()
	query := `
	WITH RECURSIVE doc_key_and_value_recursive(key, value) AS (
		SELECT
		  t.key,
		  t.value
		  FROM ` + tableName + `, jsonb_each(` + tableName + `.data) AS t
	  
		UNION ALL
	  
		SELECT
		  CONCAT(doc_key_and_value_recursive.key, '.', t.key),
		  t.value
		FROM doc_key_and_value_recursive,
		  jsonb_each(
			CASE 
			  WHEN jsonb_typeof(doc_key_and_value_recursive.value) <> 'object' THEN '{}' :: JSONB
			  ELSE doc_key_and_value_recursive.value
			END
			) AS t
	  )
	  SELECT DISTINCT key as key, jsonb_typeof(value) as valuetype
	  FROM doc_key_and_value_recursive
	  WHERE jsonb_typeof(doc_key_and_value_recursive.value) NOT IN ( 'object')   --'array',
	  ORDER BY key`
	rows, err := db.Query(query)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	data := map[string]string{}
	for rows.Next() {
		var column, columntype string
		err = rows.Scan(&column, &columntype)
		if err != nil {
			panic(err)
		}
		data[column] = columntype
	}
	return data
}

func mapLogicalype(data models.SelectCondition) string {
	switch data.LogicalType {
	case "and":
		return " AND "
	default:
		return " OR "
	}
}

func mapCondition(data models.SelectCondition) string {

	switch data.ComparisonType {
	case "lt":
		return fmt.Sprintf(` data #> '{ %s }' < '%s' `, strings.Replace(data.ColumnPath, ".", ",", -1), data.Value)
	case "let":
		return fmt.Sprintf(` data #> '{ %s }' <= '%s' `, strings.Replace(data.ColumnPath, ".", ",", -1), data.Value)
	case "gt":
		return fmt.Sprintf(` data #> '{ %s }' >= '%s' `, strings.Replace(data.ColumnPath, ".", ",", -1), data.Value)
	case "get":
		return fmt.Sprintf(` data #> '{ %s }' <= '%s' `, strings.Replace(data.ColumnPath, ".", ",", -1), data.Value)
	case "like":
		return fmt.Sprintf(` data #> '{ %s }' like '%%%s%%' `, strings.Replace(data.ColumnPath, ".", ",", -1), data.Value)
	case "eq":
		return fmt.Sprintf(` data #> '{ %s }' = '%s' `, strings.Replace(data.ColumnPath, ".", ",", -1), data.Value)
	default:
		return fmt.Sprintf(` data #> '{ %s }' = '%s' `, strings.Replace(data.ColumnPath, ".", ",", -1), data.Value)
	}
}

func getFilters(data models.SelectModel) string {
	query := " WHERE "

	for i, condition := range data.Conditions {
		if i != 0 {
			query += condition.LogicalType
		}
		query += mapCondition(condition)
	}
	return query
}

func getColumns(data models.SelectModel) string {
	query := `SELECT `
	for _, s := range data.Columns {
		//todo: do another symbol
		query += fmt.Sprintf(` data #> '{ %s}' as "%s", `, strings.Replace(s, ".", ",", -1), s)
	}
	query = query[:len(query)-1]
	query += ` FROM ` + data.TableName

	return query
}

func getQuery(data models.SelectModel) string {
	query := getColumns(data)
	if len(data.Conditions) == 0 {
		return query
	}
	query += getFilters(data)
	//get limit? where?
	return query
}
