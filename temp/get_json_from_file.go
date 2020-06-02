package main

import (
	"bufio"
	"encoding/json"
	"log"
	"os"
	"strings"

	_ "github.com/lib/pq"
)

func main() {
	var result map[string]interface{}
	var jsondtata []string
	file, err := os.Open("data.json")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		row := scanner.Text()
		first := strings.Index(row, "{")
		last := strings.LastIndex(row, "}")
		if first > 0 && last > first {
			row = row[first : last+1]
			err = json.Unmarshal([]byte(row), &result)
			if err == nil {
				jsondtata = append(jsondtata, string(row))
			}
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
