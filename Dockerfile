FROM golang:1.14
WORKDIR /go/src/golang-rest-api
RUN apt-get update
RUN apt-get install -y curl man vim less netcat dnsutils postgresql-11
COPY etc/postgresql/11/main/pg_hba.conf /etc/postgresql/11/main/pg_hba.conf 
COPY . .
RUN /etc/init.d/postgresql restart ; createuser -U postgres test_user1 ; psql -U postgres -c "ALTER USER test_user1 WITH PASSWORD '1234'" ; createdb -U postgres orders; createdb -U postgres test1 ; psql -U postgres -d orders < db/orders.sql ; psql -U postgres -d test1 < db/orders.sql
RUN go install -v ./...
RUN /etc/init.d/postgresql restart ; go test --bench=.
ENTRYPOINT ["bash","db/startapp.sh"]
