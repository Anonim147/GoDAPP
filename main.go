package main

import (
	"GODAPP/router"
	"fmt"
	"log"
	"net/http"
)

func main() {
	r := router.Router()
	fmt.Println("Start listening on 9000...")
	log.Fatal(http.ListenAndServe(":9000", r))
}
