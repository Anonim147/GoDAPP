package temp

import (
	"database/sql"
)

func createJSON() {
	conn := "user=postgres password=123456 dbname=sandbox sslmode=disable"

	db, err := sql.Open("postgres", conn)
	defer db.Close()
	if err != nil {
		panic(err)
	}
	query := `select row_to_json(d) as data
    from 
    (SELECT data #> '{ ingredients }' as "ingredients", 
            data #> '{ dimensions }' as "dimensions", 
            data #> '{ name}' as "name" 
                FROM tempo 
                WHERE  data #>> '{ dimensions,weight }' like '%50%'  
                    OR  data #>> '{ name }' = 'Pizza') d `
	rows, err := db.Query(query)
	if err != nil {
		panic(err)
	}
	jsonData := `[ `
	for rows.Next() {
		var data []byte
		err = rows.Scan(&data)
		if err != nil {
			panic(err)
		}
		jsonData += string(data) + ","
	}
	jsonData = jsonData[:len(jsonData)-1] + "]"
}
