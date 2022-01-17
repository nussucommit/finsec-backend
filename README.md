# finsec-backend
Backend for NUSSU Finsec web app. Built with Go.

## Getting Started

1. make postgres
2. make create_db
3. make run
4. Test: curl http://localhost:8080 and see if there is a hello world message to verify that it's working

Note: Steps 1 and 2 should only be done once, which is on the first time you clone this repo.
Afterwards, just do `make run` to get both the server and database up and running.

To stop the server: `make stop`

To access psql shell: `make psql_shell`