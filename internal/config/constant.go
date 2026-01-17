package config

var ENV = "prod"
var KAFKA_ACCOUNT_TOPIC = "account.commands"
var KAFKA_CARD_EVENT_TOPIC = "card-auth-events"
var KAFKA_DLQ_TOPIC = "tm.events.dlq"
var KAFKA_BROKERS = []string{"localhost:9092"}
