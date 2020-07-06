# Golang-rest-api

This is a simple REST API implementation in golang build around all the necessary elements for a stable robust API. 

As far as its features is concerned it is providing an endpoint for read,write and delete over a database. But what more important is its package structuring, responses, logging and testing. It covers all these sections quite well.


Default Port : 8082

URL          : /orders

Request Body : '{"username":"Master Gogo","meal":"dinner","date":"2019-08-12"}'

Database used : Postgres

Examples :

FETCH Records : 

    $ curl http://localhost:8082/orders
  
Add Record    : 

    $ curl -X POST -d '{"username":"Master Gogo","meal":"dinner","date":"2019-08-12"}' http://localhost:8082/orders
  
DELETE Record : 

    $ curl -X DELETE -d '{"username":"Master Gogo","meal":"dinner","date":"2019-08-12"}' http://localhost:8082/orders


# Start Application

You can start this application either with manual go commands which will require a golang installation on your system along with a postgre database.

Or

You can pull the docker image from docker hub and start the container where everything will be pre cooked.

https://hub.docker.com/r/shubham1962/golang-rest-api

# Run Manually
Clone the project repository. If its not cloned in your gopath then set a gopath for the repo folder.

Install postgres and import the database dump provided under db/orders.sql in 2 different databases. One for application and one for test cases.

Update configs user config/ accordingly with the credentials. 

Run test casea : 

    $ go test 
    $ go test -bench=.

If test cases are passed :

    $ go run main.go
    
or

    $ go build
    $ ./golang-rest-api
              
              
# Docker

  Pull image
  
    $ docker pull shubham1962/golang-rest-api

  Run container
  
    $ docker run --name rest -p 8082:8082 shubham1962/golang-rest-api
