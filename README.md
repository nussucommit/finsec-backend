# finsec-backend
Backend for NUSSU Finsec web app. Built with Go.

## Getting Started

1. Create a `.env` file and copy the contents of `.env.example` into `.env`. Refer to the group chat for the values of each of them.
2. make postgres
3. make create_db
4. make run
5. Test: run curl http://localhost:8080 on your terminal and see if there is a hello world message to verify that it's working. Or can also open your browser and enter that link

Note: Steps 1 until 3 should only be done once, which is on the first time you clone this repo.
Afterwards, just do `make run` to get both the server and database up and running.

To stop the server: `make stop`

To access psql shell: `make psql_shell`