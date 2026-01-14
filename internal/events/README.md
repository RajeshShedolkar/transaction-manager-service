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
card-auth-events	AUTH_SUCCESS, AUTH_FAILED
card-settlement-events	SETTLEMENT_STARTED, DEBIT_CONFIRMED, CREDIT_CONFIRMED, SETTLEMENT_FAILED
card-release-events	REFUND_PROCESSED, 
neft-settlement-events	NEFT_SETTLEMENT_STARTED, NEFT_SETTLEMENT_DONE
compensation-events	CREDIT_REVERSAL_*

// To execute tests in go
go test ./internal/events -v