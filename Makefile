postgres:
	docker run --name postgres -p 5432:5432 -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=password -d postgres:13

create_db:
	docker exec -it postgres createdb --username=postgres --owner=postgres finsec

psql_shell:
	docker exec -it postgres psql -U postgres -W finsec

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
