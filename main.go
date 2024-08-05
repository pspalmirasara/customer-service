package main

import (
	"github.com/CAVAh/api-tech-challenge/src/infra/db/database"
	"github.com/CAVAh/api-tech-challenge/src/infra/web/routes"
)

func main() {
	database.ConnectDB()
	routes.HandleRequests()
}
