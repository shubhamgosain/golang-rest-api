package main

import (
	"flag"
	"fmt"
	"golang-rest-api/routes"
	"log"
)

var (
	logger        *log.Logger
	controlSocket string
	configFile    string
)

func restServer(configFile string) {
	routes.RestHandler(configFile)
}

func init() {
	flag.StringVar(&configFile, "c", "config/config.json", "Configuration file")
}

func main() {
	fmt.Println("\n------------------------------\n| Welcome to golang rest api |" +
		"\n------------------------------\n\n--- Quick start guide ---\n\n" +
		"Host         : localhost\nDefault Port : 8082\nURL          : /orders\n" +
		"Request Body : '{\"username\":\"Master Gogo\",\"meal\":\"dinner\",\"date\":\"2019-08-12\"}'\n\nDatabase used : Postgres\n\n\n" +
		"Examples :\n\nFETCH Records : curl http://localhost:8082/orders\n" +
		"Add Record    : curl -X POST -d '{\"username\":\"Master Gogo\",\"meal\":\"dinner\",\"date\":\"2019-08-12\"}' http://localhost:8082/orders\n" +
		"DELETE Record : curl -X DELETE -d '{\"username\":\"Master Gogo\",\"meal\":\"dinner\",\"date\":\"2019-08-12\"}' http://localhost:8082/orders\n\n\n" +
		"---- HAPPY PLAYING AROUND ----\n\n.")

	restServer(configFile)
}
