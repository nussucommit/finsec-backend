# First time tasks
postgres:
	docker run --name finsec-backend_postgresql_1 -p 5432:5432 -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=password -d postgres:13

create_db:
	docker exec -it finsec-backend_postgresql_1 createdb --username=postgres --owner=postgres finsec
	docker kill finsec-backend_postgresql_1
	docker container rm finsec-backend_postgresql_1

drop_db:
	docker exec -it finsec-backend_postgresql_1 dropdb -U postgres -w finsec
	docker kill finsec-backend_postgresql_1
	docker container rm finsec-backend_postgresql_1

# Routine tasks
psql_shell:
	docker exec -it finsec-backend_postgresql_1 psql -U postgres -w finsec

migrate_up:
	docker run -v ${CURDIR}/migrations:/migrations --network host migrate/migrate -path=migrations/ -database "postgresql://postgres:password@localhost:5432/finsec?sslmode=disable" up

migrate_down:
	docker run -v ${CURDIR}/migrations:/migrations --network host migrate/migrate -path=migrations/ -database "postgresql://postgres:password@localhost:5432/finsec?sslmode=disable" down -all

build:
	docker-compose build

run:
	docker-compose up --build -d

stop:
	docker-compose down
