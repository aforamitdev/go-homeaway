
### command to run application 

>  migrate create -seq -ext=.sql -dir=./migrations migrations_name

> migrate -path=./migrations -database=postgres://admin:admin@localhost:5432/linux?sslmode=disable