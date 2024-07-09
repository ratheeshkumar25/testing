package main

import (
	"restaurant/database"
	"restaurant/helper"
	"restaurant/routes"
)

// Init function to handle application initialization tasks
func Init() {
	helper.LoadEnv()
	database.DBconnect()
	database.InitRedis()
}

func main() {
	//Perform application initialization
	Init()

	// Initialize routes.
	r := routes.Routes()
	r.LoadHTMLGlob("templates/*")

	//Run the engine the port 3000
	if err := r.Run(":3000"); err != nil {
		panic(err) //Handle error if unable to start server
	}

}
