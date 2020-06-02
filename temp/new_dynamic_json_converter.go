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
}
