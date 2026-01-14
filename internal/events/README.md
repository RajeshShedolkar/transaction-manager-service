Creating Kfka Topics:
docker exec -it kafka kafka-topics \
  --create \
  --topic test-topic \
  --bootstrap-server localhost:9092 \
  --partitions 1 \
  --replication-factor 1


docker exec -it kafka kafka-topics \
  --create \
  --topic test-topic1 \
  --bootstrap-server localhost:9092 \
  --partitions 1 \
  --replication-factor 1

-----------------


docker exec -it kafka kafka-topics \
  --create \
  --topic card-auth-events \
  --bootstrap-server localhost:9092 \
  --partitions 1 \
  --replication-factor 1

-----------------

List topics
docker exec -it kafka kafka-topics \
  --list \
  --bootstrap-server localhost:9092

Kafka Topics (Design)
Topic	Purpose
card-auth-events	CARD_AUTH
card-settlement-events	CARD_SETTLEMENT
card-release-events	CARD_AUTH_RELEASE
neft-settlement-events	NEFT_SETTLEMENT
compensation-events	CREDIT_REVERSAL_*

// To execute tests in go
go test ./internal/events -v