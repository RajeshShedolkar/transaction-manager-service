go get
go mod tidy
go run cmd/server/main.py

Migrate in Docker:
docker run --rm \                                  3s 23:15:38
  --network transaction_manager_default \
  -v $(pwd)/migrations:/migrations \
  migrate/migrate \
  -path /migrations \
  -database "postgres://transaction_service:pass123@postgres:5432/transactiondb?sslmode=disable" \
  up

sudo systemctl status postgres
sudo systemctl start postgres
sudo systemctl stop postgres

sudo ufw enable
sudo ufw status verbose
sudo ufw allow 8080

psql -h localhost -U <DBUSER> -d <DBNAME>