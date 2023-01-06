# Golang-rest-api

This is a simple REST API in golang build around all basic elements required in an API. 

As far as its features are concerned, it provides endpoints for read,write and delete over a postgresql database.


Default Port : 8082

URL          : /orders

Request Body : '{"username":"Master Gogo","meal":"dinner","date":"2019-08-12"}'

Database used : Postgres

Examples :

FETCH Records : 

    $ curl http://localhost:8082/orders
  
Add Record    : 

    $ curl -X POST -d '{"username":"slasher","meal":"dinner","date":"2019-08-12"}' http://localhost:8082/orders
  
DELETE Record : 

    $ curl -X DELETE -d '{"username":"slasher","meal":"dinner","date":"2019-08-12"}' http://localhost:8082/orders


# Start Application

You can start the application either manually which will require a golang installation on your system along with a postgres database.

Or

You can run the docker image.

https://hub.docker.com/r/shubham1962/golang-rest-api

# Run Manually
Clone the repository. If its not cloned in your GOPATH then set it for the repo folder.

Install postgresql and import database dump provided under db/orders.sql.

Update config under /config accordingly with the credentials. 

Run test cases : 

    $ go test 
    $ go test -bench=.

To start the application :

    $ go run main.go
    
or

    $ go build
    $ ./golang-rest-api
              
              
# Docker

  Run container
  
    $ docker run --name rest -p 8082:8082 shubham1962/golang-rest-api
