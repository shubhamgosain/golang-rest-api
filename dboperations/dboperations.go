package dboperations

import (
	"fmt"
	"database/sql"
	"log"
	"errors"
	"time"

	//For db connections to work
	_ "github.com/lib/pq"
)

//PostgresDB structure for database config 


var conn *sql.DB
	
type (
	configPath struct {
		path string
	}
)

//Createconnection creates a new database connection
func Createconnection(config PostgresDBStruct)  {
	psqlInfo := fmt.Sprintf("host=%s port=%v user=%s "+
    "password=%s dbname=%s sslmode=disable",
	config.Host, config.Port, config.User, config.Password, config.Name)
	db,_ := sql.Open("postgres", psqlInfo)
	if err := db.Ping() ; err !=nil{
		
		log.Panic(err)
	}
	log.Print("Successfully made connection to database")
	conn=db
}

//ReadRecords To Read the orders from database using GET requests
func ReadRecords() (orders []OrdersSchema){
	selDB, err := conn.Query("SELECT * FROM orders")
    if err != nil {
        panic(err.Error())
	}
	for selDB.Next() {
		var orderid int64
		var username,meal string
		var orderDate time.Time
		selDB.Scan(&orderid,&username,&orderDate,&meal)
		orders = append(orders,OrdersSchema{orderid,username,orderDate,meal})
	}
	return
}

//AddRecord Adds a new records to database
func AddRecord(data InputData) (error){
	if err :=valdidateInput(data,"check");err!=nil{
		return err
	}
	date,err := time.Parse("2006-01-02", data.Date)
	if err !=nil{
		return errors.New("Incorrect date format")
	}
	if err:=valdidateInput(data,"check-order");err==nil{
		return errors.New("Order for this date exist")
	}
	_, err = conn.Query("INSERT INTO orders(username,order_date,meal) VALUES($1,$2,$3)",data.Username,date,data.Meal)
    if err != nil {
		return err
	}
	return nil
}

//DeleteRecord deletes a new records to database
func DeleteRecord(data InputData) (error){
	if err :=valdidateInput(data,"check");err!=nil {
		return err
	}
	date,err := time.Parse("2006-01-02", data.Date)
	if err !=nil{
		return errors.New("Incorrect date format")
	}
	if err:=valdidateInput(data,"check-order");err!=nil{
		return err
	}
	_, err = conn.Query("DELETE from orders where username=$1 AND order_date=$2 AND meal=$3",data.Username,date,data.Meal)
    if err != nil {
		return err
	}
	return nil
}

func valdidateInput(data InputData,category string) (error){
	switch category{
		case "check-order":
			var orderid string
			date,_ := time.Parse("2006-01-02", data.Date)
			row := conn.QueryRow("SELECT orderid FROM orders where username=$1 and order_date=$2 and meal=$3",data.Username,date,data.Meal)
			switch err := row.Scan(&orderid); err {
			case sql.ErrNoRows:
			  return errors.New("Order does not exists")
			case nil:
			  return nil
			default:
			  panic(err)
			}
		default:
			if (data.Username == "" || data.Meal == "" || data.Date == "") {
				return errors.New("Incorrect data format")
			}else if (data.Meal != "lunch" && data.Meal != "dinner") {
				return errors.New("Meal can have values lunch/dinner")
			}
	}	
	return nil
}