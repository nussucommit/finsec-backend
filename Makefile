postgres:
	docker run --name postgres -p 8080:5432 -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=password -d postgres:13

create_db:
	docker exec -it postgres createdb --username=postgres --owner=postgres finsec

psql_shell:
	docker exec -it postgres psql -U postgres

migrate_up:
	docker run --network host migrate/migrate -path=/migrations/ -database "postgresql://postgres:password@localhost:5432/commit?sslmode=disable" up

migrate_down:
	docker run --network host migrate/migrate -path=/migrations/ -database "postgresql://postgres:password@localhost:5432/commit?sslmode=disable" down

build_and_run:
	docker-compose up --build -d

stop:
	docker-compose down