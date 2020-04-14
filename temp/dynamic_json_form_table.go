package temp

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"log"

	_ "github.com/lib/pq"
)

type Item struct {
	ID    int
	Attrs Attrs
}

type Attrs map[string]interface{}

func (a Attrs) Value() (driver.Value, error) {
	return json.Marshal(a)
}

func (a *Attrs) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &a)
}

func main() {
	db, err := sql.Open("postgres", "postgres://user:pass@localhost/db")
	if err != nil {
		log.Fatal(err)
	}

	item := new(Item)
	item.Attrs = Attrs{
		"name":        "Passata",
		"ingredients": []string{"Tomatoes", "Onion", "Olive oil", "Garlic"},
		"organic":     true,
		"dimensions": map[string]interface{}{
			"weight": 250.00,
		},
	}

	_, err = db.Exec("INSERT INTO items (attrs) VALUES($1)", item.Attrs)
	if err != nil {
		log.Fatal(err)
	}

	item = new(Item)
	err = db.QueryRow("SELECT id, attrs FROM items ORDER BY id DESC LIMIT 1").Scan(&item.ID, &item.Attrs)
	if err != nil {
		log.Fatal(err)
	}

	name, ok := item.Attrs["name"].(string)
	if !ok {
		log.Fatal("unexpected type for name")
	}
	dimensions, ok := item.Attrs["dimensions"].(map[string]interface{})
	if !ok {
		log.Fatal("unexpected type for dimensions")
	}
	weight, ok := dimensions["weight"].(float64)
	if !ok {
		log.Fatal("unexpected type for weight")
	}
	weightKg := weight / 1000
	log.Printf("%s: %.2fkg", name, weightKg)
}
