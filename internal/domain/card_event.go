package domain

type CardEvent struct {
	EventID    string
	EventType  string
	NetworkTxn string
	Amount     int64
}
