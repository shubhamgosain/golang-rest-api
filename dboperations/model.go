package dboperations

import "time"

//OrdersSchema Schema for a database record
type OrdersSchema struct {
	Orderid   int64
	Username  string
	OrderDate time.Time
	Meal      string
}

//InputData Schema for request body
type InputData struct {
	Username string
	Date     string
	Meal     string
}

//DbConfig Schema for database config
type DbConfig struct {
	App struct {
		Port int `json:"port"`
	}
	PostgresDB PostgresDBStruct `json:"PostgresDB"`
}

//PostgresDBStruct schema to hold postgres db configs
type PostgresDBStruct struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Name     string `json:"name"`
}
