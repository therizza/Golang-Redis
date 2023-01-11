package main

import (
	"redistest/database"
	"redistest/routes"
)

func main() {
	database.DbConnect()
	routes.Routes()
}
