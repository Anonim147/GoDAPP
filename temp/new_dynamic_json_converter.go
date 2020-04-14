package temp

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
)

func convert() {
	conn := "user=postgres password=123456 dbname=sandbox sslmode=disable"

	db, err := sql.Open("postgres", conn)
	defer db.Close()
	if err != nil {
		panic(err)
	}
	query := `SELECT data #> '{ organic}' as "organic", data #> '{ dimensions}' as "dimensions", 
	data #> '{ dimensions,weight}' as "dimensions.weight", data #> '{ ingredients}' as "ingredients" 
	FROM tempo WHERE  data #>> '{ dimensions,weight }' like '%50%'  OR  data #>> '{ name }' = 'Pizza'`
	rows, err := db.Query(query)
	if err != nil {
		panic(err)
	}
	columns, err := rows.Columns()
	if err != nil {
		panic(err)
	}

	count := len(columns)
	values := make([]interface{}, count)
	scanArgs := make([]interface{}, count)
	for i := range values {
		scanArgs[i] = &values[i]
	}

	masterData := make(map[string][]interface{})

	for rows.Next() {
		err := rows.Scan(scanArgs...)
		if err != nil {
			panic(err)
		}
		for i, v := range values {

			x := v.([]byte)

			//NOTE: FROM THE GO BLOG: JSON and GO - 25 Jan 2011:
			// The json package uses map[string]interface{} and []interface{} values to store arbitrary JSON objects and arrays; it will happily unmarshal any valid JSON blob into a plain interface{} value. The default concrete Go types are:
			//
			// bool for JSON booleans,
			// float64 for JSON numbers,
			// string for JSON strings, and
			// nil for JSON null.

			if nx, ok := strconv.ParseFloat(string(x), 64); ok == nil {
				masterData[columns[i]] = append(masterData[columns[i]], nx)
			} else if b, ok := strconv.ParseBool(string(x)); ok == nil {
				masterData[columns[i]] = append(masterData[columns[i]], b)
			} else if "string" == fmt.Sprintf("%T", string(x)) {
				masterData[columns[i]] = append(masterData[columns[i]], string(x))
			} else {
				fmt.Printf("Failed on if for type %T of %v\n", x, x)
			}

		}
	}

	jsonData, err := json.Marshal(masterData)
	fmt.Println(string(jsonData))
}
