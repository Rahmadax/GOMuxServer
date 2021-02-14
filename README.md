# GOMuxServer

OS independent guide:

run docker

cd api/database
docker-compose up
docker exec -it dbTest bash
mysql --user=root --password="Password123" dbTest < /migrations/migrations/init.sql

go run api/cmd/main.go