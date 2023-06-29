package main

import (
	"bookstore/internal/orders/domain"
	"encoding/json"
	"fmt"
	"log"
)

func main() {
	o := domain.Order{
		User: "lol",
		Contents: []domain.Book{
			{"id"},
			{"id"},
			{"id"},
			{"id"},
			{"id"},
		},
	}

	b, err := json.Marshal(o)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(b))
}
