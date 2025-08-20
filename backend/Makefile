# Makefile
postgres-setup:
	docker-compose up -d

create-container-image:
	docker build -t  amulyam-backend .

remove-image:
	docker rmi  amulyam-backend

postgres-start:
	docker-compose start

postgres-destroy:
	docker-compose stop
	docker-compose rm -f
	rm -rf ~/amulyam/postgres

createdb:
	docker exec -it backend-amulyam-db-1  createdb --username=root --owner=root amulyam

dropdb:
	docker exec -it backend-amulyam-db-1   dropdb amulyam


migrate-create: 
	migrate create --ext sql -dir db/migration -seq init_schema

migrate-up:
	migrate -path db/migration -database "postgresql://root:admin@localhost:5435/amulyam?sslmode=disable" -verbose up

migrate-down:
	migrate -path db/migration -database "postgresql://root:admin@localhost:5435/amulyam?sslmode=disable" -verbose  down 1

sqlc:
	sqlc generate 

test:
	go test -v -cover ./...

console:
	docker exec -it go-bank-db-1 psql -U root -d backend-amulyam-db-1

start-server:
	go run main.go

mock:
	mockgen -destination db/mock/store.go -package mockdb  amulyam/db/sqlc Store 

.PHONY: createdb dropdb postgres-destroy postgres-setup migrate-up migrate-down sqlc test start-server mock create-container-image remove-image