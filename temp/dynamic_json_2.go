package temp

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"reflect"
)

func rowsToJson() ([]byte, error) {
	conn := "user=postgres password=123456 dbname=sandbox sslmode=disable"

	db, err := sql.Open("postgres", conn)
	defer db.Close()
	if err != nil {
		panic(err)
	}

	d1 := []SelectCondition{
		SelectCondition{
			ColumnPath:     "dimensions.weight",
			ComparisonType: "like",
			Value:          "50",
		}, SelectCondition{
			ColumnPath:     "name",
			ComparisonType: "eq",
			Value:          "Pizza",
		},
	}

	data := SelectModel{
		TableName:  "tempo",
		Columns:    []string{"organic", "dimensions.weight", "name"},
		Conditions: d1,
	}

	var objects []map[string]interface{}
	query := GetQuery(data)
	fmt.Println(query)
	rows, err := db.Query(query)
	if err != nil {
		fmt.Println("fuck")
		panic(err)
	}

	for rows.Next() {
		columns, err := rows.ColumnTypes()
		if err != nil {
			print("fuck2")
			panic(err)
		}
		values := make([]interface{}, len(columns))
		object := map[string]interface{}{}
		for i, column := range columns {
			v := reflect.New(column.ScanType()).Interface()
			switch v.(type) {
			case *[]uint8:
				v = new(string)
			default:
				// use this to find the type for the field
				// you need to change
				// log.Printf("%v: %T", column.Name(), v)
			}

			object[column.Name()] = v
			values[i] = object[column.Name()]
		}
		err = rows.Scan(values...)
		if err != nil {
			panic(err)
		}
		//fmt.Println(object["name"])
		objects = append(objects, object)

	}

	return json.MarshalIndent(objects, "", "\t")
}
