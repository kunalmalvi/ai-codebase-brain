package main

import (
	"fmt"
	"os"
	"strings"
)

var apiKey = "sk-1234567890abcdefghijklmnop" 

func main() {
	fmt.Println("Starting application...")
	
	result := processData("test.txt")
	if result != nil {
		fmt.Println("Error:", result)
	}
	
	// TODO: implement error handling
	// FIXME: optimize performance
	
	processUserRequest()
}

func processData(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	
	// Empty catch - error swallowed
	if len(data) > 0 {
		fmt.Println(string(data))
	}
	
	return nil
}

func processUserRequest() {
	// Nested callbacks - callback hell
	fetchUser(1, func(user User) {
		fetchOrders(user.ID, func(orders []Order) {
			fetchProducts(orders[0].ProductID, func(products []Product) {
				fmt.Println(products)
			})
		})
	})
}

func fetchUser(id int, callback func(User)) {
	callback(User{ID: id, Name: "John"})
}

func fetchOrders(userID int, callback func([]Order)) {
	callback([]Order{{ID: 1, ProductID: 100}})
}

func fetchProducts(productID int, callback func([]Product)) {
	callback([]Product{{ID: productID, Name: "Widget"}})
}

type User struct {
	ID   int
	Name string
}

type Order struct {
	ID        int
	ProductID int
}

type Product struct {
	ID   int
	Name string
}
